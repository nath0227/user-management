package user

import "errors"

const (
	ParamID = "id"
)

var (
	ErrEmailAlreadyExists    = errors.New("Email already exists")
	ErrUserOrPasswordIsWrong = errors.New("User or password is wrong")
	ErrUserNotFound          = errors.New("User not found")
)
