package apis

import (
	"context"
	"testing"

	"github.com/sina-am/chess/database"
	"github.com/sina-am/chess/types"
	"github.com/stretchr/testify/assert"
)

func TestAuthenticate(t *testing.T) {
	user := types.NewUser("test@gmail.com", "test")
	db := database.NewMemoryDatabase(context.TODO())
	db.InsertUser(context.TODO(), user)

	auth := NewJWTAuthentication("verysecretkey", db)
	token, err := auth.ObtainToken(user)
	assert.Nil(t, err)

	authUser, err := auth.Authenticate(context.TODO(), token)
	assert.Nil(t, err)
	assert.Equal(t, user, authUser)
}

func TestObtainToken(t *testing.T) {
	user := types.NewUser("test@gmail.com", "test")
	db := database.NewMemoryDatabase(context.TODO())
	db.InsertUser(context.TODO(), user)

	auth := NewJWTAuthentication("verysecretkey", db)
	token, err := auth.ObtainToken(user)
	assert.Nil(t, err)

	id, err := auth.GetUserIdFromToken(token)
	assert.Nil(t, err)
	assert.Equal(t, user.Id, id)
}
