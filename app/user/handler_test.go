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
	"user-management/logger"
	"user-management/middleware"
	"user-management/response"

	echo "github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type mockUsecase struct {
	mock.Mock
}

func (m *mockUsecase) CreateUser(ctx context.Context, req user.CreateRequest) (*response.StdResp[any], error) {
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
	e.Use(middleware.NewLogging)
	e.Use(middleware.LoggingMiddleware)
	reqBody := user.SignInRequest{Email: "test@example.com", Password: "123456"}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	ctx := context.WithValue(c.Request().Context(), logger.LogContext, logger.NewZap())
	c.SetRequest(req.WithContext(ctx))

	mockUc := new(mockUsecase)
	handler := user.NewHandler(mockUc)

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

func TestHandlerLogin_InvalidRequest(t *testing.T) {
	e := echo.New()
	e.Use(middleware.NewLogging)
	e.Use(middleware.LoggingMiddleware)
	reqBody := user.SignInRequest{Email: "test.example.com", Password: "123456"}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	ctx := context.WithValue(c.Request().Context(), logger.LogContext, logger.NewZap())
	c.SetRequest(req.WithContext(ctx))

	mockUc := new(mockUsecase)
	handler := user.NewHandler(mockUc)

	mockUc.On("Login", mock.Anything, reqBody).Return(response.InvalidData("email"), nil)

	err := handler.Login(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Contains(t, rec.Body.String(), response.InvalidData("email").Message)
}

func TestHandlerLogin_Fail(t *testing.T) {
	e := echo.New()
	e.Use(middleware.NewLogging)
	e.Use(middleware.LoggingMiddleware)
	reqBody := user.SignInRequest{Email: "test@example.com", Password: "123456"}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	ctx := context.WithValue(c.Request().Context(), logger.LogContext, logger.NewZap())
	c.SetRequest(req.WithContext(ctx))

	mockUc := new(mockUsecase)
	handler := user.NewHandler(mockUc)

	mockUc.On("Login", mock.Anything, reqBody).Return(response.LoginFail(), nil)

	err := handler.Login(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Contains(t, rec.Body.String(), response.LoginFail().Message)
}

func TestHandlerRegister(t *testing.T) {
	e := echo.New()
	e.Use(middleware.NewLogging)
	e.Use(middleware.LoggingMiddleware)
	reqBody := user.CreateRequest{Name: "New", Email: "new@example.com", Password: "abc123"}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	ctx := context.WithValue(c.Request().Context(), logger.LogContext, logger.NewZap())
	c.SetRequest(req.WithContext(ctx))

	mockUc := new(mockUsecase)
	handler := user.NewHandler(mockUc)

	stdResp := response.SuccessWithData(user.CreateResponse{Id: "abc123"})
	mockUc.On("CreateUser", mock.Anything, reqBody).Return(stdResp, nil)

	err := handler.CreateUser(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), "abc123")
}

func TestHandlerRegister_InvalidRequest(t *testing.T) {
	e := echo.New()
	e.Use(middleware.NewLogging)
	e.Use(middleware.LoggingMiddleware)
	reqBody := user.CreateRequest{Name: "", Email: "new@example.com", Password: "abc123"}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	ctx := context.WithValue(c.Request().Context(), logger.LogContext, logger.NewZap())
	c.SetRequest(req.WithContext(ctx))

	mockUc := new(mockUsecase)
	handler := user.NewHandler(mockUc)

	mockUc.On("CreateUser", mock.Anything, reqBody).Return(response.MandatoryMissing("name"), nil)

	err := handler.CreateUser(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Contains(t, rec.Body.String(), response.MandatoryMissing("name").Message)
}

func TestHandlerFindUsers(t *testing.T) {
	e := echo.New()
	e.Use(middleware.NewLogging)
	e.Use(middleware.LoggingMiddleware)
	req := httptest.NewRequest(http.MethodGet, "/users", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	ctx := context.WithValue(c.Request().Context(), logger.LogContext, logger.NewZap())
	c.SetRequest(req.WithContext(ctx))
	createdAt := time.Now()
	mockUc := new(mockUsecase)
	handler := user.NewHandler(mockUc)
	result := []user.FindUserResponse{
		user.FindUserResponse{Id: "60d5ec49f1f1c939b4f2f0c1", Name: "Test1", Email: "test1@example.com", CreatedAt: createdAt},
		user.FindUserResponse{Id: "60d5ec49f1f1c939b4f2f0c2", Name: "Test2", Email: "test2@example.com", CreatedAt: createdAt},
	}
	mockUc.On("FindUsers", mock.Anything).Return(response.SuccessWithData(result), nil)

	err := handler.FindUsers(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	// assert.Contains(t, response.SuccessWithData(result), rec.Body.String())
}

func TestHandlerFindUserById_Invalid(t *testing.T) {
	e := echo.New()
	e.Use(middleware.NewLogging)
	e.Use(middleware.LoggingMiddleware)
	req := httptest.NewRequest(http.MethodGet, "/users/invalid", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	ctx := context.WithValue(c.Request().Context(), logger.LogContext, logger.NewZap())
	c.SetRequest(req.WithContext(ctx))
	c.SetParamNames(user.ParamID)
	c.SetParamValues("invalid")

	mockUc := new(mockUsecase)
	handler := user.NewHandler(mockUc)

	err := handler.FindUserById(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Contains(t, rec.Body.String(), "invalid")
}

func TestHandlerFindUserById_Valid(t *testing.T) {
	e := echo.New()
	e.Use(middleware.NewLogging)
	e.Use(middleware.LoggingMiddleware)
	validID := "60d5ec49f1f1c939b4f2f0c2"
	req := httptest.NewRequest(http.MethodGet, "/users/"+validID, nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	ctx := context.WithValue(c.Request().Context(), logger.LogContext, logger.NewZap())
	c.SetRequest(req.WithContext(ctx))
	c.SetParamNames(user.ParamID)
	c.SetParamValues(validID)

	mockUc := new(mockUsecase)
	handler := user.NewHandler(mockUc)

	mockUc.On("FindUserById", mock.Anything, validID).Return(response.SuccessWithData(user.User{Name: "Test"}), nil)

	err := handler.FindUserById(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), response.Success().Message)
}

func TestHandlerUpdateUser(t *testing.T) {
	e := echo.New()
	e.Use(middleware.NewLogging)
	e.Use(middleware.LoggingMiddleware)

	validID := primitive.NewObjectID()
	reqBody := user.UpdateRequest{Name: "Updated Name", Email: "updated@example.com"}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPut, "/users/"+validID.Hex(), bytes.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	ctx := context.WithValue(c.Request().Context(), logger.LogContext, logger.NewZap())
	c.SetRequest(req.WithContext(ctx))
	c.SetParamNames("id")
	c.SetParamValues(validID.Hex())

	mockUc := new(mockUsecase)
	handler := user.NewHandler(mockUc)

	mockUc.On("UpdateUser", mock.Anything, user.User{
		ID:    validID,
		Name:  reqBody.Name,
		Email: reqBody.Email,
	}).Return(response.Success(), nil)

	err := handler.UpdateUser(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), response.Success().Message)
}

func TestHandlerUpdateUser_InvalidId(t *testing.T) {
	e := echo.New()
	e.Use(middleware.NewLogging)
	e.Use(middleware.LoggingMiddleware)

	id := "invalidId"
	reqBody := user.UpdateRequest{Name: "Updated Name", Email: "updated@example.com"}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPut, "/users/"+id, bytes.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	ctx := context.WithValue(c.Request().Context(), logger.LogContext, logger.NewZap())
	c.SetRequest(req.WithContext(ctx))
	c.SetParamNames("id")
	c.SetParamValues(id)

	mockUc := new(mockUsecase)
	handler := user.NewHandler(mockUc)

	err := handler.UpdateUser(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Contains(t, rec.Body.String(), response.InvalidData(user.ParamID).Message)
}

func TestHandlerUpdateUser_InvalidBody(t *testing.T) {
	e := echo.New()
	e.Use(middleware.NewLogging)
	e.Use(middleware.LoggingMiddleware)

	validID := primitive.NewObjectID()
	reqBody := user.UpdateRequest{Name: "", Email: ""}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPut, "/users/"+validID.Hex(), bytes.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	ctx := context.WithValue(c.Request().Context(), logger.LogContext, logger.NewZap())
	c.SetRequest(req.WithContext(ctx))
	c.SetParamNames("id")
	c.SetParamValues(validID.Hex())

	mockUc := new(mockUsecase)
	handler := user.NewHandler(mockUc)

	mockUc.On("UpdateUser", mock.Anything, user.User{
		ID:    validID,
		Name:  reqBody.Name,
		Email: reqBody.Email,
	}).Return(response.Success(), nil)

	err := handler.UpdateUser(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Contains(t, rec.Body.String(), response.MandatoryMissing("name or email").Message)
}

func TestHandlerDeleteUser(t *testing.T) {
	e := echo.New()
	e.Use(middleware.NewLogging)
	e.Use(middleware.LoggingMiddleware)
	validID := "60d5ec49f1f1c939b4f2f0c2"
	req := httptest.NewRequest(http.MethodDelete, "/users/"+validID, nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	ctx := context.WithValue(c.Request().Context(), logger.LogContext, logger.NewZap())
	c.SetRequest(req.WithContext(ctx))
	c.SetParamNames(user.ParamID)
	c.SetParamValues(validID)

	mockUc := new(mockUsecase)
	handler := user.NewHandler(mockUc)

	mockUc.On("DeleteUser", mock.Anything, validID).Return(response.Success(), nil)

	err := handler.DeleteUser(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), response.Success().Message)
}

func TestHandlerDeleteUser_InvalidId(t *testing.T) {
	e := echo.New()
	e.Use(middleware.NewLogging)
	e.Use(middleware.LoggingMiddleware)
	req := httptest.NewRequest(http.MethodDelete, "/users/invalid", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	ctx := context.WithValue(c.Request().Context(), logger.LogContext, logger.NewZap())
	c.SetRequest(req.WithContext(ctx))
	c.SetParamNames(user.ParamID)
	c.SetParamValues("invalid")

	mockUc := new(mockUsecase)
	handler := user.NewHandler(mockUc)

	err := handler.DeleteUser(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Contains(t, rec.Body.String(), response.InvalidData(user.ParamID).Message)
}
