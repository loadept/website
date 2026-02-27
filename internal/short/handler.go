package short

import (
	"errors"
	"fmt"
	"net/http"
	"sync"

	"github.com/loadept/website/internal/logger"
	"github.com/loadept/website/internal/storage"
)

// local cache vars
var (
	cacheMu sync.RWMutex
	cache   = make(map[string]string)
)

type shortHandler struct {
	s *storage.ShortURLStorage
}

func NewShortHandler(s *storage.ShortURLStorage) *shortHandler {
	return &shortHandler{s: s}
}

func (h *shortHandler) RedirectURL(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	le := logger.FromContext(ctx)
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
		le.Error = fmt.Sprintf("failed to get URL from database: %v", err)
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
