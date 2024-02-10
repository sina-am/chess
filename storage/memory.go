package storage

import (
	"context"
	"errors"
	"fmt"

	"github.com/sina-am/chess/auth"
	"github.com/sina-am/chess/types"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type memoryStorage struct {
	users []*types.User
	games []*types.Game
}

func NewMemoryStorage() *memoryStorage {
	return &memoryStorage{
		users: make([]*types.User, 0),
		games: make([]*types.Game, 0),
	}
}

func (db *memoryStorage) UpdateUser(ctx context.Context, user *types.User) error {
	for i := range db.users {
		if db.users[i].Id == user.Id {
			db.users[i] = user
			return nil
		}
	}
	return ErrNoRecord
}

func (db *memoryStorage) GetAllUsers(ctx context.Context) ([]*types.User, error) {
	return db.users, nil
}

func (db *memoryStorage) InsertUser(ctx context.Context, user *types.User) error {
	user.Id = primitive.NewObjectID()
	db.users = append(db.users, user)
	return nil
}

func (db *memoryStorage) GetUserById(ctx context.Context, id primitive.ObjectID) (*types.User, error) {
	for _, user := range db.users {
		if user.Id == id {
			return user, nil
		}
	}
	return nil, ErrNoRecord
}

func (db *memoryStorage) GetUserByEmail(ctx context.Context, email string) (*types.User, error) {
	for _, user := range db.users {
		if user.Email == email {
			return user, nil
		}
	}
	return nil, ErrNoRecord
}

func (db *memoryStorage) AuthenticateUser(ctx context.Context, email string, plainPassword string) (*types.User, error) {
	user, err := db.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, ErrNoRecord) {
			return nil, ErrAuthentication
		}
		return nil, err
	}

	if err := auth.VerifyPassword(plainPassword, user.Password); err != nil {
		return nil, ErrAuthentication
	}
	return user, nil
}

func (db *memoryStorage) InsertGame(ctx context.Context, game *types.Game) error {
	if len(game.Players) != 2 {
		return fmt.Errorf("invalid number of players")
	}
	user1, err := db.GetUserById(ctx, game.Players[0].UserId)
	if err != nil {
		return err
	}
	user2, err := db.GetUserById(ctx, game.Players[1].UserId)
	if err != nil {
		return err
	}
	user1.Games = append(user1.Games, *game)
	user2.Games = append(user2.Games, *game)
	return nil
}
