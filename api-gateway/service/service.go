package service

import (
	"api-gateway/config"
	userpb "api-gateway/protos/user"

	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type IClients interface {
	User() userpb.UserServiceClient
}

type ServiceManager struct {
	userService userpb.UserServiceClient
}

func New(cfg config.Config) (*ServiceManager, error) {
	connUser, err := grpc.NewClient(
		fmt.Sprintf("%s:%s", cfg.UserService.Host, cfg.UserService.Port),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to dial user service: %v", err)
	}

	return &ServiceManager{
		userService: userpb.NewUserServiceClient(connUser),
	}, nil
}

func (s *ServiceManager) User() userpb.UserServiceClient {
	return s.userService
}
