package server

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	customMiddleware "mmp/api/middleware"
	"mmp/api/routes"
	"mmp/config"
	"net/http"
)

type Server struct {
	chiRouter *chi.Mux
	config    *config.AppConfig
}

func NewServer(
	cfg *config.AppConfig,
	mongoDB *mongo.Client,
) *Server {

	return &Server{
		chiRouter: initializeChiRouter(
			cfg,
			mongoDB,
		),
		config: cfg,
	}
}

func (server *Server) Start() error {
	httpServer := http.Server{
		Addr:        fmt.Sprintf(":%d", server.config.ServerPort),
		Handler:     server.chiRouter,
		ReadTimeout: server.config.RequestTimeout,
	}

	err := ServeHTTP(&httpServer, httpServer.Addr, 0)
	if err != nil {
		logrus.Error("failed to start the REST API server:", err)
		return err
	}

	logrus.Info("REST API server stopped")
	return nil
}

func initializeChiRouter(
	config *config.AppConfig,
	mongoDB *mongo.Client,
) *chi.Mux {
	chiRouter := chi.NewRouter()

	// Middlewares.
	chiRouter.Use(middleware.Recoverer)
	chiRouter.Use(customMiddleware.CORSMiddleware)

	// Routes.
	routes.RegisterRoutes(chiRouter, mongoDB, config)

	return chiRouter
}
