package user_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
	"user-management/app/user"
	"user-management/response"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockUsecase struct {
	mock.Mock
}

func (m *mockUsecase) CreateUser(ctx context.Context, req user.User) (*response.StdResp[any], error) {
	args := m.Called(ctx, req)
	return args.Get(0).(*response.StdResp[any]), args.Error(1)
}

func (m *mockUsecase) Login(ctx context.Context, req user.SignInRequest) (*response.StdResp[any], error) {
	args := m.Called(ctx, req)
	return args.Get(0).(*response.StdResp[any]), args.Error(1)
}

func (m *mockUsecase) FindUsers(ctx context.Context) (*response.StdResp[any], error) {
	args := m.Called(ctx)
	return args.Get(0).(*response.StdResp[any]), args.Error(1)
}

func (m *mockUsecase) FindUserById(ctx context.Context, id string) (*response.StdResp[any], error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*response.StdResp[any]), args.Error(1)
}

func (m *mockUsecase) UpdateUser(ctx context.Context, user user.User) (*response.StdResp[any], error) {
	args := m.Called(ctx, user)
	return args.Get(0).(*response.StdResp[any]), args.Error(1)
}

func (m *mockUsecase) DeleteUser(ctx context.Context, id string) (*response.StdResp[any], error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*response.StdResp[any]), args.Error(1)
}

func TestHandlerLogin(t *testing.T) {
	e := echo.New()
	reqBody := user.SignInRequest{Email: "test@example.com", Password: "123456"}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	mockUc := new(mockUsecase)
	handler := user.NewHandler(nil, mockUc)

	respData := &user.SignInResponse{
		Token:     "testtoken",
		ExpiresAt: time.Now().Add(1 * time.Hour),
	}
	stdResp := response.SuccessWithData(respData)

	mockUc.On("Login", mock.Anything, reqBody).Return(stdResp, nil)

	err := handler.Login(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), "testtoken")
}

func TestHandlerRegister(t *testing.T) {
	e := echo.New()
	reqBody := user.User{Name: "New", Email: "new@example.com", Password: "abc123"}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	mockUc := new(mockUsecase)
	handler := user.NewHandler(nil, mockUc)

	stdResp := response.SuccessWithData(map[string]string{"id": "abc123"})
	mockUc.On("CreateUser", mock.Anything, reqBody).Return(stdResp, nil)

	err := handler.CreateUser(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), "abc123")
}

func TestHandlerFindUsers(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/users", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	mockUc := new(mockUsecase)
	handler := user.NewHandler(nil, mockUc)

	mockUc.On("FindUsers", mock.Anything).Return(response.SuccessWithData([]user.User{}), nil)

	err := handler.FindUsers(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestHandlerFindUserById_Invalid(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/users/invalid", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("invalid")

	mockUc := new(mockUsecase)
	handler := user.NewHandler(nil, mockUc)

	err := handler.FindUserById(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Contains(t, rec.Body.String(), "invalid")
}

func TestHandlerFindUserById_Valid(t *testing.T) {
	e := echo.New()
	validID := "60d5ec49f1f1c939b4f2f0c2"
	req := httptest.NewRequest(http.MethodGet, "/users/"+validID, nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues(validID)

	mockUc := new(mockUsecase)
	handler := user.NewHandler(nil, mockUc)

	mockUc.On("FindUserById", mock.Anything, validID).Return(response.SuccessWithData(user.User{Name: "Test"}), nil)

	err := handler.FindUserById(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
}
