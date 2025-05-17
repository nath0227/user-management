package user_test

import (
	"testing"
	"user-management/app/user"
	"user-management/response"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestUser_RequestValidation(t *testing.T) {
	tests := []struct {
		name     string
		input    user.CreateRequest
		wantCode *response.StdResp[any]
	}{
		{"Valid input 1", user.CreateRequest{Name: "Alice", Email: "alice@example.com", Password: "pass123"}, response.Success()},
		{"Missing name", user.CreateRequest{Email: "alice@example.com", Password: "pass123"}, response.MandatoryMissing("name")},
		{"Valid input 1", user.CreateRequest{Name: "Alice", Password: "pass123"}, response.MandatoryMissing("email")},
		{"Invalid email", user.CreateRequest{Name: "Alice", Email: "bad-email", Password: "pass123"}, response.InvalidData("email")},
		{"Missing password", user.CreateRequest{Name: "Alice", Email: "alice@example.com"}, response.MandatoryMissing("password")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp := tt.input.RequestValidation()
			assert.Equal(t, tt.wantCode, resp)
		})
	}
}

func TestSignInRequest_RequestValidation(t *testing.T) {
	tests := []struct {
		name     string
		input    user.SignInRequest
		wantCode string
	}{
		{"Valid", user.SignInRequest{Email: "user@test.com", Password: "123"}, "0000"},
		{"Missing email", user.SignInRequest{Password: "123"}, "4001"},
		{"Invalid email", user.SignInRequest{Email: "bad@", Password: "123"}, "4004"},
		{"Missing password", user.SignInRequest{Email: "user@test.com"}, "4001"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp := tt.input.RequestValidation()
			assert.Equal(t, tt.wantCode, resp.Code)
		})
	}
}

func TestUpdateRequest_RequestValidation(t *testing.T) {
	tests := []struct {
		name     string
		input    user.UpdateRequest
		wantCode string
	}{
		{"Valid email", user.UpdateRequest{Email: "test@test.com"}, "0000"},
		{"Valid name", user.UpdateRequest{Name: "Bob"}, "0000"},
		{"Invalid email", user.UpdateRequest{Email: "bad"}, "4004"},
		{"Missing name and email", user.UpdateRequest{}, "4001"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp := tt.input.RequestValidation()
			assert.Equal(t, tt.wantCode, resp.Code)
		})
	}
}

func TestIdValidation(t *testing.T) {
	validID := primitive.NewObjectID().Hex()
	invalidID := "not_a_valid_id"

	assert.Equal(t, "0000", user.IdValidation(validID).Code)
	assert.Equal(t, "4003", user.IdValidation(invalidID).Code)
}
