package main

import (
	"github.com/sirupsen/logrus"
	"log"
	"mmp/config"
	"mmp/internal/app/database"
	"sync"
)

func main() {
	configFilePath := ".env.yaml"
	cfg, err := config.LoadConfig(configFilePath)
	if err != nil {
		log.Fatalf("error loading config: %s", err)
		return
	}

	mongoDB, err := database.InitMongoDB(cfg)
	if err != nil {
		panic(err)
	}

	wg := sync.WaitGroup{}

	wg.Add(1)
	go func() {
		defer wg.Done()
		logrus.Info("Starting HTTP handler")
		MainHttpHandler(cfg, mongoDB)
	}()

	wg.Wait()
}
