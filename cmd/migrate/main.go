package main

import (
	"flag"
	"log"
	"os"

	"zombiezen.com/go/sqlite"
	"zombiezen.com/go/sqlite/sqlitex"
)

const schema = `
CREATE TABLE IF NOT EXISTS short_urls (
    id INTEGER PRIMARY KEY,
    name TEXT NOT NULL,
    original_url TEXT NOT NULL,
    short_code TEXT NOT NULL,
    status TEXT NOT NULL CHECK(status IN ('active', 'deleted')) DEFAULT 'active',
    created_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TEXT
) STRICT;
CREATE INDEX IF NOT EXISTS idx_original_url ON short_urls(original_url);

CREATE UNIQUE INDEX IF NOT EXISTS idx_short_urls_name_active
ON short_urls(name) WHERE status = 'active';

CREATE UNIQUE INDEX IF NOT EXISTS idx_short_urls_short_code_active
ON short_urls(short_code) WHERE status = 'active';
`

func main() {
	log.SetFlags(0)

	var db string
	var run, check bool
	flag.StringVar(&db, "db", os.Getenv("DB_PATH"), "Path to the SQLite database")
	flag.BoolVar(&run, "run", false, "")
	flag.BoolVar(&check, "check", false, "")
	flag.Parse()

	if db == "" {
		log.Fatal("db path is required: use -db flag or set DB_PATH env var")
	}
	if run && check {
		log.Fatal("use -run to apply migrations or -status to check schema")
	}
	if !run && !check {
		flag.Usage()
		os.Exit(1)
	}

	conn, err := sqlite.OpenConn(db, sqlite.OpenReadWrite, sqlite.OpenCreate)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	if run {
		if err := sqlitex.ExecuteScript(conn, schema, nil); err != nil {
			log.Fatal("failed to apply schema:", err)
		}

		log.Println("migrations applied successfully")
		return
	}
	if check {
		var found []string
		if err := sqlitex.Execute(conn, `
			SELECT name
			FROM sqlite_master
			WHERE tbl_name = 'short_urls'
		`, &sqlitex.ExecOptions{
			ResultFunc: func(stmt *sqlite.Stmt) error {
				found = append(found, stmt.GetText("name"))
				return nil
			},
		}); err != nil {
			log.Fatal(err)
		}

		if len(found) == 0 {
			log.Fatal("migrations not applied")
		}

		for _, name := range found {
			log.Println("✓", name)
		}
	}
}
