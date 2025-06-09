package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"strconv"

	pb "github.com/e6a5/learning/backend/04-grpc-basics/proto"
	"github.com/e6a5/learning/backend/04-grpc-basics/server/internal/repository"
	"github.com/e6a5/learning/backend/04-grpc-basics/server/internal/service"
	"google.golang.org/grpc"
)

func main() {
	// Get port from environment or use default
	port, err := getPort()
	if err != nil {
		log.Fatalf("Invalid port configuration: %v", err)
	}

	// Initialize dependencies
	userRepo := repository.NewUserRepository()
	userService := service.NewUserService(userRepo)

	// Setup gRPC server
	grpcServer, listener, err := setupGRPCServer(port, userService)
	if err != nil {
		log.Fatalf("Failed to setup gRPC server: %v", err)
	}

	logServerInfo(port)

	// Start serving
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}

func getPort() (int, error) {
	portStr := os.Getenv("GRPC_PORT")
	if portStr == "" {
		return 50051, nil
	}

	port, err := strconv.Atoi(portStr)
	if err != nil {
		return 0, fmt.Errorf("invalid GRPC_PORT: %w", err)
	}

	if port < 1 || port > 65535 {
		return 0, fmt.Errorf("port must be between 1 and 65535, got %d", port)
	}

	return port, nil
}

func setupGRPCServer(port int, userService *service.UserService) (*grpc.Server, net.Listener, error) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, nil, fmt.Errorf("failed to listen on port %d: %w", port, err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterUserServiceServer(grpcServer, userService)

	return grpcServer, listener, nil
}

func logServerInfo(port int) {
	log.Printf("🚀 gRPC Server running on port %d", port)
	log.Println("Available services:")
	log.Println("  - CreateUser (unary)")
	log.Println("  - GetUser (unary)")
	log.Println("  - ListUsers (unary)")
	log.Println("  - WatchUsers (server streaming)")
	log.Println("  - BatchCreateUsers (client streaming)")
}
