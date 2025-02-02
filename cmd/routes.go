package main

import (
	"net/http"

	"github.com/gorilla/mux"
	"google.golang.org/grpc"

	"github.com/Businge931/sba-api-gateway/handlers"
)

func setupRoutes(oddsConn, authConn *grpc.ClientConn) *mux.Router {
	router := mux.NewRouter()

	// Register handlers
	oddsHandler := handlers.NewOddsHandler(oddsConn)
	authHandler := handlers.NewAuthHandler(authConn)

	// Public endpoints
	router.HandleFunc("/login", authHandler.Login).Methods("POST")
	router.HandleFunc("/register", authHandler.Register).Methods("POST")

	// Protected endpoints (require token verification)
	router.Handle("/create", authHandler.VerifyTokenMiddleware(http.HandlerFunc(oddsHandler.CreateOdds))).Methods("POST")
	router.Handle("/read", authHandler.VerifyTokenMiddleware(http.HandlerFunc(oddsHandler.ReadOdds))).Methods("GET")
	router.Handle("/update", authHandler.VerifyTokenMiddleware(http.HandlerFunc(oddsHandler.UpdateOdds))).Methods("PUT")
	router.Handle("/delete", authHandler.VerifyTokenMiddleware(http.HandlerFunc(oddsHandler.DeleteOdds))).Methods("DELETE")

	return router
}
