package storage

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"zombiezen.com/go/sqlite"
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

var (
	ErrShortURLNotFound         = errors.New("short url does not exist or invalid short code")
	ErrShortURLInsertNotRow     = errors.New("insert short url returned no rows")
	ErrShortURLCodeExists       = errors.New("the provided shortcode already exists")
	ErrShortURLNameExists       = errors.New("the provided name already exists")
	ErrShortURLConstraintUnique = errors.New("a unique constraint was violated")
)

const (
	URLStatusActive  = "active"
	URLStatusDeleted = "deleted"
)

type ShortURL struct {
	ID          int64  `json:"id,omitempty"`
	Name        string `json:"name"`
	OriginalURL string `json:"original_url"`
	ShortCode   string `json:"short_code,omitempty"`
}

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

func (s *ShortURLStorage) SaveURL(ctx context.Context, data *ShortURL) (*ShortURL, error) {
	conn, err := s.db.Take(ctx)
	if err != nil {
		return nil, err
	}
	defer s.db.Put(conn)

	stmt, err := conn.Prepare(`
		INSERT INTO short_urls(name, original_url, short_code)
		VALUES ($1, $2, $3)
		RETURNING id, name, original_url, short_code
	`)
	if err != nil {
		return nil, err
	}
	defer stmt.Reset()

	stmt.SetText("$1", data.Name)
	stmt.SetText("$2", data.OriginalURL)
	stmt.SetText("$3", data.ShortCode)

	if hasRow, err := stmt.Step(); err != nil {
		if sqlite.ErrCode(err) == sqlite.ResultConstraintUnique {
			if strings.Contains(err.Error(), "short_urls.short_code") {
				return nil, ErrShortURLCodeExists
			}
			if strings.Contains(err.Error(), "short_urls.name") {
				return nil, ErrShortURLNameExists
			}
			return nil, ErrShortURLConstraintUnique
		}
		return nil, err
	} else if !hasRow {
		return nil, ErrShortURLInsertNotRow
	}

	insertedData := &ShortURL{
		ID:          stmt.GetInt64("id"),
		Name:        stmt.GetText("name"),
		OriginalURL: stmt.GetText("original_url"),
		ShortCode:   stmt.GetText("short_code"),
	}
	return insertedData, nil
}
