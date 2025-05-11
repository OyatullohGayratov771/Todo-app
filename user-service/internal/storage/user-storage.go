package storage

import (
	"context"
	"database/sql"
	"errors"
	"user-service/protos/user/userpb"

	"github.com/lib/pq"
)

type Storage interface {
	InsertUser(ctx context.Context, req *userpb.RegisterUserRequest) (int, error)
}

type PostgresStorage struct {
	db *sql.DB
}

func NewPostgresStorage(db *sql.DB) *PostgresStorage {
	return &PostgresStorage{
		db: db,
	}
}

func (s *PostgresStorage) InsertUser(ctx context.Context, req *userpb.RegisterUserRequest) (int, error) {
	var userID int
	err := s.db.QueryRowContext(ctx, "INSERT INTO users (username, email, password) VALUES ($1, $2, $3) RETURNING id", req.Username, req.Email, req.Password).Scan(&userID)
	if err != nil {
		if pgErr, ok := err.(*pq.Error); ok {
			if pgErr.Code == "23505" {
				return 0, errors.New("username already exists")
			}
		}
		return 0, err
	}
	return userID, nil
}
