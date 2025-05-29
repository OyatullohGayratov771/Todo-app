package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"task-service/internal/storage"
	taskpb "task-service/protos/task"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UserService struct {
	storage storage.Storage
	taskpb.UnimplementedTaskServiceServer
}

func NewUserService(s *storage.PostgresStorage) *UserService {
	return &UserService{storage: s}
}

func (s *UserService) CreateTask(ctx context.Context, req *taskpb.CreateTaskRequest) (*taskpb.CreateTaskResponse, error) {
	userID, ok := ctx.Value("userID").(string)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "user not authenticated")
	}
	if req.Title == "" {
		return nil, status.Errorf(codes.InvalidArgument, "title is required")
	}

	id, err := s.storage.CreateTask(ctx, userID, req)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create task: %v", err)
	}

	return &taskpb.CreateTaskResponse{
		Id: fmt.Sprintf("%d", id),
	}, nil
}

func (s *UserService) GetTask(ctx context.Context, req *taskpb.GetTaskRequest) (*taskpb.Task, error) {
	userID, ok := ctx.Value("userID").(string)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "user not authenticated")
	}
	task, err := s.storage.GetTask(req.Id, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, status.Errorf(codes.NotFound, "task not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to get task: %v", err)
	}
	return task, nil
}

func (s *UserService) ListTasks(ctx context.Context, req *taskpb.ListTasksRequest) (*taskpb.ListTasksResponse, error) {
	userID, ok := ctx.Value("userID").(string)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "user not authenticated")
	}
	tasks, err := s.storage.ListTasks(ctx, userID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list tasks: %v", err)
	}
	return &taskpb.ListTasksResponse{
		Tasks: tasks,
	}, nil
}

func (s *UserService) UpdateTask(ctx context.Context, req *taskpb.UpdateTaskRequest) (*taskpb.UpdateTaskResponse, error) {
	userID, ok := ctx.Value("userID").(string)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "user not authenticated")
	}
	res, err := s.storage.UpdateTask(req, userID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update task: %v", err)
	}
	return res, nil
}

func (s *UserService) DeleteTask(ctx context.Context, req *taskpb.DeleteTaskRequest) (*taskpb.DeleteTaskResponse, error) {
	userID, ok := ctx.Value("userID").(string)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "user not authenticated")
	}
	res, err := s.storage.DeleteTask(req.Id, userID)
	if err != nil {
		return nil, err
	}
	return res, nil
}
