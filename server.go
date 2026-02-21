package main

import (
	"context"
	"embed"
	"errors"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path"
	"syscall"
	"time"

	"github.com/loadept/website/internal/shortener"
	"github.com/loadept/website/internal/storage"
	"zombiezen.com/go/sqlite"
	"zombiezen.com/go/sqlite/sqlitex"
)

//go:embed all:web/static
var staticFiles embed.FS

func main() {
	log.SetFlags(0)

	pool, err := sqlitex.NewPool(getEnv("DB_PATH"), sqlitex.PoolOptions{
		PoolSize: 3,
		PrepareConn: func(conn *sqlite.Conn) error {
			return sqlitex.ExecuteScript(conn, fmt.Sprintf(`
				PRAGMA foreign_keys = ON;
				PRAGMA busy_timeout = %d;
				PRAGMA journal_mode = WAL;
			`, 5000), nil)
		},
	})
	fatalIfErr(err)

	mux := http.NewServeMux()
	s, err := storage.NewShortURLStorage(pool)
	fatalIfErr(err)
	sa, err := storage.NewAuthStorage(pool)
	fatalIfErr(err)
	shortHandler := shortener.NewShortHandler(s, sa)

	subFS, err := fs.Sub(staticFiles, "web/static")
	fatalIfErr(err)
	staticFS := neuteredFS{fs: http.FS(subFS)}
	mux.Handle("GET /", http.FileServer(staticFS))
	mux.HandleFunc("GET /s/{code}", shortHandler.RedirectURL)
	mux.HandleFunc("POST /s/shorten", shortHandler.CreateURL)

	server := http.Server{
		Addr:         getEnv("ADDR"),
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}
	go func() {
		log.Println("server listen on", server.Addr)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("server error: %v", err)
		}
	}()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	<-ctx.Done()

	shotdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	log.Println("shutting down server...")
	fatalIfErr(server.Shutdown(shotdownCtx))
}

func getEnv(key string) string {
	value := os.Getenv(key)
	if value == "" {
		log.Fatalf("environment variable %s is required", key)
	}
	return value
}

func fatalIfErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

// FileServer Wrapper
type neuteredFS struct{ fs http.FileSystem }

func (n neuteredFS) Open(name string) (http.File, error) {
	f, err := n.fs.Open(name)
	if err != nil {
		return nil, err
	}

	stat, err := f.Stat()
	if err != nil {
		if errClose := f.Close(); errClose != nil {
			return nil, errors.Join(err, errClose)
		}
		return nil, err
	}

	if stat.IsDir() {
		index := path.Join(name, "index.html")
		indexFile, err := n.fs.Open(index)
		if err != nil {
			if errClose := f.Close(); errClose != nil {
				return nil, errors.Join(err, errClose)
			}
			return nil, fs.ErrNotExist
		}
		if errClose := indexFile.Close(); errClose != nil {
			return nil, errors.Join(err, errClose)
		}
	}

	return f, nil
}
