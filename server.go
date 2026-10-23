package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/loadept/website/internal/config"
	"github.com/loadept/website/internal/middleware"
	"github.com/loadept/website/internal/short"
	"github.com/loadept/website/internal/website"
	"github.com/loadept/website/web/static"
	"zombiezen.com/go/sqlite"
	"zombiezen.com/go/sqlite/sqlitex"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	slog.SetDefault(logger)

	cfg, err := config.Load()
	fatalIfErr(err)

	pool, err := sqlitex.NewPool(cfg.DBPath, sqlitex.PoolOptions{
		PoolSize: 3,
		PrepareConn: func(conn *sqlite.Conn) error {
			return sqlitex.ExecuteScript(conn, `
				PRAGMA foreign_keys = ON;
				PRAGMA busy_timeout = 5000;
				PRAGMA journal_mode = WAL;
			`, nil)
		},
	})
	fatalIfErr(err)
	defer pool.Close()

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	fileServ := http.FileServerFS(static.FS)
	shortHandler := short.NewHandler(short.NewRepo(pool))
	webHandler := website.NewHandler(cfg.AppURL)

	r.Get("/", webHandler.Index)
	r.Get("/go/{pkg}", webHandler.ServePkg)
	r.Get("/s/{code}", shortHandler.RedirectURL)
	r.Handle("GET /static/*", http.StripPrefix("/static/", fileServ))

	server := http.Server{
		Addr:         cfg.Addr,
		Handler:      r,
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

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	log.Println("shutting down server...")
	fatalIfErr(server.Shutdown(shutdownCtx))
}

func fatalIfErr(err error) {
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}
