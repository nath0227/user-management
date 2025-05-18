package user

import (
	"context"
	"errors"
	"time"
	"user-management/config"
	"user-management/response"

	jwt "github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type Repository interface {
	CreateUser(ctx context.Context, user User) (string, error)
	FindUserByEmail(ctx context.Context, email string) (User, error)
	FindUserById(ctx context.Context, id string) (FindUserResponse, error)
	FindUsers(ctx context.Context) ([]FindUserResponse, error)
	UpdateUser(ctx context.Context, user User) (int64, error)
	DeleteUser(ctx context.Context, id string) (int64, error)
}

type usecase struct {
	cfgCrypto config.CryptoCredential
	repo      Repository
}

func NewUsecase(cfg config.CryptoCredential, r Repository) *usecase {
	return &usecase{
		cfgCrypto: cfg,
		repo:      r,
	}
}

func (u *usecase) CreateUser(ctx context.Context, req CreateRequest) (*response.StdResp[any], error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	uid, err := u.repo.CreateUser(ctx, User{
		Name:     req.Name,
		Email:    req.Email,
		Password: string(hashedPassword),
	})
	if err != nil {
		if errors.Is(err, ErrEmailAlreadyExists) {
			return response.DuplicatedRegistration(), nil
		}
		return nil, err
	}
	return response.SuccessWithData(CreateResponse{uid}), nil
}

func (u *usecase) Login(ctx context.Context, req SignInRequest) (*response.StdResp[any], error) {
	result, err := u.repo.FindUserByEmail(ctx, req.Email)
	if err != nil {
		if errors.Is(err, ErrUserOrPasswordIsWrong) {
			return response.LoginFail(), nil
		}
		return nil, err
	}

	if !ValidPassword(result.Password, req.Password) {
		return response.LoginFail(), nil
	}

	expAt := time.Now().Add(u.cfgCrypto.JwtExpireDuration)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.RegisteredClaims{
			Subject:   result.Email,
			ExpiresAt: jwt.NewNumericDate(expAt),
		},
	)

	signedToken, err := token.SignedString([]byte(u.cfgCrypto.JwtKey))
	if err != nil {
		return nil, err
	}
	return response.SuccessWithData(&SignInResponse{
		Token:     signedToken,
		ExpiresAt: expAt,
	}), err
}

func (u *usecase) FindUsers(ctx context.Context) (*response.StdResp[any], error) {
	users, err := u.repo.FindUsers(ctx)
	if err != nil {
		return nil, err
	}
	return response.SuccessWithData(users), nil
}

func (u *usecase) FindUserById(ctx context.Context, id string) (*response.StdResp[any], error) {
	user, err := u.repo.FindUserById(ctx, id)
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			return response.UserNotFound(), nil
		}
		return nil, err
	}
	return response.SuccessWithData(user), nil
}

func (u *usecase) UpdateUser(ctx context.Context, user User) (*response.StdResp[any], error) {
	updateCount, err := u.repo.UpdateUser(ctx, user)
	if err != nil {
		if errors.Is(err, ErrEmailAlreadyExists) {
			return response.DuplicatedRegistration(), nil
		}
		return nil, err
	}
	if updateCount == 0 {
		return response.UserNotFound(), nil
	}
	return response.Success(), err
}

func (u *usecase) DeleteUser(ctx context.Context, id string) (*response.StdResp[any], error) {
	delCount, err := u.repo.DeleteUser(ctx, id)
	if err != nil {
		return nil, err
	}
	if delCount == 0 {
		return response.UserNotFound(), nil
	}
	return response.Success(), nil
}

func ValidPassword(hashedPassword, plainPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plainPassword))
	return err == nil
}
