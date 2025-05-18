package storage

import (
	"context"
	"fmt"
	"user-management/config"
	"user-management/logger"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DatabaseConn interface {
	Disconnect(ctx context.Context)
	Collection(name string) CollectionInterface
}

type CollectionInterface interface {
	Find(ctx context.Context, filter interface{}, opts ...*options.FindOptions) (*mongo.Cursor, error)
	FindOne(ctx context.Context, filter interface{}, opts ...*options.FindOneOptions) *mongo.SingleResult
	InsertOne(ctx context.Context, document interface{}, opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error)
	UpdateByID(ctx context.Context, id interface{}, update interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error)
	DeleteOne(ctx context.Context, filter interface{}, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error)
	CountDocuments(ctx context.Context, filter interface{}, opts ...*options.CountOptions) (int64, error)
}

type MongoConn struct {
	client   *mongo.Client
	database *mongo.Database
}

func NewMongoConn(client *mongo.Client, database *mongo.Database) *MongoConn {
	return &MongoConn{
		client:   client,
		database: database,
	}
}

func InitMongoConnection(ctx context.Context, cfg config.MongoConfig) *MongoConn {
	zlog := logger.NewZap()

	uri := fmt.Sprintf(cfg.Uri, cfg.Username, cfg.Password, cfg.Database)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		zlog.Fatal(err.Error())
	}
	if err := client.Ping(ctx, nil); err != nil {
		zlog.Fatal(err.Error())
	}

	return NewMongoConn(client, client.Database(cfg.Database))
}

func (m *MongoConn) Disconnect(ctx context.Context) {
	zlog := logger.NewZap()
	if err := m.client.Disconnect(ctx); err != nil {
		zlog.Sugar().Errorf("Error disconnecting MongoDB:", err)
	} else {
		zlog.Info("MongoDB connection closed.")
	}
}

func (m *MongoConn) Collection(name string) CollectionInterface {
	return &MongoCollection{coll: m.database.Collection(name)}
}

type MongoCollection struct {
	coll *mongo.Collection
}

func (c *MongoCollection) Find(ctx context.Context, filter interface{}, opts ...*options.FindOptions) (*mongo.Cursor, error) {
	return c.coll.Find(ctx, filter, opts...)
}

func (c *MongoCollection) FindOne(ctx context.Context, filter interface{}, opts ...*options.FindOneOptions) *mongo.SingleResult {
	return c.coll.FindOne(ctx, filter, opts...)
}

func (c *MongoCollection) InsertOne(ctx context.Context, document interface{}, opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error) {
	return c.coll.InsertOne(ctx, document, opts...)
}

func (c *MongoCollection) UpdateByID(ctx context.Context, id interface{}, update interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	return c.coll.UpdateByID(ctx, id, update, opts...)
}

func (c *MongoCollection) DeleteOne(ctx context.Context, filter interface{}, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error) {
	return c.coll.DeleteOne(ctx, filter, opts...)
}

func (c *MongoCollection) CountDocuments(ctx context.Context, filter interface{}, opts ...*options.CountOptions) (int64, error) {
	return c.coll.CountDocuments(ctx, filter, opts...)
}
