package types

import (
	"github.com/sina-am/chess/chess"
	"github.com/sina-am/chess/services/auth"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Gender string

const (
	MaleGender   Gender = "male"
	FemaleGender Gender = "female"
	OtherGender  Gender = "other"
)

type Player struct {
	UserId primitive.ObjectID `json:"user_id" bson:"user_id"`
	Name   string             `json:"name"`
	Color  chess.Color        `json:"color" bson:"color"`
}
type Game struct {
	Id      primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Players []Player           `json:"players" bson:"players"`
	Winner  string             `json:"winner" bson:"winner"`
	Reason  string             `json:"reason" bson:"reason"`
}

func NewUserId() primitive.ObjectID {
	return primitive.NewObjectID()
}

type User struct {
	Id          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Email       string             `json:"email" bson:"email,omitempty"`
	Password    string             `json:"-" bson:"password,omitempty"`
	Picture     string             `json:"picture" bson:"picture,omitempty"`
	Gender      Gender             `json:"gender" bson:"gender,omitempty"`
	Name        string             `json:"name" bson:"name"`
	Nationality string             `json:"nationality" bson:"nationality,omitempty"`
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
