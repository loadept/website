package shortener

import (
	"crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/loadept/loadept.com/internal/storage"
)

type shortHandler struct {
	s  *storage.ShortURLStorage
	sa *storage.AuthStorage
}

func NewShortHandler(s *storage.ShortURLStorage, sa *storage.AuthStorage) *shortHandler {
	return &shortHandler{s: s, sa: sa}
}

var (
	cacheMu    sync.RWMutex
	cache      = make(map[string]string)
	httpClient = &http.Client{
		Timeout: 10 * time.Second,
	}
)

func (h *shortHandler) RedirectURL(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	shortCode := r.PathValue("code")

	cacheMu.RLock()
	if cachedURL, ok := cache[shortCode]; ok {
		cacheMu.RUnlock()
		http.Redirect(w, r, cachedURL, http.StatusFound)
		return
	}
	cacheMu.RUnlock()

	originalURL, err := h.s.GetURL(ctx, shortCode)
	if err != nil {
		if errors.Is(err, storage.ErrShortURLNotFound) {
			http.NotFound(w, r)
			return
		}
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	cacheMu.Lock()
	cache[shortCode] = originalURL
	cacheMu.Unlock()

	http.Redirect(w, r, originalURL, http.StatusFound)
}

type m map[string]any

func (h *shortHandler) CreateURL(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	headerToken := r.Header.Get("Authorization")
	if !strings.HasPrefix(headerToken, "Bearer") {
		writeJSON(w, m{"detail": "invalid token format"}, http.StatusUnauthorized)
		return
	}
	token := strings.TrimPrefix(headerToken, "Bearer ")

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		"https://api.github.com/user",
		nil,
	)
	if err != nil {
		writeJSON(w, m{"detail": "failed to create request"}, http.StatusInternalServerError)
		return
	}

	req.Header.Set("Authorization", "Bearer "+token)
	resp, err := httpClient.Do(req)
	if err != nil {
		writeJSON(w, m{"detail": "failed to fetch user info"}, http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		writeJSON(w, m{"detail": "failed to fetch user info"}, resp.StatusCode)
		return
	}

	var githubUser struct {
		ID    int64  `json:"id"`
		Login string `json:"login"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&githubUser); err != nil {
		writeJSON(w, m{"detail": "failed to decode user info"}, http.StatusInternalServerError)
		return
	}

	_, err = h.sa.GetUserByGithubID(ctx, githubUser.ID)
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			writeJSON(w, m{"detail": "user not found"}, http.StatusUnauthorized)
			return
		}
		writeJSON(w, m{"detail": "failed to fetch user"}, http.StatusInternalServerError)
		return
	}

	var body storage.ShortURL
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeJSON(w, m{"detail": "invalid request body"}, http.StatusBadRequest)
		return
	}
	if body.Name == "" || body.OriginalURL == "" {
		writeJSON(w, m{"detail": "name and original_url are required"}, http.StatusBadRequest)
		return
	}

	parsedURL, err := url.Parse(body.OriginalURL)
	if err != nil || (parsedURL.Scheme != "http" && parsedURL.Scheme != "https") {
		writeJSON(w, m{"detail": "invalid original_url"}, http.StatusBadRequest)
		return
	}
	body.OriginalURL = parsedURL.String()

	if body.ShortCode == "" {
		code, err := generateBase62Code(4)
		if err != nil {
			writeJSON(w, m{"detail": "failed to generate short code"}, http.StatusInternalServerError)
			return
		}
		body.ShortCode = code
	} else if len(body.ShortCode) < 4 || len(body.ShortCode) > 10 {
		writeJSON(w, m{"detail": "short_code must be between 4 and 10 characters"}, http.StatusBadRequest)
		return
	} else if !isValidBase62(body.ShortCode) {
		writeJSON(w, m{"detail": "short_code must be alphanumeric"}, http.StatusBadRequest)
		return
	}

	shortURL, err := h.s.SaveURL(ctx, &body)
	if err != nil {
		if errors.Is(err, storage.ErrShortURLCodeExists) {
			writeJSON(w, m{"detail": "short_code already exists"}, http.StatusConflict)
			return
		}
		if errors.Is(err, storage.ErrShortURLNameExists) {
			writeJSON(w, m{"detail": "name already exists"}, http.StatusConflict)
			return
		}
		writeJSON(w, m{"detail": "failed to create short url"}, http.StatusInternalServerError)
		return
	}

	writeJSON(w, m{
		"message": "success",
		"url":     shortURL,
	}, http.StatusCreated)
}

func writeJSON(w http.ResponseWriter, response any, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	responseBytes, err := json.Marshal(response)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, `{"detail": "unexpected internal error"}`)
		return
	}
	w.WriteHeader(statusCode)
	w.Write(responseBytes)
}

const base62Chars = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

func generateBase62Code(length int) (string, error) {
	b := make([]byte, length)
	for i := range b {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(base62Chars))))
		if err != nil {
			return "", err
		}
		b[i] = base62Chars[num.Int64()]
	}
	return string(b), nil
}

func isValidBase62(s string) bool {
	for _, c := range s {
		if !strings.ContainsRune(base62Chars, c) {
			return false
		}
	}
	return true
}
