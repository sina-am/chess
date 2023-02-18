package types

import (
	"github.com/go-playground/validator"
	"github.com/sina-am/chess/core"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var validate *validator.Validate

func NewValidator() {
	validate = validator.New()
}

type RequestModel interface {
	Validate() error
}

type User struct {
	Id       primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Email    string             `json:"email"`
	Password string             `json:"-"`
	Picture  string             `json:"picture"`
	Gender   string             `json:"gender"`
	Name     string             `json:"name"`
	Score    int                `json:"score" `
	Games    []Game             `json:"games" `
}

func NewUser(email, plainPassword string) *User {
	return &User{
		Email:    email,
		Password: core.HashPassword(plainPassword),
		Games:    make([]Game, 0),
	}
}

type RegistrationRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

func (u *RegistrationRequest) Validate() error {
	return validate.Struct(u)
}

type AuthenticationRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

func (u *AuthenticationRequest) Validate() error {
	return validate.Struct(u)
}

type UpdateUserRequest struct {
	Gender  string `json:"gender"`
	Picture string `json:"picture"`
	Name    string `json:"name"`
}

func (u *UpdateUserRequest) Validate() error {
	return validate.Struct(u)
}
