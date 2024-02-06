package types

import (
	"github.com/sina-am/chess/auth"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Gender string

const (
	MaleGender   Gender = "male"
	FemaleGender Gender = "female"
	OtherGender  Gender = "other"
)

type Game struct{}

func NewUserId() primitive.ObjectID {
	return primitive.NewObjectID()
}

type UserType string

const (
	Guest  UserType = "GUEST"
	Simple UserType = "SIMPLE"
)

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

func NewUser(email, name, plainPassword string) *User {
	return &User{
		Email:    email,
		Name:     name,
		Password: auth.HashPassword(plainPassword),
		Games:    make([]Game, 0),
	}
}

func (u *User) IsAuthenticated() bool {
	return true
}

func (u *User) GetName() string {
	return u.Name
}

func (u *User) GetId() primitive.ObjectID {
	return u.Id
}

func (u *User) GetType() UserType {
	return u.Type
}
