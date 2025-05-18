package main

import (
	api "api-gateway/api/v1"
	"api-gateway/config"
	_ "api-gateway/docs"
	"api-gateway/service"
	"flag"
	"log/slog"
	"net/http"
	"os"

	ginSwagger "github.com/swaggo/gin-swagger"
	swaggerFiles "github.com/swaggo/files"
)

// @title           TODO-APP
// @version         1.0
// @description     This is a sample server for todo-app.
// @host            localhost:8080
// @BasePath        /
// @schemes         http
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	config.LoadConfig()
	cfg := config.AppConfig

	addr := flag.String("addr", cfg.Http.Host+":"+cfg.Http.Port, "HTTP Server address")
	logger := slog.New(slog.NewJSONHandler(os.Stdout,
		&slog.HandlerOptions{
			AddSource: true,
		}))

	flag.Parse()

	s, err := service.New(cfg)
	if err != nil {
		logger.Error("Failed to initialize services", "error", err)
		os.Exit(1)
	}

	router := api.Router(api.Options{
		Service: s,
	})

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	srv := &http.Server{
		Addr:    *addr,
		Handler: router,
	}

	logger.Info("Starting server", "addr", srv.Addr)
	if err := srv.ListenAndServe(); err != nil {
		logger.Error("Failed to start server", "error", err)
		os.Exit(1)
	}

}
