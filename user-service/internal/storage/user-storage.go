package storage

import (
	"context"
	"database/sql"
	"errors"
	"user-service/internal/utils"
	userpb "user-service/protos/user"

	"github.com/lib/pq"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Storage interface {
	InsertUser(ctx context.Context, req *userpb.RegisterUserReq) (int, error)
	LoginSql(ctx context.Context, req *userpb.LoginUserReq) (int, error)
	UpdateUserName(ctx context.Context, userID, newUserName string) error
	UpdatePassword(ctx context.Context, userID, oldPassword, newPassword string) error
	UpdateEmail(ctx context.Context, userID, newEmail string) error
}

type PostgresStorage struct {
	db *sql.DB
}

func NewPostgresStorage(db *sql.DB) *PostgresStorage {
	return &PostgresStorage{
		db: db,
	}
}

func (s *PostgresStorage) InsertUser(ctx context.Context, req *userpb.RegisterUserReq) (int, error) {
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

func (s *PostgresStorage) LoginSql(ctx context.Context, req *userpb.LoginUserReq) (int, error) {
	var storedPassword string
	var userID int

	err := s.db.QueryRowContext(ctx, "SELECT password,id FROM users WHERE username = $1", req.Username).
		Scan(&storedPassword, &userID)
	if err != nil {
		return 0, status.Error(codes.NotFound, "user not found")
	}

	if !utils.CheckPasswordHash(req.Password, storedPassword) {
		return 0, status.Error(codes.Unauthenticated, "invalid password")
	}

	return userID, nil
}

func (s *PostgresStorage) UpdateUserName(ctx context.Context, userID, newUserName string) error {
	var currentUsername string
	err := s.db.QueryRow("SELECT username FROM users WHERE id = $1", userID).Scan(&currentUsername)
	if err != nil {
		return err
	}

	if currentUsername == newUserName {
		return status.Error(codes.AlreadyExists, "username already exists")
	}

	_, err = s.db.Exec("UPDATE users SET username = $1 WHERE id = $2", newUserName, userID)
	return err
}

func (s *PostgresStorage) UpdatePassword(ctx context.Context, userID, oldPassword, newPassword string) error {
	var currentPassword string
	err := s.db.QueryRow("SELECT password FROM users WHERE id = $1", userID).Scan(&currentPassword)
	if err != nil {
		return err
	}
	check := utils.CheckPasswordHash(oldPassword, currentPassword)
	if !check {
		return status.Errorf(codes.FailedPrecondition, "The old password was entered incorrectly.")
	}
	check = utils.CheckPasswordHash(newPassword, currentPassword)
	if check {
		return status.Errorf(codes.FailedPrecondition, "The new password cannot be the same as the old one.")
	}
	hashnewpassword, err := utils.HashPassword(newPassword)

	_, err = s.db.Exec("UPDATE users SET password = $1 WHERE id = $2", hashnewpassword, userID)
	return err
}

func (s *PostgresStorage) UpdateEmail(ctx context.Context, userID, newEmail string) error {
	var currentEmail string
	err := s.db.QueryRow("SELECT email FROM users WHERE id = $1", userID).Scan(&currentEmail)
	if err != nil {
		return err
	}

	if currentEmail == newEmail {
		return status.Error(codes.AlreadyExists, "email already exists")
	}

	_, err = s.db.Exec("UPDATE users SET email = $1 WHERE id = $2", newEmail, userID)
	return err
}
