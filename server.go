package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os/signal"
	"path"
	"syscall"
	"time"

	"github.com/loadept/loadept.com/internal/config"
	"github.com/loadept/loadept.com/internal/shortener"
	"github.com/loadept/loadept.com/internal/storage"
	"zombiezen.com/go/sqlite"
	"zombiezen.com/go/sqlite/sqlitex"
)

func main() {
	log.SetFlags(0)

	var confFile string
	flag.StringVar(&confFile, "config", "", "config file for server")
	flag.Parse()

	if confFile == "" {
		log.Fatal("config file not specified")
	}
	cfg, err := config.Load(confFile)
	fatalIfErr(err)

	pool, err := sqlitex.NewPool(cfg.Database.DBPath, sqlitex.PoolOptions{
		PoolSize: cfg.Database.PoolSize,
		PrepareConn: func(conn *sqlite.Conn) error {
			return sqlitex.ExecuteScript(conn, fmt.Sprintf(`
				PRAGMA foreign_keys = ON;
				PRAGMA busy_timeout = %d;
				PRAGMA journal_mode = WAL;
			`, cfg.Database.BusyTimeout), nil)
		},
	})
	fatalIfErr(err)

	mux := http.NewServeMux()
	s, err := storage.NewShortURLStorage(pool)
	fatalIfErr(err)
	sa, err := storage.NewAuthStorage(pool)
	fatalIfErr(err)
	shortHandler := shortener.NewShortHandler(s, sa)

	staticFs := neuteredFS{fs: http.Dir(cfg.App.StaticFiles)}
	mux.Handle("GET /", http.FileServer(staticFs))
	mux.HandleFunc("GET /{code}", shortHandler.RedirectURL)
	mux.HandleFunc("POST /shorten", shortHandler.CreateURL)

	server := http.Server{
		Addr:         cfg.HTTP.Addr,
		Handler:      mux,
		ReadTimeout:  time.Duration(cfg.HTTP.ReadTimeoutSeconds) * time.Second,
		WriteTimeout: time.Duration(cfg.HTTP.WriteTimeoutSeconds) * time.Second,
		IdleTimeout:  time.Duration(cfg.HTTP.IdleTimeoutSeconds) * time.Second,
	}
	go func() {
		log.Printf("server listen on %s\n", server.Addr)
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

func fatalIfErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

// FileSystem Wrapper
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
