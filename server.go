package main

import (
	"context"
	"errors"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
	"github.com/loadept/website/internal/gopackages"
	"github.com/loadept/website/internal/middleware"
	"github.com/loadept/website/internal/short"
	"github.com/loadept/website/ui"
	"zombiezen.com/go/sqlite"
	"zombiezen.com/go/sqlite/sqlitex"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	slog.SetDefault(logger)

	_ = godotenv.Overload()

	db := os.Getenv("DB_PATH")
	addr := os.Getenv("ADDR")
	githubUser := os.Getenv("GITHUB_USER")
	if db == "" {
		log.Fatalf("env var DB_PATH is required")
	}
	if addr == "" {
		log.Fatalf("env var ADDR is required")
	}
	if githubUser == "" {
		log.Fatalf("env var GITHUB_USER is required")
	}

	pool, err := sqlitex.NewPool(db, sqlitex.PoolOptions{
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

	shortHandler := short.NewHandler(short.NewRepo(pool))
	goPkg := gopackages.NewHandler(githubUser)

	fileServ := http.FileServerFS(ui.FS)
	r.Get("/s/{code}", shortHandler.RedirectURL)
	r.Get("/go/{pkg}", goPkg.ServePkg)
	r.Get("/*", func(w http.ResponseWriter, r *http.Request) { fileServ.ServeHTTP(w, r) })

	server := http.Server{
		Addr:         addr,
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
		log.Fatal(err)
	}
}
