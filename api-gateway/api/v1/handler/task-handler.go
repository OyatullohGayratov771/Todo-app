package handler

import (
	"api-gateway/api/v1/models"
	"api-gateway/internal/utils"
	taskpb "api-gateway/protos/task"
	"fmt"
	"net/http"
	"strconv"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/gin-gonic/gin"
)

// CreateTask godoc
// @Summary      Create a new task
// @Description  Allows an authenticated user to create a new task
// @Tags         task
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        body  body      taskpb.CreateTaskRequest  true  "Task details"
// @Success      200   {object}  taskpb.CreateTaskResponse
// @Failure      400   {object}  map[string]string
// @Failure      401   {object}  map[string]string
// @Failure      404   {object}  map[string]string
// @Failure      500   {object}  map[string]string
// @Router       /task/create [post]
func (h *Handler) CreateTask(c *gin.Context) {
	ctx, err := utils.InjectTokenToContext(c.Request)
	if err != nil {
		c.JSON(401, gin.H{"error": "Unauthorized: invalid or missing token"})
		return
	}

	req := &taskpb.CreateTaskRequest{}
	if err := c.ShouldBindJSON(req); err != nil {
		c.JSON(400, gin.H{
			"error": "Invalid request format: " + err.Error(),
		})
		return
	}

	res, err := h.service.Task().CreateTask(ctx, req)
	if err != nil {
		grpcCode := status.Code(err)
		errMsg := status.Convert(err).Message()
		fmt.Println(err)
		switch grpcCode {
		case codes.Unauthenticated:
			c.JSON(401, gin.H{"error": "Unauthorized: invalid or missing token"})
		case codes.NotFound:
			c.JSON(404, gin.H{"error": errMsg})
		default:
			c.JSON(500, gin.H{"error": errMsg})
		}
		return
	}
	c.JSON(200, res)
}

// GetTasks godoc
// @Summary      Retrieve a task
// @Description  Allows an authenticated user to retrieve a task by ID
// @Tags         task
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      string  true  "Task ID"
// @Success      200  {object}  map[string]string
// @Failure      400  {object}  map[string]string
// @Failure      401  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /task/get/{id} [get]
func (h *Handler) GetTask(c *gin.Context) {
	// Tokenni tekshirish va contextga qo‘yish
	ctx, err := utils.InjectTokenToContext(c.Request)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized: invalid or missing token"})
		return
	}

	// Task ID olish
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Task ID is required"})
		return
	}
	if _, err := strconv.Atoi(id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Task ID must be a valid number"})
		return
	}

	// gRPC chaqiruvi
	res, err := h.service.Task().GetTask(ctx, &taskpb.GetTaskRequest{Id: id})
	if err != nil {
		grpcCode := status.Code(err)
		errMsg := status.Convert(err).Message()

		switch grpcCode {
		case codes.Unauthenticated:
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized: " + errMsg})
		case codes.NotFound:
			c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error: " + errMsg})
		}
		return
	}
	r := models.TaskResponse{
		ID:          res.Id,
		Title:       res.Title,
		Description: res.Description,
		Done:        res.Done,
		UserID:      res.UserId,
	}
	c.JSON(http.StatusOK, r)
}

// ListTasks godoc
// @Summary      Retrieve a list of tasks
// @Description  Allows an authenticated user to retrieve a list of tasks
// @Tags         task
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  taskpb.ListTasksResponse
// @Failure      401  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /task/list [get]
func (h *Handler) ListTasks(c *gin.Context) {
	ctx, err := utils.InjectTokenToContext(c.Request)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized: invalid or missing token"})
		return
	}

	res, err := h.service.Task().ListTasks(ctx, &taskpb.ListTasksRequest{})
	if err != nil {
		grpcCode := status.Code(err)
		errMsg := status.Convert(err).Message()
		switch grpcCode {
		case codes.Unauthenticated:
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized: " + errMsg})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error: " + errMsg})
		}
		return
	}
	Tasks := []models.TaskResponse{}
	for _, v := range res.Tasks {
		Tasks = append(Tasks, models.TaskResponse{
			ID:          v.Id,
			Title:       v.Title,
			Description: v.Description,
			Done:        v.Done,
			UserID:      v.UserId,
		})
	}

	c.JSON(http.StatusOK, gin.H{"tasks": Tasks})
}

// UpdateTask godoc
// @Summary      Update a task
// @Description  Allows an authenticated user to update a task by ID
// @Tags         task
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      string  true  "Task ID"
// @Param        body  body      models.TaskRequest  true  "Task details"
// @Success      200  {object}  map[string]string
// @Failure      400  {object}  map[string]string
// @Failure      401  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /task/update/{id} [put]
func (h *Handler) UpdateTask(c *gin.Context) {
	ctx, err := utils.InjectTokenToContext(c.Request)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized: invalid or missing token"})
		return
	}

	req := taskpb.UpdateTaskRequest{}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format: " + err.Error()})
		return
	}

	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Task ID is required"})
		return
	}
	if _, err := strconv.Atoi(id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Task ID must be a valid number"})
		return
	}

	req.Id = id

	res, err := h.service.Task().UpdateTask(ctx, &req)
	if err != nil {
		grpcCode := status.Code(err)
		errMsg := status.Convert(err).Message()
		switch grpcCode {
		case codes.Unauthenticated:
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized: " + errMsg})
		case codes.NotFound:
			c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error: " + errMsg})
		}
		return
	}

	c.JSON(http.StatusOK, res)
}

// DeleteTask godoc
// @Summary      delete a task
// @Description  Allows an authenticated user to delete a task by ID
// @Tags         task
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      string  true  "Task ID"
// @Success      200  {object}  map[string]string
// @Failure      400  {object}  map[string]string
// @Failure      401  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /task/delete/{id} [delete]
func (h *Handler) DeleteTask(c *gin.Context) {
	ctx, err := utils.InjectTokenToContext(c.Request)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized: invalid or missing token"})
		return
	}

	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Task ID is required"})
		return
	}
	if _, err := strconv.Atoi(id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Task ID must be a valid number"})
		return
	}

	req := taskpb.DeleteTaskRequest{
		Id: id,
	}
	res, err := h.service.Task().DeleteTask(ctx, &req)
	if err != nil {
		grpcCode := status.Code(err)
		errMsg := status.Convert(err).Message()
		switch grpcCode {
		case codes.Unauthenticated:
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized: " + errMsg})
		case codes.NotFound:
			c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error: " + errMsg})
		}
		return
	}

	c.JSON(http.StatusOK, res)
}
