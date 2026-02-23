package storage

import (
	"context"
	"errors"
	"fmt"

	"zombiezen.com/go/sqlite/sqlitex"
)

type ShortURLStorage struct {
	db *sqlitex.Pool
}

func NewShortURLStorage(db *sqlitex.Pool) (*ShortURLStorage, error) {
	if db == nil {
		return nil, fmt.Errorf("db pool is nil")
	}
	return &ShortURLStorage{db: db}, nil
}

var ErrShortURLNotFound = errors.New("short url does not exist or invalid short code")

const (
	URLStatusActive  = "active"
	URLStatusDeleted = "deleted"
)

func (s *ShortURLStorage) GetURL(ctx context.Context, shortCode string) (string, error) {
	conn, err := s.db.Take(ctx)
	if err != nil {
		return "", err
	}
	defer s.db.Put(conn)

	stmt, err := conn.Prepare(`
		SELECT original_url
		FROM short_urls
		WHERE short_code = $1
		AND status = $2
	`)
	if err != nil {
		return "", err
	}
	defer stmt.Reset()

	stmt.SetText("$1", shortCode)
	stmt.SetText("$2", URLStatusActive)

	if hasRow, err := stmt.Step(); err != nil {
		return "", err
	} else if !hasRow {
		return "", ErrShortURLNotFound
	}
	return stmt.GetText("original_url"), nil
}
