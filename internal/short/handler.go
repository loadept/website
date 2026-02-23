package short

import (
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/loadept/website/internal/storage"
)

// local cache vars
var (
	cacheMu sync.RWMutex
	cache   = make(map[string]string)
	logMu   sync.Mutex
	enc     = json.NewEncoder(os.Stdout)
)

type shortHandler struct {
	s *storage.ShortURLStorage
}

func NewShortHandler(s *storage.ShortURLStorage) *shortHandler {
	return &shortHandler{s: s}
}

type logEntry struct {
	Timestamp string `json:"timestamp"`
	Method    string `json:"method,omitempty"`
	Path      string `json:"path,omitempty"`
	CacheHit  bool   `json:"cache_hit,omitempty"`
	Error     string `json:"error,omitempty"`
	Addr      string `json:"ip,omitempty"`
	Country   string `json:"country,omitempty"`
	RayID     string `json:"ray_id,omitempty"`
}

func newLogEntry(r *http.Request) *logEntry {
	// log entry (addr will be obtained from headers set by Cloudflare)
	return &logEntry{
		Timestamp: time.Now().Format(time.RFC3339),
		Method:    r.Method,
		Path:      r.URL.Path,
		Addr:      r.Header.Get("CF-Connecting-IP"),
		Country:   r.Header.Get("CF-IPCountry"),
		RayID:     r.Header.Get("CF-Ray"),
	}
}

func (h *shortHandler) RedirectURL(w http.ResponseWriter, r *http.Request) {
	le := newLogEntry(r)
	defer writeLog(le)
	ctx := r.Context()

	shortCode := r.PathValue("code")

	cacheMu.RLock()
	if cachedURL, ok := cache[shortCode]; ok {
		cacheMu.RUnlock()
		le.CacheHit = true
		http.Redirect(w, r, cachedURL, http.StatusFound)
		return
	}
	cacheMu.RUnlock()

	originalURL, err := h.s.GetURL(ctx, shortCode)
	if err != nil {
		le.Error = "failed to get URL from database: " + err.Error()
		if errors.Is(err, storage.ErrShortURLNotFound) {
			http.NotFound(w, r)
			return
		}
		http.Error(w, "500 internal error", http.StatusInternalServerError)
		return
	}
	cacheMu.Lock()
	cache[shortCode] = originalURL
	cacheMu.Unlock()

	http.Redirect(w, r, originalURL, http.StatusFound)
}

func writeLog(le *logEntry) {
	logMu.Lock()
	defer logMu.Unlock()
	enc.Encode(le)
}
