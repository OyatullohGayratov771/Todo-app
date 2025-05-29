package v1

import (
	"api-gateway/api/v1/handler"
	"api-gateway/service"

	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type Options struct {
	Service service.IClients
}

func Router(opt Options) *gin.Engine {
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"}, // frontend URL
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	h := handler.NewHandler(&handler.HandlerConfig{
		Service: opt.Service,
	})

	user := r.Group("/user")
	{
		user.POST("/register", h.RegisterUser)
		user.POST("/login", h.LoginUser)
		user.PUT("/updatename", h.UpdateUserName)
		user.PUT("/updatepassword", h.UpdatePassword)
		user.PUT("/updateemail", h.UpdateEmail)
	}

	task := r.Group("/task")
	{
		task.POST("/create", h.CreateTask)
		task.GET("/get/:id", h.GetTask)
		task.GET("/list", h.ListTasks)
		task.PUT("/update/:id", h.UpdateTask)
		task.DELETE("/delete/:id", h.DeleteTask)
	}

	return r
}
