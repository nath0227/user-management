package user

import (
	"context"
	"errors"
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
	ior, err := collectionUser(r).InsertOne(ctx, user)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return "", errors.New(EmailAlreadyExists)
		}
		return "", err
	}
	return ior.InsertedID.(primitive.ObjectID).Hex(), err
}

func (r *repository) FindUserByEmail(ctx context.Context, email string) (User, error) {
	var user User
	err := collectionUser(r).FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return user, errors.New(UserOrPasswordIsWrong)
		}
		return user, err
	}
	return user, err
}

func (r *repository) FindUserById(ctx context.Context, id string) (User, error) {
	var user User
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return user, err
	}
	err = collectionUser(r).FindOne(ctx, bson.M{"_id": oid}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return user, errors.New(UserNotFound)
		}
		return user, err
	}
	return user, err
}

func (r *repository) FindUsers(ctx context.Context) ([]User, error) {
	sort := bson.D{{"created_at", 1}}
	opts := options.Find().SetSort(sort)
	cursor, err := collectionUser(r).Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, err
	}
	var users []User
	err = cursor.All(ctx, &users)
	return users, err
}

func (r *repository) UpdateUser(ctx context.Context, user User) error {
	updateFields := bson.M{}
	if user.Name != "" {
		updateFields["name"] = user.Name
	}
	if user.Email != "" {
		updateFields["email"] = user.Email
	}
	update := bson.M{"$set": updateFields}
	_, err := collectionUser(r).UpdateByID(ctx, user.ID, update)
	if mongo.IsDuplicateKeyError(err) {
		return errors.New(EmailAlreadyExists)
	}
	return err
}

func (r *repository) DeleteUser(ctx context.Context, id string) (int64, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return 0, err
	}
	result, err := collectionUser(r).DeleteOne(ctx, bson.M{"_id": oid})
	return result.DeletedCount, err
}

func (r *repository) CountUsers(ctx context.Context) (int64, error) {
	return collectionUser(r).CountDocuments(ctx, bson.M{})
}

func collectionUser(r *repository) *mongo.Collection {
	return r.mc.GetDatabaseCollection(r.cfg.Database, r.cfg.UserCollection)
}
