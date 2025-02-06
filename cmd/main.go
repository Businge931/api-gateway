package main

import (
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/Businge931/sba-api-gateway/internal/api"
	"github.com/Businge931/sba-api-gateway/internal/api/middleware"
	"github.com/Businge931/sba-api-gateway/internal/app/service"
	
	gRPC "github.com/Businge931/sba-api-gateway/internal/client/grpc"
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

	// Initialize gRPC clients
	authClient := gRPC.NewAuthClient(authConn)
	oddsClient := gRPC.NewOddsClient(oddsConn)

	// Initialize services
	authService := service.NewAuthService(authClient)
	oddsService := service.NewOddsService(oddsClient)

	// Setup routes
	router := api.SetupRoutes(authService, oddsService)

	// Register middleware
	router.Use(middleware.LoggingMiddleware)

	// Add CORS middleware
	corsHandler := middleware.GetCORSConfig()
	handler := corsHandler.Handler(router)

	// Start the server
	api.Start(handler, ":8080")
}
