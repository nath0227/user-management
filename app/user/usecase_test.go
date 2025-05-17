package user_test

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"log"
	"testing"
	"time"
	"user-management/app/user"
	"user-management/config"
	"user-management/response"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

type mockRepo struct {
	mock.Mock
}

func (m *mockRepo) CreateUser(ctx context.Context, user user.User) (string, error) {
	args := m.Called(ctx, user)
	return args.String(0), args.Error(1)
}

func (m *mockRepo) FindUserByEmail(ctx context.Context, email string) (user.User, error) {
	args := m.Called(ctx, email)
	return args.Get(0).(user.User), args.Error(1)
}

func (m *mockRepo) FindUserById(ctx context.Context, id string) (user.FindUserResponse, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(user.FindUserResponse), args.Error(1)
}

func (m *mockRepo) FindUsers(ctx context.Context) ([]user.FindUserResponse, error) {
	args := m.Called(ctx)
	return args.Get(0).([]user.FindUserResponse), args.Error(1)
}

func (m *mockRepo) UpdateUser(ctx context.Context, user user.User) (int64, error) {
	args := m.Called(ctx, user)
	return args.Get(0).(int64), args.Error(1)
}

func (m *mockRepo) DeleteUser(ctx context.Context, id string) (int64, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(int64), args.Error(1)
}

func newUsecaseWithMock(repo *mockRepo) user.Usecase {
	return user.NewUsecase(config.CryptoCredential{
		JwtKey:            "testsecret",
		JwtExpireDuration: time.Minute,
	}, repo)
}

func TestUsecaseCreateUser(t *testing.T) {
	repo := new(mockRepo)
	uc := newUsecaseWithMock(repo)

	input := user.CreateRequest{Name: "Test", Email: "test@example.com", Password: "pass123"}
	repo.On("CreateUser", mock.Anything, mock.MatchedBy(func(u user.User) bool {
		return u.Email == input.Email
	})).Return("abc123", nil)

	resp, err := uc.CreateUser(context.Background(), input)

	assert.NoError(t, err)
	assert.Equal(t, "abc123", resp.Data.(user.CreateResponse).Id)
	repo.AssertExpectations(t)
}

func TestUsecaseCreateUser_DuplicatedRegistration(t *testing.T) {
	repo := new(mockRepo)
	uc := newUsecaseWithMock(repo)

	input := user.CreateRequest{Name: "Test", Email: "test@example.com", Password: "pass123"}
	repo.On("CreateUser", mock.Anything, mock.MatchedBy(func(u user.User) bool {
		return u.Email == input.Email
	})).Return("", user.ErrEmailAlreadyExists)

	resp, err := uc.CreateUser(context.Background(), input)

	assert.NoError(t, err)
	assert.Equal(t, response.DuplicatedRegistration(), resp)
	repo.AssertExpectations(t)
}

func TestUsecaseLoginSuccess(t *testing.T) {
	repo := new(mockRepo)
	uc := newUsecaseWithMock(repo)

	hashed, _ := bcrypt.GenerateFromPassword([]byte("pass123"), bcrypt.DefaultCost)
	userData := user.User{Email: "test@example.com", Password: string(hashed)}

	repo.On("FindUserByEmail", mock.Anything, "test@example.com").Return(userData, nil)

	resp, err := uc.Login(context.Background(), user.SignInRequest{
		Email:    "test@example.com",
		Password: "pass123",
	})

	assert.NoError(t, err)
	token := resp.Data.(*user.SignInResponse).Token
	parsed, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		return []byte("testsecret"), nil
	})
	assert.NoError(t, err)
	assert.True(t, parsed.Valid)
	repo.AssertExpectations(t)
}

func TestUsecaseLogin_Fail(t *testing.T) {
	repo := new(mockRepo)
	uc := newUsecaseWithMock(repo)
	repo.On("FindUserByEmail", mock.Anything, "test@example.com").Return(user.User{}, user.ErrUserOrPasswordIsWrong)

	resp, err := uc.Login(context.Background(), user.SignInRequest{
		Email:    "test@example.com",
		Password: "pass123",
	})

	assert.NoError(t, err)
	assert.Equal(t, response.LoginFail(), resp)
	repo.AssertExpectations(t)
}

func TestUsecaseLoginInvalidPassword(t *testing.T) {
	repo := new(mockRepo)
	uc := newUsecaseWithMock(repo)

	hashed, _ := bcrypt.GenerateFromPassword([]byte("pass123"), bcrypt.DefaultCost)
	userData := user.User{Email: "test@example.com", Password: string(hashed)}

	repo.On("FindUserByEmail", mock.Anything, "test@example.com").Return(userData, nil)

	resp, err := uc.Login(context.Background(), user.SignInRequest{
		Email:    "test@example.com",
		Password: "wrongpassword",
	})

	assert.NoError(t, err)
	assert.Equal(t, response.LoginFail(), resp)
	repo.AssertExpectations(t)
}

func TestUsecaseFindUsers(t *testing.T) {
	repo := new(mockRepo)
	uc := newUsecaseWithMock(repo)

	users := []user.FindUserResponse{{Name: "User1"}, {Name: "User2"}}
	repo.On("FindUsers", mock.Anything).Return(users, nil)

	resp, err := uc.FindUsers(context.Background())

	assert.NoError(t, err)
	assert.Len(t, resp.Data.([]user.FindUserResponse), 2)
	repo.AssertExpectations(t)
}

func TestUsecaseFindUserById_NotFound(t *testing.T) {
	repo := new(mockRepo)
	uc := newUsecaseWithMock(repo)

	repo.On("FindUserById", mock.Anything, "notfound").Return(user.FindUserResponse{}, user.ErrUserNotFound)

	resp, err := uc.FindUserById(context.Background(), "notfound")

	assert.NoError(t, err)
	assert.Equal(t, response.UserNotFound().Code, resp.Code)
	repo.AssertExpectations(t)
}

func TestUsecaseUpdateUser(t *testing.T) {
	repo := new(mockRepo)
	uc := newUsecaseWithMock(repo)

	input := user.User{ID: primitive.NewObjectID(), Name: "Updated Name"}
	repo.On("UpdateUser", mock.Anything, input).Return(int64(1), nil)

	resp, err := uc.UpdateUser(context.Background(), input)

	assert.NoError(t, err)
	assert.Equal(t, response.Success(), resp)
	repo.AssertExpectations(t)
}

func TestUsecaseUpdateUser_NotFound(t *testing.T) {
	repo := new(mockRepo)
	uc := newUsecaseWithMock(repo)

	input := user.User{ID: primitive.NewObjectID(), Name: "Updated Name"}
	repo.On("UpdateUser", mock.Anything, input).Return(int64(0), nil)

	resp, err := uc.UpdateUser(context.Background(), input)

	assert.NoError(t, err)
	assert.Equal(t, response.UserNotFound(), resp)
	repo.AssertExpectations(t)
}

func TestUsecaseUpdateUser_DuplicatedEmail(t *testing.T) {
	repo := new(mockRepo)
	uc := newUsecaseWithMock(repo)

	input := user.User{ID: primitive.NewObjectID(), Email: "duplicate@example.com", Name: "Updated Name"}
	repo.On("UpdateUser", mock.Anything, input).Return(int64(0), user.ErrEmailAlreadyExists)

	resp, err := uc.UpdateUser(context.Background(), input)

	assert.NoError(t, err)
	assert.Equal(t, response.DuplicatedRegistration(), resp)
	repo.AssertExpectations(t)
}

func TestUsecaseDeleteUser(t *testing.T) {
	repo := new(mockRepo)
	uc := newUsecaseWithMock(repo)

	id := "abc123"
	repo.On("DeleteUser", mock.Anything, id).Return(int64(1), nil)

	resp, err := uc.DeleteUser(context.Background(), id)

	assert.NoError(t, err)
	assert.Equal(t, response.Success(), resp)
	repo.AssertExpectations(t)
}

func TestUsecaseDeleteUser_NotFound(t *testing.T) {
	repo := new(mockRepo)
	uc := newUsecaseWithMock(repo)

	id := "abc1234"
	repo.On("DeleteUser", mock.Anything, id).Return(int64(0), nil)

	resp, err := uc.DeleteUser(context.Background(), id)

	assert.NoError(t, err)
	assert.Equal(t, response.UserNotFound(), resp)
	repo.AssertExpectations(t)
}

func TestUsecaseCheckPassword(t *testing.T) {
	plain := "passwordstring"
	hashed, err := bcrypt.GenerateFromPassword([]byte(plain), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal(err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(hashed), []byte(plain))
	assert.NoError(t, err)
}

func TestUsecaseGenerateHMAC256Key(t *testing.T) {
	key := make([]byte, 32) // 256-bit key
	_, err := rand.Read(key)
	assert.NoError(t, err)

	secret := base64.StdEncoding.EncodeToString(key)
	assert.NotEmpty(t, secret)
}
