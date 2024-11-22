package routes

import (
	"github.com/go-chi/chi/v5"
	"go.mongodb.org/mongo-driver/mongo"
	"mmp/api/middleware"
	"mmp/config"
	"mmp/internal/delivery/rest/response"
	"net/http"
)

// RegisterRoutes configures the main routes for the application.
func RegisterRoutes(apiV1Router *chi.Mux, mongoDB *mongo.Client, cfg *config.AppConfig) {
	// Create a sub-router for API version 1
	apiV1Router.Mount("/api/v1", apiV1Router)

	// Create a sub-router for Loans
	apiV1WithAuthRouter := chi.NewRouter()
	apiV1WithAuthRouter.Use(middleware.BasicAuth(cfg))

	apiV1Router.Mount("/", apiV1WithAuthRouter)
	apiV1Router.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		httpResponse := response.BuildSuccessResponseWithData(response.Ok, "pong")
		response.JSON(w, httpResponse.StatusCode, httpResponse)
	})
}
