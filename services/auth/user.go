package auth

import "go.mongodb.org/mongo-driver/bson/primitive"

type User interface {
	IsAuthenticated() bool
	GetId() primitive.ObjectID
	GetName() string
}

type anonymousUser struct {
	id primitive.ObjectID
}

func NewAnonymousUser() *anonymousUser {
	return &anonymousUser{
		id: primitive.NewObjectID(),
	}
}

func (u *anonymousUser) GetId() primitive.ObjectID {
	return u.id
}
func (u *anonymousUser) GetName() string {
	return "anonymous"
}

func (u *anonymousUser) IsAuthenticated() bool {
	return false
}

func UserIdFromString(s string) (primitive.ObjectID, error) {
	id, err := primitive.ObjectIDFromHex(s)
	return id, err
}
