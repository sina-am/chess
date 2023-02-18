package database

import (
	"context"
	"fmt"

	"github.com/sina-am/chess/core"
	"github.com/sina-am/chess/types"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type memoryDatabase struct {
	users map[primitive.ObjectID]interface{}
}

func NewMemoryDatabase(ctx context.Context) *memoryDatabase {
	return &memoryDatabase{
		users: make(map[primitive.ObjectID]interface{}, 10),
	}
}
func (db *memoryDatabase) GetUserById(ctx context.Context, id primitive.ObjectID) (*types.User, error) {
	if user, found := db.users[id]; found {
		return user.(*types.User), nil
	}
	return nil, ErrNoRecord
}

func (db *memoryDatabase) GetUserByEmail(ctx context.Context, email string) (*types.User, error) {
	for _, user := range db.users {
		if user.(*types.User).Email == email {
			return user.(*types.User), nil
		}
	}
	return nil, ErrNoRecord
}

func (db *memoryDatabase) InsertUser(ctx context.Context, user *types.User) error {
	if _, err := db.GetUserByEmail(ctx, user.Email); err == nil {
		return fmt.Errorf("user with email %s already exist", user.Email)
	}
	user.Id = primitive.NewObjectID()
	db.users[user.Id] = user
	return nil
}

func (db *memoryDatabase) AuthenticateUser(ctx context.Context, email string, plainPassword string) (*types.User, error) {
	user, err := db.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, ErrAuthentication
	}

	if err := core.VerifyPassword(plainPassword, user.Password); err != nil {
		return nil, ErrAuthentication
	}

	return user, nil
}

func (db *memoryDatabase) InsertGame(ctx context.Context, game *types.Game) error {
	return nil
}

func (db *memoryDatabase) GetUserGame(ctx context.Context, userId string, gameId string) (*types.Game, error) {
	return nil, nil
}
func (db *memoryDatabase) UpdateGame(ctx context.Context, game *types.Game) error {
	return nil
}

func (db *memoryDatabase) GetAllUsers(ctx context.Context) ([]*types.User, error) {
	userList := make([]*types.User, 0)
	for _, user := range db.users {
		userList = append(userList, user.(*types.User))
	}
	return userList, nil
}

func (db *memoryDatabase) UpdateUser(ctx context.Context, user *types.User) error {
	if _, found := db.users[user.Id]; !found {
		return ErrNoRecord
	}

	db.users[user.Id] = user
	return nil
}
