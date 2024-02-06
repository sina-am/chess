package users

import "github.com/sina-am/chess/types"

type RegistrationRequest struct {
	Email    string       `json:"email" validate:"required,email"`
	Password string       `json:"password" validate:"required,min=6"`
	Name     string       `json:"name" validate:"required"`
	Gender   types.Gender `json:"gender" validate:"required,gender"`
}

type AuthenticationRequest struct {
	Email    string `form:"email" validate:"required,email"`
	Password string `form:"password" validate:"required"`
}

type UpdateUserRequest struct {
	Gender  types.Gender `json:"gender"`
	Picture string       `json:"picture"`
	Name    string       `json:"name"`
}
