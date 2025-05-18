package v1

import (
	"api-gateway/service"

	"api-gateway/api/v1/handler"

	"github.com/gin-gonic/gin"
)

type Options struct {
	Service service.IClients
}

func Router(opt Options) *gin.Engine {
	r := gin.Default()

	h := handler.NewHandler(&handler.HandlerCinfig{
		Service: opt.Service,
	})

	user := r.Group("/user")

	user.POST("/register", h.RegisterUser)
	user.POST("/login", h.LoginUser)
	user.PUT("/updatename", h.UpdateUserName)
	user.PUT("/updatepassword", h.UpdatePassword)
	user.PUT("/updateemail", h.UpdateEmail)

	return r
}
