package main

import (
	"context"
	"log"
	"net"
	"user-service/config"
	"user-service/internal/service"

	db "user-service/internal/storage"

	pb "user-service/protos/user/userpb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	// Loading configuration (.env)
	config.LoadConfig()

	// Connecting to PostgreSQL
	ConnectToDB,err := db.ConnectToDB(config.AppConfig)
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}

	// Launching the Storage Layer
	postgres := db.NewPostgresStorage(ConnectToDB)

	// Opening a TCP listener for the gRPC server
	lis, err := net.Listen("tcp", config.AppConfig.Http.Host+":"+config.AppConfig.Http.Port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	// Creating a gRPC server and adding an interceptor 
	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(grpcInterceptor),
	)
	// gRPC reflection (needed for tools like grpcurl)
	reflection.Register(grpcServer)

	// Registering a UserService server
	pb.RegisterUserServiceServer(grpcServer, service.NewUserService(postgres))

	// Starting the server
	log.Println("User service running on port ", config.AppConfig.Http.Port)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

// gRPC Interceptor: to log all RPC calls
func grpcInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	m, err := handler(ctx, req)
	if err != nil {
		log.Printf("RPC failed with error: %v", err)
		return nil, err
	}
	return m, nil
}
