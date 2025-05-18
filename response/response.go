package response

import (
	"fmt"
	"net/http"
)

// response
const (
	success                = "0000"
	unexpectedRequest      = "4000"
	mandatoryMissing       = "4001"
	duplicatedRegistration = "4002"
	loginFailed            = "4003"
	invalidData            = "4004"
	userNotFound           = "4005"
	invalidAuthToken       = "4006"
	internalServerError    = "5000"
)

var message = map[string]string{
	success:                "Success",
	unexpectedRequest:      "Unexpected request",
	mandatoryMissing:       "%s is required",
	duplicatedRegistration: "An email has already been used",
	loginFailed:            "Login failed",
	invalidData:            "%s is invalid data",
	userNotFound:           "User not found",
	invalidAuthToken:       "Invalid authentication token",
	internalServerError:    "Internal server error",
}

var mapHTTPStatus = map[string]int{
	success:                http.StatusOK,
	unexpectedRequest:      http.StatusBadRequest,
	mandatoryMissing:       http.StatusBadRequest,
	duplicatedRegistration: http.StatusBadRequest,
	loginFailed:            http.StatusBadRequest,
	invalidData:            http.StatusBadRequest,
	userNotFound:           http.StatusNotFound,
	invalidAuthToken:       http.StatusUnauthorized,
	internalServerError:    http.StatusInternalServerError,
}

type StdResp[T any] struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Data    T      `json:"data,omitempty"`
}

func (s StdResp[T]) IsSuccess() bool {
	return s.Code == success
}

func (s StdResp[T]) WithHTTPStatus() (int, StdResp[T]) {
	return mapHTTPStatus[s.Code], s
}

func Success() *StdResp[any] {
	return &StdResp[any]{
		Code:    success,
		Message: message[success],
	}
}

func SuccessWithData(data any) *StdResp[any] {
	return &StdResp[any]{
		Code:    success,
		Message: message[success],
		Data:    data,
	}
}

func UnexpectedRequest() *StdResp[any] {
	return &StdResp[any]{
		Code:    unexpectedRequest,
		Message: message[unexpectedRequest],
	}
}

func UserNotFound() *StdResp[any] {
	return &StdResp[any]{
		Code:    userNotFound,
		Message: message[userNotFound],
	}
}

func DuplicatedRegistration() *StdResp[any] {
	return &StdResp[any]{
		Code:    duplicatedRegistration,
		Message: message[duplicatedRegistration],
	}
}

func LoginFail() *StdResp[any] {
	return &StdResp[any]{
		Code:    loginFailed,
		Message: message[loginFailed],
	}
}

func MandatoryMissing(field string) *StdResp[any] {
	return &StdResp[any]{
		Code:    mandatoryMissing,
		Message: fmt.Sprintf(message[mandatoryMissing], field),
	}
}

func InvalidData(field string) *StdResp[any] {
	return &StdResp[any]{
		Code:    invalidData,
		Message: fmt.Sprintf(message[invalidData], field),
	}
}

func Unauthorized() *StdResp[any] {
	return &StdResp[any]{
		Code:    invalidAuthToken,
		Message: message[invalidAuthToken],
	}
}

func InternalServerError() *StdResp[any] {
	return &StdResp[any]{
		Code:    internalServerError,
		Message: message[internalServerError],
	}
}
