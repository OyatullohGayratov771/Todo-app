package handler

import (
	"api-gateway/internal/utils"
	userpb "api-gateway/protos/user"
	"fmt"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)


// RegisterUser godoc
// @Summary      Register new user
// @Description  Create a new user and return JWT token
// @Tags         user
// @Accept       json
// @Produce      json
// @Param        user  body      userpb.RegisterUserReq  true  "User registration info"
// @Success      200   {object}  map[string]interface{}
// @Failure      400   {object}  map[string]string
// @Failure      409   {object}  map[string]string
// @Failure      500   {object}  map[string]string
// @Router       /user/register [post]
func (h *Handler) RegisterUser(c *gin.Context) {
	req := &userpb.RegisterUserReq{}

	if err := c.ShouldBindJSON(req); err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}

	res, err := h.service.User().Register(c, req)
	if err != nil {
		grpcCode := status.Code(err)
		errMsg := status.Convert(err).Message()

		switch grpcCode {
		case codes.AlreadyExists:
			c.JSON(409, gin.H{"error": errMsg})
		case codes.InvalidArgument:
			c.JSON(400, gin.H{"error": errMsg})
		default:
			c.JSON(500, gin.H{"error": errMsg})
		}

		return
	}

	c.JSON(200, gin.H{
		"message": res.Message,
		"token":   res.Token,
	})

}


// LoginUser godoc
// @Summary      Login user
// @Description  Login user and return JWT token
// @Tags         user
// @Accept       json
// @Produce      json
// @Param        user  body      userpb.LoginUserReq  true  "User login info"
// @Success      200   {object}  map[string]interface{}
// @Failure      400   {object}  map[string]string
// @Failure      404   {object}  map[string]string
// @Failure      500   {object}  map[string]string
// @Router       /user/login [post]
func (h *Handler) LoginUser(c *gin.Context) {
	req := &userpb.LoginUserReq{}

	if err := c.ShouldBindJSON(req); err != nil {
		c.JSON(400, gin.H{
			"error": "invalid request format: " + err.Error(),
		})
		return
	}

	res, err := h.service.User().Login(c, req)
	if err != nil {
		grpcCode := status.Code(err)
		errMsg := status.Convert(err).Message()
		fmt.Println(err)
		switch grpcCode {
		case codes.NotFound:
			c.JSON(404, gin.H{"error": "user not found"})
		case codes.Unauthenticated:
			c.JSON(401, gin.H{"error": "invalid password"})
		case codes.InvalidArgument:
			c.JSON(400, gin.H{"error": errMsg})
		default:
			c.JSON(500, gin.H{"error": "internal server error"})
		}

		return
	}

	c.JSON(200, gin.H{
		"token": res.Token,
	})
}

// UpdateUserName godoc
// @Summary      Update user's name
// @Description  Allows an authenticated user to update their username
// @Tags         user
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        body  body      userpb.UpdateUserNameReq  true  "New username"
// @Success      200   {object}  map[string]string
// @Failure      400   {object}  map[string]string
// @Failure      401   {object}  map[string]string
// @Failure      404   {object}  map[string]string
// @Failure      409   {object}  map[string]string
// @Failure      500   {object}  map[string]string
// @Router       /user/updatename [put]
func (h *Handler) UpdateUserName(c *gin.Context) {
	req := &userpb.UpdateUserNameReq{}

	// JSON parsing va validatsiya
	if err := c.ShouldBindJSON(req); err != nil {
		c.JSON(400, gin.H{
			"error": "Invalid request format: " + err.Error(),
		})
		return
	}

	// Tokenni olib, contextga joylash
	ctx, err := utils.InjectTokenToContext(c.Request)
	if err != nil {
		c.JSON(401, gin.H{"error": "Unauthorized: invalid or missing token"})
		return
	}

	// RPC chaqiruv
	res, err := h.service.User().UpdateUserName(ctx, req)
	if err != nil {
		grpcCode := status.Code(err)
		errMsg := status.Convert(err).Message() // gRPC error message
		switch grpcCode {
		case codes.NotFound:
			c.JSON(404, gin.H{"error": "User not found"})
		case codes.AlreadyExists:
			c.JSON(409, gin.H{"error": "Username already exists"})
		case codes.InvalidArgument:
			c.JSON(400, gin.H{"error": "Invalid request: " + errMsg})
		case codes.Unauthenticated:
			c.JSON(401, gin.H{"error": "Unauthorized: " + errMsg})
		default:
			c.JSON(500, gin.H{"error": "Internal server error: " + errMsg})
		}
		return
	}

	// Muvaffaqiyatli javob
	c.JSON(200, gin.H{
		"message": res.Message,
	})
}

// UpdatePassword godoc
// @Summary      Update user's name
// @Description  Allows an authenticated user to update their password
// @Tags         user
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        body  body      userpb.UpdatePasswordReq  true  "New password"
// @Success      200   {object}  map[string]string
// @Failure      400   {object}  map[string]string
// @Failure      401   {object}  map[string]string
// @Failure      404   {object}  map[string]string
// @Failure      500   {object}  map[string]string
// @Router       /user/updatepassword [put]
func (h *Handler) UpdatePassword(c *gin.Context) {
	req := &userpb.UpdatePasswordReq{}

	// JSON parsing va validatsiya
	if err := c.ShouldBindJSON(req); err != nil {
		c.JSON(400, gin.H{
			"error": "Invalid request format: " + err.Error(),
		})
		return
	}

	// Tokenni olib, contextga joylash
	ctx, err := utils.InjectTokenToContext(c.Request)
	if err != nil {
		c.JSON(401, gin.H{"error": "Unauthorized: invalid or missing token"})
		return
	}

	// RPC chaqiruv
	res, err := h.service.User().UpdatePassword(ctx, req)
	if err != nil {
		grpcCode := status.Code(err)
		errMsg := status.Convert(err).Message() // gRPC error message
		switch grpcCode {
		case codes.NotFound:
			c.JSON(404, gin.H{"error": "User not found"})
		case codes.InvalidArgument:
			c.JSON(400, gin.H{"error": "Invalid request: " + errMsg})
		case codes.Unauthenticated:
			c.JSON(401, gin.H{"error": "Unauthorized: " + errMsg})
		case codes.FailedPrecondition:
			c.JSON(400, gin.H{"error": "Invalid password: " + errMsg})
		default:
			c.JSON(500, gin.H{"error": "Internal server error: " + errMsg})
		}
		return
	}

	// Muvaffaqiyatli javob
	c.JSON(200, gin.H{
		"message": res.Message,
	})
}

// UpdateEmail godoc
// @Summary      Update user's email
// @Description  Allows an authenticated user to update their email
// @Tags         user
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        body  body      userpb.UpdateEmailReq  true  "New email"
// @Success      200   {object}  map[string]string
// @Failure      400   {object}  map[string]string
// @Failure      401   {object}  map[string]string
// @Failure      404   {object}  map[string]string
// @Failure      409   {object}  map[string]string
// @Failure      500   {object}  map[string]string
// @Router       /user/updateemail [put]
func (h *Handler) UpdateEmail(c *gin.Context) {
	req := &userpb.UpdateEmailReq{}

	// JSON parsing va validatsiya
	if err := c.ShouldBindJSON(req); err != nil {
		c.JSON(400, gin.H{
			"error": "Invalid request format: " + err.Error(),
		})
		return
	}

	// Tokenni olib, contextga joylash
	ctx, err := utils.InjectTokenToContext(c.Request)
	if err != nil {
		c.JSON(401, gin.H{"error": "Unauthorized: invalid or missing token"})
		return
	}

	// RPC chaqiruv
	res, err := h.service.User().UpdateEmail(ctx, req)
	if err != nil {
		grpcCode := status.Code(err)
		errMsg := status.Convert(err).Message() // gRPC error message
		switch grpcCode {
		case codes.NotFound:
			c.JSON(404, gin.H{"error": "User not found"})
		case codes.AlreadyExists:
			c.JSON(409, gin.H{"error": "Email already exists"})
		case codes.InvalidArgument:
			c.JSON(400, gin.H{"error": "Invalid request: " + errMsg})
		case codes.Unauthenticated:
			c.JSON(401, gin.H{"error": "Unauthorized: " + errMsg})
		default:
			c.JSON(500, gin.H{"error": "Internal server error: " + errMsg})
		}
		return
	}

	// Muvaffaqiyatli javob
	c.JSON(200, gin.H{
		"message": res.Message,
	})
}
