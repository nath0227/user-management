package storage

import (
	"context"
	"fmt"
	"log"
	"user-management/config"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DatabaseConn interface {
	GetDatabaseCollection(databaseName, collectionName string) *mongo.Collection
	Disconnect(ctx context.Context)
}

type MongoConn struct {
	Client *mongo.Client
}

func NewMongoConnection(ctx context.Context, cfg config.MongoConfig) *MongoConn {
	uri := fmt.Sprintf(cfg.Uri, cfg.Username, cfg.Password, cfg.Database)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatal(err)
	}
	if err := client.Ping(ctx, nil); err != nil {
		log.Fatal(err)
	}

	return &MongoConn{
		Client: client,
	}
}

func (m *MongoConn) Disconnect(ctx context.Context) {
	if err := m.Client.Disconnect(ctx); err != nil {
		log.Println("Error disconnecting MongoDB:", err)
	} else {
		log.Println("MongoDB connection closed.")
	}
}

func (m *MongoConn) GetDatabaseCollection(databaseName, collectionName string) *mongo.Collection {
	return m.Client.Database(databaseName).Collection(collectionName)
}
