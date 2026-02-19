package storage

import (
	"context"
	"fmt"

	"zombiezen.com/go/sqlite/sqlitex"
)

type AuthStorage struct {
	db *sqlitex.Pool
}

func NewAuthStorage(db *sqlitex.Pool) (*AuthStorage, error) {
	if db == nil {
		return nil, fmt.Errorf("db pool is nil")
	}
	return &AuthStorage{db: db}, nil
}

var ErrUserNotFound = fmt.Errorf("user not found")

const (
	UserStatusActive  = "active"
	UserStatusDeleted = "deleted"
)

type User struct {
	ID           int64
	Identifier   string
	GithubID     int64
	GithubLogin  string
	GithubNodeID string
}

func (s *AuthStorage) GetUserByGithubID(ctx context.Context, githubID int64) (*User, error) {
	conn, err := s.db.Take(ctx)
	if err != nil {
		return nil, err
	}
	defer s.db.Put(conn)

	stmt, err := conn.Prepare(`
		SELECT
			id,
			github_id,
			identifier
		FROM users
		WHERE github_id = $1
		AND status = $2
	`)
	if err != nil {
		return nil, err
	}
	defer stmt.Reset()

	stmt.SetInt64("$1", githubID)
	stmt.SetText("$2", UserStatusActive)

	if hasRow, err := stmt.Step(); err != nil {
		return nil, err
	} else if !hasRow {
		return nil, ErrUserNotFound
	}
	return &User{
		ID:         stmt.GetInt64("id"),
		Identifier: stmt.GetText("identifier"),
		GithubID:   stmt.GetInt64("github_id"),
	}, nil
}
