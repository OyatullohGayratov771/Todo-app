package storage

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"

	taskpb "task-service/protos/task"
)

type Storage interface {
	CreateTask(ctx context.Context, userID string, req *taskpb.CreateTaskRequest) (int, error)
	GetTask(id, userID string) (*taskpb.Task, error)
	ListTasks(ctx context.Context, userID string) ([]*taskpb.Task, error)
	UpdateTask(req *taskpb.UpdateTaskRequest, userID string) (*taskpb.UpdateTaskResponse, error)
	DeleteTask(id, userID string) (*taskpb.DeleteTaskResponse, error)
}

type PostgresStorage struct {
	db *sql.DB
}

func NewPostgresStorage(db *sql.DB) *PostgresStorage {
	return &PostgresStorage{
		db: db,
	}
}

func (s *PostgresStorage) CreateTask(ctx context.Context, userID string, req *taskpb.CreateTaskRequest) (int, error) {
	var id int
	query := `
		INSERT INTO tasks (user_id, title, description, done, created_at, updated_at)
		VALUES ($1, $2, $3, false, NOW(), NOW())
		RETURNING id
	`

	i, err := strconv.Atoi(userID)
	if err != nil {
		return 0, err
	}

	err = s.db.QueryRowContext(ctx, query,
		i,
		req.Title,
		req.Description,
	).Scan(&id)

	if err != nil {
		return 0, err
	}

	return id, nil
}

func (p *PostgresStorage) GetTask(id, userID string) (*taskpb.Task, error) {
	var task taskpb.Task
	var taskID, dbUserID int

	err := p.db.QueryRow(`
		SELECT id, user_id, title, description, done 
		FROM tasks WHERE id = $1
	`, id).Scan(&taskID, &dbUserID, &task.Title, &task.Description, &task.Done)

	if err != nil {
		return nil, err
	}

	if fmt.Sprintf("%d", dbUserID) != userID {
		return nil, fmt.Errorf("access denied: task does not belong to user")
	}

	task.Id = fmt.Sprintf("%d", taskID)
	task.UserId = fmt.Sprintf("%d", dbUserID)

	return &task, nil
}

func (p *PostgresStorage) ListTasks(ctx context.Context, userID string) ([]*taskpb.Task, error) {
	rows, err := p.db.Query(`
		SELECT id, title, description, done 
		FROM tasks WHERE user_id = $1
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []*taskpb.Task

	for rows.Next() {
		var task taskpb.Task
		var id int

		err := rows.Scan(&id, &task.Title, &task.Description, &task.Done)
		if err != nil {
			return nil, err
		}

		task.Id = fmt.Sprintf("%d", id)
		task.UserId = userID

		tasks = append(tasks, &task)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return tasks, nil
}

func (p *PostgresStorage) UpdateTask(req *taskpb.UpdateTaskRequest, userID string) (*taskpb.UpdateTaskResponse, error) {
	var existing taskpb.Task

	// Eski qiymatlarni faqat user_id tekshiruvi bilan olish
	err := p.db.QueryRow(`SELECT title, description FROM tasks WHERE id = $1 AND user_id = $2`, req.Id, userID).
		Scan(&existing.Title, &existing.Description)
	if err != nil {
		return nil, fmt.Errorf("task not found or you are not authorized: %v", err)
	}

	// Bo'sh joylarni saqlab qolamiz
	if req.Title == "" {
		req.Title = existing.Title
	}
	if req.Description == "" {
		req.Description = existing.Description
	}

	// Update faqat user_id to‘g‘ri bo‘lsa ishlaydi
	result, err := p.db.Exec(`
		UPDATE tasks 
		SET title = $1, description = $2, done = $3, updated_at = NOW()
		WHERE id = $4 AND user_id = $5
	`, req.Title, req.Description, req.Done, req.Id, userID)

	if err != nil {
		return nil, fmt.Errorf("failed to update task: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, fmt.Errorf("failed to check rows affected: %v", err)
	}

	if rowsAffected == 0 {
		return nil, fmt.Errorf("no task updated")
	}

	return &taskpb.UpdateTaskResponse{
		Message: "Task updated successfully",
	}, nil
}

func (p *PostgresStorage) DeleteTask(id, userID string) (*taskpb.DeleteTaskResponse, error) {
	result, err := p.db.Exec(`
		DELETE FROM tasks 
		WHERE id = $1 AND user_id = $2
	`, id, userID)

	if err != nil {
		return nil, fmt.Errorf("failed to delete task: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, fmt.Errorf("failed to check rows affected: %v", err)
	}

	if rowsAffected == 0 {
		return nil, fmt.Errorf("no task deleted")
	}

	return &taskpb.DeleteTaskResponse{
		Message: "Task deleted successfully",
	}, nil
}
