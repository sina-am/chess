package users

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Gender string

const (
	MaleGender   Gender = "male"
	FemaleGender Gender = "female"
	OtherGender  Gender = "other"
)

type RequestModel interface {
	Validate() error
}

type Game struct{}

func NewUserId() primitive.ObjectID {
	return primitive.NewObjectID()
}
func UserIdFromString(s string) (primitive.ObjectID, error) {
	id, err := primitive.ObjectIDFromHex(s)
	return id, err
}

type UserType string

const (
	Guest  UserType = "GUEST"
	Simple UserType = "SIMPLE"
)

type UserI interface {
	IsAuthenticated() bool
	GetName() string
}

type anonymousUser struct {
}

func NewAnonymousUser() *anonymousUser {
	return &anonymousUser{}
}
func (u *anonymousUser) GetName() string {
	return "player #1090"
}

func (u *anonymousUser) IsAuthenticated() bool {
	return false
}

type User struct {
	Id          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Type        UserType           `json:"type" bson:"type"`
	Email       string             `json:"email" bson:"email,omitempty"`
	Password    string             `json:"-" bson:"password,omitempty"`
	Picture     string             `json:"picture" bson:"picture,omitempty"`
	Gender      Gender             `json:"gender" bson:"gender,omitempty"`
	Name        string             `json:"name" bson:"name"`
	Nationality string             `json:"nationality" bson:"nationality,omitempty"`
	Score       int                `json:"score" bson:"score"`
	Games       []Game             `json:"games" bson:"games"`
}

func (u *User) IsAuthenticated() bool {
	return true
}
func (u *User) GetName() string {
	return u.Name
}
func NewUser(email, name, plainPassword string) *User {
	return &User{
		Email:    email,
		Name:     name,
		Password: HashPassword(plainPassword),
		Games:    make([]Game, 0),
	}
}

func (u *User) GetType() UserType {
	return u.Type
}
func (u *User) GetId() primitive.ObjectID {
	return u.Id
}

type RegistrationRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
	Name     string `json:"name" validate:"required"`
	Gender   Gender `json:"gender" validate:"required,gender"`
}

type AuthenticationRequest struct {
	Email    string `form:"email" validate:"required,email"`
	Password string `form:"password" validate:"required"`
}

type UpdateUserRequest struct {
	Gender  Gender `json:"gender"`
	Picture string `json:"picture"`
	Name    string `json:"name"`
}
