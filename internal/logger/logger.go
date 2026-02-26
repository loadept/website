// Package logger provides middleware for adding logs throughout the application
package logger

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"sync"
	"time"
)

var (
	logMu sync.Mutex
	enc   = json.NewEncoder(os.Stdout)
)

type contextKey struct{}

var logEntryKey = contextKey{}

type resWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *resWriter) WriteHeader(statusCode int) {
	rw.statusCode = statusCode
	rw.ResponseWriter.WriteHeader(statusCode)
}

type LogEntry struct {
	Timestamp  string `json:"timestamp"`
	Method     string `json:"method,omitempty"`
	Path       string `json:"path,omitempty"`
	StatusCode int    `json:"status_code,omitempty"`
	CacheHit   bool   `json:"cache_hit,omitempty"`
	Error      string `json:"error,omitempty"`
	Addr       string `json:"ip,omitempty"`
	Country    string `json:"country,omitempty"`
	RayID      string `json:"ray_id,omitempty"`
}

func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rw := &resWriter{ResponseWriter: w, statusCode: http.StatusOK}

		// log entry (addr, country, rayid will be obtained from headers set by Cloudflare)
		le := &LogEntry{
			Timestamp: time.Now().Format(time.RFC3339),
			Method:    r.Method,
			Path:      r.URL.Path,
			Addr:      r.Header.Get("CF-Connecting-IP"),
			Country:   r.Header.Get("CF-IPCountry"),
			RayID:     r.Header.Get("CF-Ray"),
		}
		defer writeLog(le)

		ctx := context.WithValue(r.Context(), logEntryKey, le)
		next.ServeHTTP(rw, r.WithContext(ctx))
		le.StatusCode = rw.statusCode
	})
}

func FromContext(ctx context.Context) *LogEntry {
	le := ctx.Value(logEntryKey).(*LogEntry)
	return le
}

func writeLog(le *LogEntry) {
	logMu.Lock()
	defer logMu.Unlock()
	enc.Encode(le)
}
