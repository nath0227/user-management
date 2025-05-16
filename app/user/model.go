package user

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name      string             `bson:"name" json:"name"`
	Email     string             `bson:"email" json:"email"`
	Password  string             `bson:"password" json:"password"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
}

type SignInRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type SignInResponse struct {
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expire_at"`
}

type CreateResponse struct {
	Id string `json:"id"`
}

type UpdateRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}
