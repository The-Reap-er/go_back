package database

import (
	"context"
	"log"
	"time"

	"go_back/config"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	Client              *mongo.Client
	UserCollection      *mongo.Collection
	URLCollection       *mongo.Collection
	ScanLogCollection   *mongo.Collection // Collection for scan logs
	APITargetCollection *mongo.Collection // Collection for API targets
	APIAlertCollection  *mongo.Collection // Collection for API alerts
)

func ConnectDB(cfg config.Config) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clientOptions := options.Client().ApplyURI(cfg.MongoURI)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal("Failed to connect to MongoDB:", err)
	}

	// Ping the database to verify connection
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal("Failed to ping MongoDB:", err)
	}

	Client = client
	db := client.Database(cfg.MongoDB)

	UserCollection = db.Collection("users")
	URLCollection = db.Collection("urls")
	ScanLogCollection = db.Collection("scan_logs")
	APITargetCollection = db.Collection("api_targets")
	APIAlertCollection = db.Collection("api_alerts") // Initialize the new collection

	// Create indexes if necessary
	_, err = UserCollection.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bson.M{"username": 1},
		Options: options.Index().SetUnique(true),
	})
	if err != nil {
		log.Fatal("Failed to create index on users collection:", err)
	}
}
