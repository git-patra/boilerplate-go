package main

import (
	"go.mongodb.org/mongo-driver/mongo"
	"mmp/config"
	"mmp/internal/app/server"

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
