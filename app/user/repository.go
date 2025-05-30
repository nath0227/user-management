package user

import (
	"context"
	"time"
	"user-management/config"
	"user-management/storage"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type CountUsersRepository interface {
	CountUsers(ctx context.Context) (int64, error)
}

type repository struct {
	mc  storage.DatabaseConn
	cfg config.MongoConfig
}

func NewRepository(mc storage.DatabaseConn, cfg config.MongoConfig) *repository {
	return &repository{
		mc:  mc,
		cfg: cfg,
	}
}

func (r *repository) CreateUser(ctx context.Context, user User) (string, error) {
	user.CreatedAt = time.Now()
	ior, err := r.mc.Collection(r.cfg.UserCollection).InsertOne(ctx, user)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return "", ErrEmailAlreadyExists
		}
		return "", err
	}
	return ior.InsertedID.(primitive.ObjectID).Hex(), err
}

func (r *repository) FindUserByEmail(ctx context.Context, email string) (User, error) {
	var user User
	err := r.mc.Collection(r.cfg.UserCollection).FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return user, ErrUserOrPasswordIsWrong
		}
		return user, err
	}
	return user, err
}

func (r *repository) FindUserById(ctx context.Context, id string) (FindUserResponse, error) {
	var user FindUserResponse
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return user, err
	}
	err = r.mc.Collection(r.cfg.UserCollection).FindOne(ctx, bson.M{"_id": oid}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return user, ErrUserNotFound
		}
		return user, err
	}
	return user, err
}

func (r *repository) FindUsers(ctx context.Context) ([]FindUserResponse, error) {
	sort := bson.D{bson.E{Key: "created_at", Value: 1}}
	opts := options.Find().SetSort(sort)
	cursor, err := r.mc.Collection(r.cfg.UserCollection).Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, err
	}
	var users []FindUserResponse
	err = cursor.All(ctx, &users)
	return users, err
}

func (r *repository) UpdateUser(ctx context.Context, user User) (int64, error) {
	updateFields := bson.M{}
	if user.Name != "" {
		updateFields["name"] = user.Name
	}
	if user.Email != "" {
		updateFields["email"] = user.Email
	}
	update := bson.M{"$set": updateFields}
	result, err := r.mc.Collection(r.cfg.UserCollection).UpdateByID(ctx, user.ID, update)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return 0, ErrEmailAlreadyExists
		}
		return 0, err
	}
	return result.MatchedCount, nil
}

func (r *repository) DeleteUser(ctx context.Context, id string) (int64, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return 0, err
	}
	result, err := r.mc.Collection(r.cfg.UserCollection).DeleteOne(ctx, bson.M{"_id": oid})
	return result.DeletedCount, err
}

func (r *repository) CountUsers(ctx context.Context) (int64, error) {
	return r.mc.Collection(r.cfg.UserCollection).CountDocuments(ctx, bson.M{})
}
