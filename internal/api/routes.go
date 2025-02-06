package api

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/Businge931/sba-api-gateway/internal/api/handlers"
	"github.com/Businge931/sba-api-gateway/internal/app/service"
)

func SetupRoutes(authService service.AuthService, oddsService service.OddsService) *mux.Router {
	router := mux.NewRouter()

	// Register handlers
	authHandler := handlers.NewAuthHandler(authService)
	oddsHandler := handlers.NewOddsHandler(oddsService)

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
