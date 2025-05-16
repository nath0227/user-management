package user_test

import (
	"context"
	"errors"
	"testing"
	"user-management/app/user"
	usergrpc "user-management/app/user/grpc/gen/go/user/v1"
	"user-management/response"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestGrpcHandler_CreateUser(t *testing.T) {
	mockUc := new(mockUsecase)
	handler := user.NewGrpcHandler(mockUc)

	ctx := context.Background()
	req := &usergrpc.CreateUserRequest{
		Name:     "John Doe",
		Email:    "john.doe@example.com",
		Password: "password123",
	}

	mockUc.On("CreateUser", ctx, mock.AnythingOfType("user.User")).Return(
		response.SuccessWithData(user.CreateResponse{Id: "12345"}), nil)

	resp, err := handler.CreateUser(ctx, req)
	assert.NoError(t, err)
	assert.Equal(t, "0000", resp.Code)
	assert.Equal(t, "Success", resp.Message)
	assert.Equal(t, "12345", resp.Data.Id)

	mockUc.AssertExpectations(t)
}

func TestGrpcHandler_CreateUser_Duplicated(t *testing.T) {
	mockUc := new(mockUsecase)
	handler := user.NewGrpcHandler(mockUc)

	ctx := context.Background()
	req := &usergrpc.CreateUserRequest{
		Name:     "John Doe",
		Email:    "john.doe@example.com",
		Password: "password123",
	}

	mockUc.On("CreateUser", ctx, mock.AnythingOfType("user.User")).Return(response.DuplicatedRegistration(), errors.New(user.EmailAlreadyExists))

	resp, err := handler.CreateUser(ctx, req)
	assert.NoError(t, err)
	assert.Equal(t, response.DuplicatedRegistration().Code, resp.Code)
	assert.Equal(t, response.DuplicatedRegistration().Message, resp.Message)

	mockUc.AssertExpectations(t)
}

func TestGrpcHandler_CreateUser_Error(t *testing.T) {
	mockUc := new(mockUsecase)
	handler := user.NewGrpcHandler(mockUc)

	ctx := context.Background()
	req := &usergrpc.CreateUserRequest{
		Name:     "John Doe",
		Email:    "john.doe@example.com",
		Password: "password123",
	}

	mockUc.On("CreateUser", ctx, mock.AnythingOfType("user.User")).Return(&response.StdResp[any]{}, errors.New("internal error"))

	resp, err := handler.CreateUser(ctx, req)
	assert.Error(t, err)
	assert.Nil(t, resp)

	mockUc.AssertExpectations(t)
}

func TestGrpcHandler_GetUser(t *testing.T) {
	mockUc := new(mockUsecase)
	handler := user.NewGrpcHandler(mockUc)

	ctx := context.Background()
	oid := primitive.NewObjectID()
	req := &usergrpc.GetUserRequest{Id: oid.Hex()}
	mockUc.On("FindUserById", ctx, oid.Hex()).Return(response.SuccessWithData(user.User{
		ID:    oid,
		Name:  "John Doe",
		Email: "john.doe@example.com",
	}), nil)

	resp, err := handler.GetUser(ctx, req)
	assert.NoError(t, err)
	assert.Equal(t, "0000", resp.Code)
	assert.Equal(t, "Success", resp.Message)
	assert.Equal(t, oid.Hex(), resp.Data.Id)
	assert.Equal(t, "John Doe", resp.Data.Name)
	assert.Equal(t, "john.doe@example.com", resp.Data.Email)

	mockUc.AssertExpectations(t)
}

func TestGrpcHandler_GetUser_InvalidId(t *testing.T) {
	mockUc := new(mockUsecase)
	handler := user.NewGrpcHandler(mockUc)

	ctx := context.Background()
	req := &usergrpc.GetUserRequest{Id: "12345"}
	
	resp, err := handler.GetUser(ctx, req)
	assert.NoError(t, err)
	assert.Equal(t, response.InvalidData("id").Code, resp.Code)
	assert.Equal(t, response.InvalidData("id").Message, resp.Message)
	
	mockUc.AssertExpectations(t)
}
func TestGrpcHandler_GetUser_NotFound(t *testing.T) {
	mockUc := new(mockUsecase)
	handler := user.NewGrpcHandler(mockUc)

	ctx := context.Background()
	oid := primitive.NewObjectID()
	req := &usergrpc.GetUserRequest{Id: oid.Hex()}
	mockUc.On("FindUserById", ctx, oid.Hex()).Return(response.UserNotFound(), nil)

	resp, err := handler.GetUser(ctx, req)
	assert.NoError(t, err)
	assert.Equal(t, response.UserNotFound().Code, resp.Code)
	assert.Equal(t, response.UserNotFound().Message, resp.Message)
	
	mockUc.AssertExpectations(t)
}
