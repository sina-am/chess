package storage

import (
	"context"
	"errors"

	"github.com/sina-am/chess/types"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	ErrNoRecord       = errors.New("no record found")
	ErrAuthentication = errors.New("email or password is not correct")
)

type Storage interface {
	GetAllUsers(ctx context.Context) ([]*types.User, error)
	InsertUser(ctx context.Context, user *types.User) error
	UpdateUser(ctx context.Context, user *types.User) error
	GetUserById(ctx context.Context, id primitive.ObjectID) (*types.User, error)
	GetUserByEmail(ctx context.Context, email string) (*types.User, error)
	AuthenticateUser(ctx context.Context, email string, plainPassword string) (*types.User, error)
	InsertGame(ctx context.Context, game *types.Game) error
}
