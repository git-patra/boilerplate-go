package main

import (
	"boilerplate-go/config"
	"boilerplate-go/internal/app/server"
	"go.mongodb.org/mongo-driver/mongo"

	"log"
)

func MainHttpHandler(cfg *config.AppConfig, mongoDB *mongo.Client) {
	serverHttp := server.NewServer(cfg, mongoDB)
	err := serverHttp.Start()

	if err != nil {
		log.Fatalf("error starting server: %s", err)
		return
	}

}
