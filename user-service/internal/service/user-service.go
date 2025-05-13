// internal/service/user_service
package service

import (
	"context"
	"fmt"
	"user-service/internal/storage"
	"user-service/internal/utils"
	userpb "user-service/protos/user"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UserService struct {
	storage storage.Storage
	userpb.UnimplementedUserServiceServer
}

func NewUserService(s *storage.PostgresStorage) *UserService {
	return &UserService{storage: s}
}

func (s *UserService) Register(ctx context.Context, req *userpb.RegisterUserReq) (*userpb.RegisterUserRes, error) {
	// hashing password
	hashpassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}
	req.Password = hashpassword

	id, err := s.storage.InsertUser(ctx, req)
	if err != nil {
		return &userpb.RegisterUserRes{}, err
	}
	t, err := utils.GenerateJWT(fmt.Sprintf("%d", id))
	if err != nil {
		return nil, err
	}

	return &userpb.RegisterUserRes{
		Message: "registration successful",
		Token:   t,
	}, nil
}

func (s *UserService) Login(ctx context.Context, req *userpb.LoginUserReq) (*userpb.LoginUserRes, error) {
	userID, err := s.storage.LoginSql(ctx, req)
	if err != nil {
		return nil, err
	}
	if userID == 0 {
		return nil, err
	}

	gentoken, err := utils.GenerateJWT(fmt.Sprintf("%d", userID))
	if err != nil {
		return nil, err
	}

	return &userpb.LoginUserRes{
		Token: gentoken,
	}, nil
}

func (s *UserService) UpdateUserName(ctx context.Context, req *userpb.UpdateUserNameReq) (*userpb.UpdateRes, error) {
	if req.Newusername == "" {
		return &userpb.UpdateRes{Message: "username cannot be empty"}, nil
	}
	userID, ok := ctx.Value("userID").(string)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "user not authenticated")
	}

	err := s.storage.UpdateUserName(ctx, userID, req.Newusername)
	if err != nil {
		return &userpb.UpdateRes{Message: "failed update"}, err
	}
	return &userpb.UpdateRes{Message: "update user name successful"}, nil
}

func (s *UserService) UpdatePassword(ctx context.Context, req *userpb.UpdatePasswordReq) (*userpb.UpdateRes, error) {
	if req.Newpassword == "" {
		return &userpb.UpdateRes{Message: "enter new password"}, nil
	}

	userID, ok := ctx.Value("userID").(string)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "user not authenticated")
	}

	err := s.storage.UpdatePassword(ctx, userID, req.Oldpassword, req.Newpassword)
	if err != nil {
		return &userpb.UpdateRes{Message: "error in storage"}, err
	}

	return &userpb.UpdateRes{Message: "update password successful"}, nil
}

func (s *UserService) UpdateEmail(ctx context.Context, req *userpb.UpdateEmailReq) (*userpb.UpdateRes, error) {
	if req.Newemail == "" {
		return &userpb.UpdateRes{Message: "The email field is required."}, nil
	}

	userID, ok := ctx.Value("userID").(string)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "user not authenticated")
	}
	err := s.storage.UpdateEmail(ctx, userID, req.Newemail)
	if err != nil {
		return &userpb.UpdateRes{Message: "error in storage"}, err
	}

	return &userpb.UpdateRes{Message: "update email successful"}, nil
}
