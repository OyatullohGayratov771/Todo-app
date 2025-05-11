// internal/service/user_service
package service

import (
	"context"
	"fmt"
	"user-service/internal/storage"
	"user-service/protos/user/userpb"
)

type UserService struct {
	storage storage.Storage
	userpb.UnimplementedUserServiceServer
}

func NewUserService(s *storage.PostgresStorage) *UserService {
	return &UserService{storage: s}
}

func (s *UserService) Register(ctx context.Context, req *userpb.RegisterUserRequest) (*userpb.RegisterUserResponse, error) {
	id, err := s.storage.InsertUser(ctx, req)
	if err != nil {
		return &userpb.RegisterUserResponse{}, err
	}
	return &userpb.RegisterUserResponse{
		UserID:   fmt.Sprintf("%d", id),
		Username: req.Username,
		Email:    req.Email,
	}, nil
}
