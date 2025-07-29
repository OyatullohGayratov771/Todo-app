package service

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"task-service/internal/storage"
	taskpb "task-service/protos/task"
	"time"

	"github.com/go-redis/redis/v8"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type TaskService struct {
	storage storage.Storage
	rd      *redis.Client
	taskpb.UnimplementedTaskServiceServer
}

func NewTaskService(s *storage.PostgresStorage, rd *redis.Client) *TaskService {
	return &TaskService{storage: s, rd: rd}
}

func (s *TaskService) CreateTask(ctx context.Context, req *taskpb.CreateTaskRequest) (*taskpb.CreateTaskResponse, error) {
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

	// Cache invalidatsiya qilish (ListTasks qayta yuklanadi)
	s.rd.Del(ctx, "tasks:"+userID)

	return &taskpb.CreateTaskResponse{
		Id: fmt.Sprintf("%d", id),
	}, nil
}

func (s *TaskService) GetTask(ctx context.Context, req *taskpb.GetTaskRequest) (*taskpb.Task, error) {
	userID, ok := ctx.Value("userID").(string)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "user not authenticated")
	}

	cacheKey := "task:" + userID + ":" + req.Id
	cached, err := s.rd.Get(ctx, cacheKey).Result()
	if err == nil {
		var task taskpb.Task
		json.Unmarshal([]byte(cached), &task)
		return &task, nil
	}

	task, err := s.storage.GetTask(req.Id, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, status.Errorf(codes.NotFound, "task not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to get task: %v", err)
	}

	data, _ := json.Marshal(task)
	s.rd.Set(ctx, cacheKey, data, time.Minute*5)

	return task, nil
}

func (s *TaskService) ListTasks(ctx context.Context, req *taskpb.ListTasksRequest) (*taskpb.ListTasksResponse, error) {
	userID, ok := ctx.Value("userID").(string)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "user not authenticated")
	}

	cacheKey := "tasks:" + userID
	cached, err := s.rd.Get(ctx, cacheKey).Result()
	if err == nil {
		var tasks []*taskpb.Task
		json.Unmarshal([]byte(cached), &tasks)
		return &taskpb.ListTasksResponse{Tasks: tasks}, nil
	}

	tasks, err := s.storage.ListTasks(ctx, userID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list tasks: %v", err)
	}

	data, _ := json.Marshal(tasks)
	s.rd.Set(ctx, cacheKey, data, time.Minute*5)

	return &taskpb.ListTasksResponse{
		Tasks: tasks,
	}, nil
}

func (s *TaskService) UpdateTask(ctx context.Context, req *taskpb.UpdateTaskRequest) (*taskpb.UpdateTaskResponse, error) {
	userID, ok := ctx.Value("userID").(string)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "user not authenticated")
	}

	res, err := s.storage.UpdateTask(req, userID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update task: %v", err)
	}

	// Cache invalidatsiya
	s.rd.Del(ctx, "tasks:"+userID)
	s.rd.Del(ctx, "task:"+userID+":"+req.Id)

	return res, nil
}

func (s *TaskService) DeleteTask(ctx context.Context, req *taskpb.DeleteTaskRequest) (*taskpb.DeleteTaskResponse, error) {
	userID, ok := ctx.Value("userID").(string)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "user not authenticated")
	}

	res, err := s.storage.DeleteTask(req.Id, userID)
	if err != nil {
		return nil, err
	}

	// Cache invalidatsiya
	s.rd.Del(ctx, "tasks:"+userID)
	s.rd.Del(ctx, "task:"+userID+":"+req.Id)

	return res, nil
}
