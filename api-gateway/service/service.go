package service

import (
	"api-gateway/config"
	taskpb "api-gateway/protos/task"
	userpb "api-gateway/protos/user"

	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type IClients interface {
	User() userpb.UserServiceClient
	Task() taskpb.TaskServiceClient
}

type ServiceManager struct {
	userService userpb.UserServiceClient
	taskService taskpb.TaskServiceClient
}

func New(cfg config.Config) (*ServiceManager, error) {
	connUser, err := grpc.NewClient(
		fmt.Sprintf("%s:%s", cfg.UserService.Host, cfg.UserService.Port),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to dial user service: %v", err)
	}

	connTask, err := grpc.NewClient(
		fmt.Sprintf("%s:%s", cfg.TaskService.Host, cfg.TaskService.Port),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to dial task service: %v", err)
	}

	return &ServiceManager{
		userService: userpb.NewUserServiceClient(connUser),
		taskService: taskpb.NewTaskServiceClient(connTask),
	}, nil
}

func (s *ServiceManager) User() userpb.UserServiceClient {
	return s.userService
}

func (s *ServiceManager) Task() taskpb.TaskServiceClient {
	return s.taskService
}
