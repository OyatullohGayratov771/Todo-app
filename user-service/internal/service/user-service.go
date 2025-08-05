// internal/service/user_service
package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
	"user-service/internal/kafka"
	"user-service/internal/storage"
	"user-service/internal/utils"
	userpb "user-service/protos/user"

	"github.com/go-redis/redis/v8"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UserService struct {
	storage storage.Storage
	rd      *redis.Client
	kf      *kafka.Producer
	userpb.UnimplementedUserServiceServer
}

func NewUserService(s *storage.PostgresStorage, rd *redis.Client, kf *kafka.Producer) *UserService {
	return &UserService{storage: s, rd: rd, kf: kf}
}

func (s *UserService) Register(ctx context.Context, req *userpb.RegisterUserReq) (*userpb.RegisterUserRes, error) {
	if req.Username == "" {
		return nil, status.Error(codes.InvalidArgument, "username cannot be empty")
	}
	if req.Email == "" {
		return nil, status.Error(codes.InvalidArgument, "email cannot be empty")
	}
	if req.Password == "" {
		return nil, status.Error(codes.InvalidArgument, "password cannot be empty")
	}

	hashpassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to hash password: %v", err)
	}
	req.Password = hashpassword

	id, err := s.storage.InsertUser(ctx, req)
	if err != nil {
		return nil, err
	}

	t, err := utils.GenerateJWT(fmt.Sprintf("%d", id))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to generate token: %v", err)
	}

	event := map[string]string{
		"subject":   "user_registered",
		"userEmail": req.Email,
		"message":   "Welcome " + req.Username + "! Your registration is successful.",
	}
	data, _ := json.Marshal(event)
	s.kf.Publish(data)

	return &userpb.RegisterUserRes{
		Message: "registration successful",
		Token:   t,
	}, nil
}

func (s *UserService) Login(ctx context.Context, req *userpb.LoginUserReq) (*userpb.LoginUserRes, error) {
	ipKey := "login_attempts:" + req.Username

	attempts, err := s.rd.Incr(ctx, ipKey).Result()
	if err != nil {
		return nil, status.Error(codes.Internal, "redis error")
	}

	if attempts == 1 {
		s.rd.Expire(ctx, ipKey, time.Minute)
	}

	if attempts > 5 {
		return nil, status.Error(codes.ResourceExhausted, "Too many login attempts. Try again later.")
	}

	userID, email, err := s.storage.LoginSql(ctx, req)
	if err != nil {
		return nil, err
	}

	gentoken, err := utils.GenerateJWT(fmt.Sprintf("%d", userID))
	if err != nil {
		return nil, err
	}

	event := map[string]string{
		"subject":   "user_logged_in",
		"userEmail": email,
		"message":   "User successfully logged in!",
	}
	data, _ := json.Marshal(event)
	s.kf.Publish(data)

	return &userpb.LoginUserRes{Token: gentoken}, nil
}

func (s *UserService) UpdateUserName(ctx context.Context, req *userpb.UpdateUserNameReq) (*userpb.UpdateRes, error) {
	if req.Newusername == "" {
		return &userpb.UpdateRes{Message: "enter new username"}, status.Error(codes.InvalidArgument, "username cannot be empty")
	}
	userID, ok := ctx.Value("userID").(string)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "user not authenticated")
	}

	err := s.storage.UpdateUserName(ctx, userID, req.Newusername)
	if err != nil {
		return &userpb.UpdateRes{Message: "failed update"}, err
	}
	s.rd.Del(ctx, "user:"+userID)
	return &userpb.UpdateRes{Message: "update user name successful"}, nil
}

func (s *UserService) UpdatePassword(ctx context.Context, req *userpb.UpdatePasswordReq) (*userpb.UpdateRes, error) {
	if req.Newpassword == "" {
		return &userpb.UpdateRes{Message: "enter new password"}, status.Errorf(codes.InvalidArgument, "password cannot be empty")
	}

	userID, ok := ctx.Value("userID").(string)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "user not authenticated")
	}

	err := s.storage.UpdatePassword(ctx, userID, req.Currentpassword, req.Newpassword)
	if err != nil {
		return &userpb.UpdateRes{Message: "error in storage"}, err
	}

	return &userpb.UpdateRes{Message: "update password successful"}, nil
}

func (s *UserService) UpdateEmail(ctx context.Context, req *userpb.UpdateEmailReq) (*userpb.UpdateRes, error) {
	if req.Newemail == "" {
		return &userpb.UpdateRes{Message: "The email field is required."}, status.Errorf(codes.InvalidArgument, "email cannot be empty")
	}

	userID, ok := ctx.Value("userID").(string)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "user not authenticated")
	}
	err := s.storage.UpdateEmail(ctx, userID, req.Newemail)
	if err != nil {
		return &userpb.UpdateRes{Message: "error in storage"}, err
	}
	s.rd.Del(ctx, "user:"+userID)
	return &userpb.UpdateRes{Message: "update email successful"}, nil
}

func (s *UserService) GetProfile(ctx context.Context, req *userpb.GetProfileReq) (*userpb.GetProfileRes, error) {
	userID, ok := ctx.Value("userID").(string)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "user not authenticated")
	}

	// 1. Redis'dan tekshiramiz (cache bor yoki yo‘q)
	cached, err := s.rd.Get(ctx, "user:"+userID).Result()
	if err == nil {
		// Agar Redis'da bor bo‘lsa — JSON dan struct'ga parse qilib yuboramiz
		var cachedUser userpb.GetProfileRes

		json.Unmarshal([]byte(cached), &cachedUser)

		return &userpb.GetProfileRes{
			Username: cachedUser.Username,
			Email:    cachedUser.Email,
		}, nil
	}

	// 2. Redis'da topilmadi → PostgreSQL'dan olib kelamiz
	user, err := s.storage.GetUserByID(ctx, userID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get user: %v", err)
	}

	// 3. Redis'ga saqlab qo‘yamiz (cache qilish)
	data, _ := json.Marshal(user)
	s.rd.Set(ctx, "user:"+userID, data, time.Minute*10)

	return &userpb.GetProfileRes{
		Username: user.Username,
		Email:    user.Email,
	}, nil
}
