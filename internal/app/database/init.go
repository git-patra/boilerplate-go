package database

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"mmp/config"
	"time"
)

func InitMongoDB(cfg *config.AppConfig) (*mongo.Client, error) {
	// MongoDB connection URI
	uri := fmt.Sprintf("mongodb://%s:%s@%s:%s",
		cfg.MongoDBConfig.UserName,
		cfg.MongoDBConfig.Password,
		cfg.MongoDBConfig.Host,
		cfg.MongoDBConfig.Port,
	)

	// Create a new client with options
	clientOptions := options.Client().ApplyURI(uri)

	// Set a timeout for the connection
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Connect to MongoDB
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		logrus.Errorf("Failed to connect to MongoDB: %v", err)
		return nil, err
	}

	// Check the connection
	err = client.Ping(ctx, nil)
	if err != nil {
		logrus.Errorf("Failed to ping MongoDB: %v", err)
		return nil, err
	}

	fmt.Println("Connected to MongoDB!")
	return client, err
}
