package database

import (
	"context"
	"testing"

	"github.com/sina-am/chess/types"
	"github.com/stretchr/testify/assert"
)

func getTestMemoryDatabase() *memoryDatabase {
	return NewMemoryDatabase(context.TODO())
}

func TestMemoryInsertUser(t *testing.T) {
	db := getTestMemoryDatabase()

	user := types.NewUser("test@gmail.com", "1234")
	assert.Nil(t, db.InsertUser(context.TODO(), user))

	dbUser, err := db.GetUserByEmail(context.TODO(), "test@gmail.com")
	assert.Nil(t, err)
	assert.NotNil(t, dbUser.Id)
	assert.Equal(t, dbUser.Password, user.Password)
	assert.Equal(t, dbUser.Email, user.Email)
}

func TestMemoryUpdateUser(t *testing.T) {
	db := getTestMemoryDatabase()

	user := types.NewUser("test@gmail.com", "1234")
	assert.Nil(t, db.InsertUser(context.TODO(), user))

	updatedUser, err := db.GetUserByEmail(context.TODO(), "test@gmail.com")
	assert.Nil(t, err)

	updatedUser.Gender = "male"
	updatedUser.Name = "test user"
	assert.Nil(t, db.UpdateUser(context.TODO(), updatedUser))

	dbUser, err := db.GetUserByEmail(context.TODO(), "test@gmail.com")
	assert.Nil(t, err)
	assert.Equal(t, dbUser, updatedUser)
}

func TestMemoryGetUserById(t *testing.T) {
	db := getTestMemoryDatabase()

	user := types.NewUser("test@gmail.com", "1234")
	assert.Nil(t, db.InsertUser(context.TODO(), user))

	dbUser, err := db.GetUserByEmail(context.TODO(), "test@gmail.com")
	assert.Nil(t, err)

	dbUser2, err := db.GetUserById(context.TODO(), dbUser.Id)
	assert.Nil(t, err)
	assert.Equal(t, dbUser, dbUser2)
}

func TestMemoryGetUserByEmail(t *testing.T) {
	db := getTestMemoryDatabase()

	user := types.NewUser("test@gmail.com", "1234")
	assert.Nil(t, db.InsertUser(context.TODO(), user))

	dbUser, err := db.GetUserByEmail(context.TODO(), "test@gmail.com")
	assert.Nil(t, err)
	assert.Equal(t, dbUser.Password, user.Password)
	assert.Equal(t, dbUser.Email, user.Email)

	// not found record
	_, err = db.GetUserByEmail(context.TODO(), "invalid@gmail.com")
	assert.ErrorIs(t, err, ErrNoRecord)
}

func TestMemoryAuthenticateUser(t *testing.T) {
	db := getTestMemoryDatabase()

	user := types.NewUser("test@gmail.com", "1234")
	assert.Nil(t, db.InsertUser(context.TODO(), user))

	dbUser, err := db.AuthenticateUser(context.TODO(), "test@gmail.com", "1234")
	assert.Nil(t, err)
	assert.NotNil(t, dbUser)

	// wrong email
	_, err = db.AuthenticateUser(context.TODO(), "invalid@gmail.com", "1234")
	assert.ErrorIs(t, err, ErrAuthentication)

	// wrong password
	_, err = db.AuthenticateUser(context.TODO(), "test@gmail.com", "invalid")
	assert.ErrorIs(t, err, ErrAuthentication)
}

func TestMemoryGetAllUsers(t *testing.T) {

}

func TestMemoryInsertGame(t *testing.T) {

}

func TestMemoryGetUserGame(t *testing.T) {

}

func TestMemoryUpdateGame(t *testing.T) {

}
