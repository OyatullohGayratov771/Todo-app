package main

import (
	"context"
	"log"
	"net"
	"notification-service/config"
	"notification-service/kafka"
	notifactionpb "notification-service/protos/notification"
	"notification-service/service"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	config.LoadConfig()
	log.Println(config.AppConfig)

	consumer := kafka.ConnKafka(config.AppConfig)
	defer consumer.Close()

	service := service.NewNotificationServer(consumer, config.AppConfig.Kafka.Topic)
	go service.ListenForKafkaMessages()

	grpcPort := config.AppConfig.Http.Port
	grpcHost := config.AppConfig.Http.Host

	lis, err := net.Listen("tcp", grpcHost+":"+grpcPort)
	if err != nil {
		log.Fatalf("Failed to listen on %s:%s, error: %v", grpcHost, grpcPort, err)
	}

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(grpcInterceptor),
	)

	reflection.Register(grpcServer)

	notifactionpb.RegisterNotificationServiceServer(grpcServer, service)

	log.Println("🚀 gRPC Notification service running on", grpcHost+":"+grpcPort)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatal("Failed to start gRPC server:", err)
	}
}

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
