package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	indexgrpc "bleve-proj/internal/grpc"
	proto "bleve-proj/proto"

	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found. Using default values.")
	}

	// Get values from environment or use defaults
	indexDir := getEnv("INDEX_DIRECTORY", "./index")
	grpcPort := getEnv("GRPC_PORT", ":50051")

	// Start gRPC server
	server := startGRPCServer(indexDir, grpcPort)

	// Graceful shutdown handling
	gracefulShutdown(server)
}

// startGRPCServer starts the gRPC server
func startGRPCServer(indexDir, grpcPort string) *grpc.Server {
	lis, err := net.Listen("tcp", grpcPort)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	s := grpc.NewServer()
	proto.RegisterIndexServiceServer(s, indexgrpc.NewServer(indexDir))
	reflection.Register(s)

	fmt.Printf("Starting gRPC server on port %s with index directory: %s\n", grpcPort, indexDir)
	go func() {
		if err := s.Serve(lis); err != nil {
			log.Fatalf("Failed to serve: %v", err)
		}
	}()
	return s
}

// gracefulShutdown handles graceful shutdown of the server
func gracefulShutdown(server *grpc.Server) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	fmt.Println("Shutting down gRPC server...")

	// Give the server 10 seconds to gracefully stop
	_, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	server.GracefulStop()

	fmt.Println("gRPC server stopped.")
	// Optional: Close any resources here.
}

// getEnv retrieves an environment variable or returns a default value
func getEnv(key string, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		log.Printf("Environment variable %s not set. Using default value: %s\n", key, defaultValue)
		return defaultValue
	}
	return value
}
