package main

import (
	"crypto/rand"
	"flag"
	"log"
	"math/big"
	netURL "net/url"
	"os"
	"strings"

	"zombiezen.com/go/sqlite"
	"zombiezen.com/go/sqlite/sqlitex"
)

const base62Chars = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

func main() {
	log.SetFlags(0)

	var db, name, url, code string
	flag.StringVar(&db, "db", os.Getenv("DB_PATH"), "Path to the SQLite database (if not set, will use DB_PATH env var)")
	flag.StringVar(&name, "n", "", "Name of short url")
	flag.StringVar(&name, "name", "", "Name of short url")
	flag.StringVar(&url, "u", "", "Original URL to shorten")
	flag.StringVar(&url, "url", "", "Original URL to shorten")
	flag.StringVar(&code, "c", "", "Custom short code (optional)")
	flag.StringVar(&code, "code", "", "Custom short code (optional)")
	flag.Parse()

	if db == "" {
		log.Fatal("db path is required: use -db flag or set DB_PATH env var")
	}
	if name == "" {
		log.Fatal("-name or -n flag is required")
	}
	if url == "" {
		log.Fatal("-url or -u flag is required")
	}

	purl, err := netURL.ParseRequestURI(url)
	if err != nil {
		log.Fatal(err)
	}

	if purl.Scheme != "http" && purl.Scheme != "https" {
		log.Fatal("url must have http or https scheme")
	}
	url = purl.String()

	if code != "" {
		if !isValidBase62(code) {
			log.Fatal("custom code must be a valid base62 string")
		}
		if len(code) > 10 {
			log.Fatal("custom code must be at most 10 characters long")
		}
	} else {
		var err error
		code, err = generateBase62Code(4)
		if err != nil {
			log.Fatalf("failed to generate short code: %v", err)
		}
	}

	conn, err := sqlite.OpenConn(db, sqlite.OpenReadWrite)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	if err := sqlitex.Execute(conn, `
		INSERT INTO short_urls(name, original_url, short_code)
		VALUES ($1, $2, $3)
		RETURNING id, name, original_url, short_code
	`, &sqlitex.ExecOptions{
		Named: map[string]any{"$1": name, "$2": url, "$3": code},
		ResultFunc: func(stmt *sqlite.Stmt) error {
			log.Println("short url created successfully")
			log.Printf("id=%d name=%s code=%s url=%s",
				stmt.GetInt64("id"),
				stmt.GetText("name"),
				stmt.GetText("short_code"),
				stmt.GetText("original_url"),
			)
			return nil
		},
	}); err != nil {
		if sqlite.ErrCode(err) == sqlite.ResultConstraintUnique {
			if strings.Contains(err.Error(), "short_urls.short_code") {
				log.Fatal("the provided short code already exists")
			}
			if strings.Contains(err.Error(), "short_urls.name") {
				log.Fatal("the provided name already exists")
			}
			log.Fatal("unexpected constraint violation")
		}
		log.Fatalf("failed to insert short url: %v", err)
	}
}

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
