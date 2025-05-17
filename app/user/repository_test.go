package user_test

import (
	"context"
	"testing"
	"user-management/app/user"
	"user-management/config"
	"user-management/storage"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
)

func TestRepository_CreateUser(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	// defer mt.Close()

	mt.Run("success", func(mt *mtest.T) {
		dbConn := storage.NewMongoConn(mt.Client, mt.Client.Database("testdb"))
		repo := user.NewRepository(dbConn, config.MongoConfig{
			Database:       "testdb",
			UserCollection: "users",
		})

		userData := user.User{
			ID:    primitive.NewObjectID(),
			Name:  "Test User",
			Email: "test@example.com",
		}

		mt.AddMockResponses(mtest.CreateSuccessResponse())
		id, err := repo.CreateUser(context.Background(), userData)

		assert.NoError(t, err)
		assert.NotEmpty(t, id)
	})
}

func TestRepository_FindUserByEmail(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("success", func(mt *mtest.T) {
		dbConn := storage.NewMongoConn(mt.Client, mt.Client.Database("testdb"))
		repo := user.NewRepository(dbConn, config.MongoConfig{
			Database:       "testdb",
			UserCollection: "users",
		})

		expectedUser := user.User{
			ID:    primitive.NewObjectID(),
			Name:  "Test User",
			Email: "test@example.com",
		}

		mt.AddMockResponses(mtest.CreateCursorResponse(1, "testdb.users", mtest.FirstBatch, bson.D{
			bson.E{Key: "_id", Value: expectedUser.ID},
			bson.E{Key: "name", Value: expectedUser.Name},
			bson.E{Key: "email", Value: expectedUser.Email},
		}))

		result, err := repo.FindUserByEmail(context.Background(), expectedUser.Email)

		assert.NoError(t, err)
		assert.Equal(t, expectedUser, result)
	})
}

func TestRepository_FindUserById(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("success", func(mt *mtest.T) {
		dbConn := storage.NewMongoConn(mt.Client, mt.Client.Database("testdb"))
		repo := user.NewRepository(dbConn, config.MongoConfig{
			Database:       "testdb",
			UserCollection: "users",
		})

		expectedUser := user.User{
			ID:    primitive.NewObjectID(),
			Name:  "Test User",
			Email: "test@example.com",
		}

		mt.AddMockResponses(mtest.CreateCursorResponse(1, "testdb.users", mtest.FirstBatch, bson.D{
			bson.E{Key: "_id", Value: expectedUser.ID},
			bson.E{Key: "name", Value: expectedUser.Name},
			bson.E{Key: "email", Value: expectedUser.Email},
		}))

		result, err := repo.FindUserById(context.Background(), expectedUser.ID.Hex())

		assert.NoError(t, err)
		assert.Equal(t, expectedUser, result)
	})
}

func TestRepository_FindUsers(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("success", func(mt *mtest.T) {
		dbConn := storage.NewMongoConn(mt.Client, mt.Client.Database("testdb"))
		repo := user.NewRepository(dbConn, config.MongoConfig{
			Database:       "testdb",
			UserCollection: "users",
		})

		expectedUsers := []user.User{
			{ID: primitive.NewObjectID(), Name: "User1", Email: "user1@example.com"},
			{ID: primitive.NewObjectID(), Name: "User2", Email: "user2@example.com"},
		}

		// Add mock responses for the cursor
		mt.AddMockResponses(
			mtest.CreateCursorResponse(1, "testdb.users", mtest.FirstBatch,
				bson.D{
					bson.E{Key: "_id", Value: expectedUsers[0].ID},
					bson.E{Key: "name", Value: expectedUsers[0].Name},
					bson.E{Key: "email", Value: expectedUsers[0].Email},
				},
			),
			mtest.CreateCursorResponse(0, "testdb.users", mtest.NextBatch,
				bson.D{
					bson.E{Key: "_id", Value: expectedUsers[1].ID},
					bson.E{Key: "name", Value: expectedUsers[1].Name},
					bson.E{Key: "email", Value: expectedUsers[1].Email},
				},
			),
		)

		result, err := repo.FindUsers(context.Background())

		assert.NoError(t, err)
		assert.Equal(t, expectedUsers, result)
	})
}

func TestRepository_UpdateUser(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("success", func(mt *mtest.T) {
		dbConn := storage.NewMongoConn(mt.Client, mt.Client.Database("testdb"))
		repo := user.NewRepository(dbConn, config.MongoConfig{
			Database:       "testdb",
			UserCollection: "users",
		})

		userID := primitive.NewObjectID()
		updateUser := user.User{
			ID:    userID,
			Name:  "Updated Name",
			Email: "updated@example.com",
		}

		// Mock a successful update response with MatchedCount = 1 and ModifiedCount = 1
		mt.AddMockResponses(bson.D{
			bson.E{Key: "ok", Value: 1},
			bson.E{Key: "n", Value: 1}, // MatchedCount
		})

		count, err := repo.UpdateUser(context.Background(), updateUser)

		assert.NoError(t, err)
		assert.Equal(t, int64(1), count)
	})

	mt.Run("duplicate email", func(mt *mtest.T) {
		dbConn := storage.NewMongoConn(mt.Client, mt.Client.Database("testdb"))
		repo := user.NewRepository(dbConn, config.MongoConfig{
			Database:       "testdb",
			UserCollection: "users",
		})

		userID := primitive.NewObjectID()
		updateUser := user.User{
			ID:    userID,
			Name:  "Updated Name",
			Email: "duplicate@example.com",
		}

		mt.AddMockResponses(mtest.CreateWriteErrorsResponse(mtest.WriteError{
			Index:   0,
			Code:    11000, // Duplicate key error
			Message: "duplicate key error",
		}))
		count, err := repo.UpdateUser(context.Background(), updateUser)

		assert.Error(t, err)
		assert.Equal(t, int64(0), count)
		assert.Equal(t, user.EmailAlreadyExists, err.Error())
	})
}

func TestRepository_DeleteUser(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("success", func(mt *mtest.T) {
		dbConn := storage.NewMongoConn(mt.Client, mt.Client.Database("testdb"))
		repo := user.NewRepository(dbConn, config.MongoConfig{
			Database:       "testdb",
			UserCollection: "users",
		})

		userID := primitive.NewObjectID()

		mt.AddMockResponses(bson.D{
			bson.E{Key: "ok", Value: 1},
			bson.E{Key: "n", Value: 1}, // MatchedCount
		})
		count, err := repo.DeleteUser(context.Background(), userID.Hex())

		assert.NoError(t, err)
		assert.Equal(t, int64(1), count)
	})

	mt.Run("invalid ID", func(mt *mtest.T) {
		dbConn := storage.NewMongoConn(mt.Client, mt.Client.Database("testdb"))
		repo := user.NewRepository(dbConn, config.MongoConfig{
			Database:       "testdb",
			UserCollection: "users",
		})

		invalidID := "invalidObjectID"

		count, err := repo.DeleteUser(context.Background(), invalidID)

		assert.Error(t, err)
		assert.Equal(t, int64(0), count)
	})
}

func TestRepository_CountUsers(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("success", func(mt *mtest.T) {
		dbConn := storage.NewMongoConn(mt.Client, mt.Client.Database("testdb"))
		repo := user.NewRepository(dbConn, config.MongoConfig{
			Database:       "testdb",
			UserCollection: "users",
		})

		// Mock a successful count response
		mt.AddMockResponses(mtest.CreateCursorResponse(
			1,
			"test.users",
			mtest.FirstBatch,
			bson.D{
				bson.E{"n", int64(3)},
			},
		))

		count, err := repo.CountUsers(context.Background())

		assert.NoError(t, err)
		assert.Equal(t, int64(3), count)
	})

	mt.Run("error", func(mt *mtest.T) {
		dbConn := storage.NewMongoConn(mt.Client, mt.Client.Database("testdb"))
		repo := user.NewRepository(dbConn, config.MongoConfig{
			Database:       "testdb",
			UserCollection: "users",
		})

		// Mock an error response
		mt.AddMockResponses(bson.D{
			{Key: "ok", Value: 0},
			{Key: "errmsg", Value: "count error"},
		})

		count, err := repo.CountUsers(context.Background())

		assert.Error(t, err)
		assert.Equal(t, int64(0), count)
	})
}
