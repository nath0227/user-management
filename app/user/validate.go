package user

import (
	"net/mail"
	"strings"
	"user-management/response"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (r CreateRequest) RequestValidation() *response.StdResp[any] {
	if checkLen(r.Name) == 0 {
		return response.MandatoryMissing("name")
	}
	if checkLen(r.Email) == 0 {
		return response.MandatoryMissing("email")
	}
	if !isValidEmail(r.Email) {
		return response.InvalidData("email")
	}
	if checkLen(r.Password) == 0 {
		return response.MandatoryMissing("password")
	}
	return response.Success()
}

func (r SignInRequest) RequestValidation() *response.StdResp[any] {
	if checkLen(r.Email) == 0 {
		return response.MandatoryMissing("email")
	}
	if !isValidEmail(r.Email) {
		return response.InvalidData("email")
	}
	if checkLen(r.Password) == 0 {
		return response.MandatoryMissing("password")
	}
	return response.Success()
}

func (r UpdateRequest) RequestValidation() *response.StdResp[any] {
	if checkLen(r.Email) == 0 && checkLen(r.Name) == 0 {
		return response.MandatoryMissing("name or email")
	}
	if checkLen(r.Email) != 0 && !isValidEmail(r.Email) {
		return response.InvalidData("email")
	}
	return response.Success()
}

func IdValidation(id string) *response.StdResp[any] {
	_, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return response.InvalidData(ParamID)
	}
	return response.Success()
}

func checkLen(s string) int {
	return len([]rune(strings.TrimSpace(s)))
}

func isValidEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}
