package main

import (
	"flag"
	"log"

	"zombiezen.com/go/sqlite"
	"zombiezen.com/go/sqlite/sqlitex"
)

const schema = `
CREATE TABLE IF NOT EXISTS users (
    id INTEGER PRIMARY KEY,
    identifier TEXT NOT NULL UNIQUE,
    github_id INTEGER NOT NULL UNIQUE,
    github_login TEXT NOT NULL,
    github_node_id TEXT NOT NULL,
    status TEXT NOT NULL CHECK(status IN ('active', 'deleted')) DEFAULT 'active',
    created_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TEXT 
) STRICT;
CREATE INDEX IF NOT EXISTS idx_users_identifier ON users(identifier);

CREATE TABLE IF NOT EXISTS short_urls (
    id INTEGER PRIMARY KEY,
    name TEXT NOT NULL UNIQUE,
    original_url TEXT NOT NULL,
    short_code TEXT NOT NULL UNIQUE,
    status TEXT NOT NULL CHECK(status IN ('active', 'deleted')) DEFAULT 'active',
    created_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TEXT
) STRICT;
CREATE INDEX IF NOT EXISTS idx_short_code ON short_urls(short_code);
CREATE INDEX IF NOT EXISTS idx_original_url ON short_urls(original_url);
`

func main() {
	log.SetFlags(0)

	dbPath := flag.String("db", "db.sqlite3", "Path to the SQLite database")
	flag.Parse()

	conn, err := sqlite.OpenConn(*dbPath, sqlite.OpenReadWrite, sqlite.OpenCreate)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	if err := sqlitex.ExecuteScript(conn, schema, nil); err != nil {
		log.Fatal(err)
	}
	log.Println("all migrations applied")
}
