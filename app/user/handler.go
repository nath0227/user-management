package user

import (
	"context"
	"log"
	"net/http"
	"user-management/response"

	echo "github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Usecase interface {
	CreateUser(ctx context.Context, req User) (*response.StdResp[any], error)
	Login(ctx context.Context, req SignInRequest) (*response.StdResp[any], error)
	FindUsers(ctx context.Context) (*response.StdResp[any], error)
	FindUserById(ctx context.Context, id string) (*response.StdResp[any], error)
	UpdateUser(ctx context.Context, user User) (*response.StdResp[any], error)
	DeleteUser(ctx context.Context, id string) (*response.StdResp[any], error)
}

type Handler interface {
	Login(c echo.Context) error
	CreateUser(c echo.Context) error
	FindUsers(c echo.Context) error
	FindUserById(c echo.Context) error
	UpdateUser(c echo.Context) error
	DeleteUser(c echo.Context) error
}

type handler struct {
	logger  *log.Logger
	usecase Usecase
}

func NewHandler(logger *log.Logger, u Usecase) *handler {
	return &handler{
		logger:  logger,
		usecase: u,
	}
}

func (h *handler) Login(c echo.Context) error {
	ctx := c.Request().Context()
	var request SignInRequest
	err := c.Bind(&request)
	if err != nil {
		h.logger.Println("[Handler] Bind request error:", err.Error())
		return c.JSON(response.UnexpectedRequest().WithHTTPStatus())
	}

	if resp := request.RequestValidation(); !resp.IsSuccess() {
		return c.JSON(resp.WithHTTPStatus())
	}
	sr, err := h.usecase.Login(ctx, request)
	if err != nil {
		h.logger.Println("[Handler] Service error:", err.Error())
		if err.Error() == UserOrPasswordIsWrong {
			return c.JSON(response.LoginFail().WithHTTPStatus())
		}
		return c.JSON(response.InternalServerError().WithHTTPStatus())
	}
	if !sr.IsSuccess() {
		return c.JSON(sr.WithHTTPStatus())
	}

	c.SetCookie(newCookie(sr.Data.(*SignInResponse)))
	// return c.JSON(response.Success().WithHTTPStatus())
	return c.JSON(response.SuccessWithData(sr).WithHTTPStatus())
}

func (h *handler) CreateUser(c echo.Context) error {
	ctx := c.Request().Context()
	var request User
	err := c.Bind(&request)
	if err != nil {
		h.logger.Println("[Handler] Bind request error:", err.Error())
		return c.JSON(response.UnexpectedRequest().WithHTTPStatus())
	}

	if resp := request.RequestValidation(); !resp.IsSuccess() {
		return c.JSON(resp.WithHTTPStatus())
	}

	resp, err := h.usecase.CreateUser(ctx, request)
	if err != nil {
		h.logger.Println("[Handler] Service error:", err.Error())
		if err.Error() == EmailAlreadyExists {
			return c.JSON(response.DuplicatedRegistration().WithHTTPStatus())
		}
		return c.JSON(response.InternalServerError().WithHTTPStatus())
	}

	return c.JSON(resp.WithHTTPStatus())
}

func (h *handler) FindUsers(c echo.Context) error {
	ctx := c.Request().Context()
	resp, err := h.usecase.FindUsers(ctx)
	if err != nil {
		h.logger.Println("[Handler] Service error:", err.Error())
		return c.JSON(response.InternalServerError().WithHTTPStatus())
	}
	return c.JSON(resp.WithHTTPStatus())
}

func (h *handler) FindUserById(c echo.Context) error {
	ctx := c.Request().Context()
	paramId := c.Param(ParamID)
	if respValidate := IdValidation(paramId); !respValidate.IsSuccess() {
		return c.JSON(response.InvalidData(ParamID).WithHTTPStatus())
	}

	resp, err := h.usecase.FindUserById(ctx, paramId)
	if err != nil {
		h.logger.Println("[Handler] Service error:", err.Error())
		return c.JSON(response.InternalServerError().WithHTTPStatus())
	}
	return c.JSON(resp.WithHTTPStatus())
}

func (h *handler) UpdateUser(c echo.Context) error {
	ctx := c.Request().Context()
	var request UpdateRequest
	paramId := c.Param(ParamID)
	userID, err := primitive.ObjectIDFromHex(paramId)
	if err != nil {
		return c.JSON(response.LoginFail().WithHTTPStatus())
	}
	err = c.Bind(&request)
	if err != nil {
		h.logger.Println("[Handler] Bind request error:", err.Error())
		return c.JSON(response.UnexpectedRequest().WithHTTPStatus())
	}

	if resp := request.RequestValidation(); !resp.IsSuccess() {
		return c.JSON(resp.WithHTTPStatus())
	}

	resp, err := h.usecase.UpdateUser(ctx, User{
		ID:    userID,
		Name:  request.Name,
		Email: request.Email,
	})
	if err != nil {
		h.logger.Println("[Handler] Service error:", err.Error())
		return c.JSON(response.InternalServerError().WithHTTPStatus())
	}
	return c.JSON(resp.WithHTTPStatus())
}

func (h *handler) DeleteUser(c echo.Context) error {
	ctx := c.Request().Context()
	paramId := c.Param(ParamID)
	if respValidate := IdValidation(paramId); !respValidate.IsSuccess() {
		return c.JSON(response.InvalidData(ParamID).WithHTTPStatus())
	}

	resp, err := h.usecase.DeleteUser(ctx, paramId)
	if err != nil {
		h.logger.Println("[Handler] Service error:", err.Error())
		return c.JSON(response.InternalServerError().WithHTTPStatus())
	}
	return c.JSON(resp.WithHTTPStatus())
}

func newCookie(s *SignInResponse) *http.Cookie {
	return &http.Cookie{
		Name:     "token",
		Value:    s.Token,
		Expires:  s.ExpiresAt,
		HttpOnly: true,
		Secure:   true,
		Path:     "/",
		SameSite: http.SameSiteStrictMode,
	}
}
