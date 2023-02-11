package types

import (
	"github.com/go-playground/validator"
	"github.com/sina-am/chess/core"
)

var validate *validator.Validate

func NewValidator() {
	validate = validator.New()
}

type RequestModel interface {
	Validate() error
}

type User struct {
	Id       string `json:"id" bson:"_id,omitempty"`
	Email    string `json:"email"`
	Password string `json:"-"`
	Picture  string `json:"picture"`
	Gender   string `json:"gender"`
	Name     string `json:"name"`
	Score    int    `json:"score" `
	Games    []Game `json:"games" `
}

func NewUser(email, plainPassword string) *User {
	return &User{
		Email:    email,
		Password: core.HashPassword(plainPassword),
		Games:    make([]Game, 0),
	}
}

type RegisterUserRequest struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}

func (u *RegisterUserRequest) Validate() error {
	return validate.Struct(u)
}

type AuthenticateUserRequest struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}

func (u *AuthenticateUserRequest) Validate() error {
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
