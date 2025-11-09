package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ConnectMongoDB(uri, dbName string) (*mongo.Database, *mongo.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clientOpts := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(ctx, clientOpts)
	if err != nil {
		log.Fatal("failed to connect to MongoDB:", err)
		return nil, nil, err
	}

	if err := client.Ping(ctx, nil); err != nil {
		log.Fatal("MongoDB ping failed:", err)
		return nil, nil, err
	}

	fmt.Println("Connected to MongoDB...")
	db := client.Database(dbName)
	return db, client, nil
}
