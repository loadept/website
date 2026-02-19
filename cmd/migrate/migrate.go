package main

import (
	"flag"
	"log"
	"os"
	"path/filepath"
	"sort"

	"github.com/loadept/loadept.com/internal/config"
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
	if err != nil {
		log.Fatalf("internal error: %v", err)
	}

	conn, err := sqlite.OpenConn(cfg.Database.DBPath, sqlite.OpenReadWrite, sqlite.OpenCreate)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	dirEntry, err := os.ReadDir(cfg.Database.MigrationsPath)
	if err != nil {
		log.Fatal(err)
	}
	sort.Slice(dirEntry, func(i, j int) bool {
		return dirEntry[i].Name() < dirEntry[j].Name()
	})

	for _, entry := range dirEntry {
		if entry.IsDir() || filepath.Ext(entry.Name()) != ".sql" {
			continue
		}

		migration := filepath.Join(cfg.Database.MigrationsPath, entry.Name())
		content, err := os.ReadFile(migration)
		if err != nil {
			log.Fatal(err)
		}
		if err := sqlitex.ExecuteScript(conn, string(content), nil); err != nil {
			log.Fatal(err)
		}
		log.Printf("applied migration: %s", entry.Name())
	}
	log.Println("all migrations applied")
}
