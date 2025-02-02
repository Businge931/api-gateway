package main

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	log "github.com/sirupsen/logrus"
)

func main() {
	// Connect to gRPC services
	oddsConn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Printf("Failed to connect to odds service: %v", err)
	}
	defer oddsConn.Close()

	authConn, err := grpc.NewClient("localhost:50052", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Printf("Failed to connect to auth service: %v", err)
	}
	defer authConn.Close()

	// Setup routes
	r := setupRoutes(oddsConn, authConn)

	// Register middleware
	r.Use(loggingMiddleware)

	// Add CORS middleware
	corsHandler := getCORSConfig()
	handler := corsHandler.Handler(r)

	// Start the server
	start(handler, ":8080")
}
