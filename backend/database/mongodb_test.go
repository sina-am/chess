package database

import (
	"context"
	"testing"

	"github.com/sina-am/chess/types"
	"github.com/stretchr/testify/assert"
)

func getTestMongoDatabase() *mongoDatabase {
	client, err := NewMongoDatabase(context.TODO(), "mongodb://localhost", "chess_test")
	if err != nil {
		panic(err)
	}
	return client
}

func deleteTestMongoDatabase(db *mongoDatabase) {
	db.client.Database("chess_test").Drop(context.TODO())
}

func TestMongoInsertUser(t *testing.T) {
	db := getTestMongoDatabase()
	defer deleteTestMongoDatabase(db)

	user := types.NewUser("test@gmail.com", "1234")
	assert.Nil(t, db.InsertUser(context.TODO(), user))

	dbUser, err := db.GetUserByEmail(context.TODO(), "test@gmail.com")
	assert.Nil(t, err)
	assert.NotNil(t, dbUser.Id)
	assert.Equal(t, dbUser.Password, user.Password)
	assert.Equal(t, dbUser.Email, user.Email)
}

func TestMongoUpdateUser(t *testing.T) {
	db := getTestMongoDatabase()
	defer deleteTestMongoDatabase(db)

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

func TestMongoGetUserById(t *testing.T) {
	db := getTestMongoDatabase()
	defer deleteTestMongoDatabase(db)

	user := types.NewUser("test@gmail.com", "1234")
	assert.Nil(t, db.InsertUser(context.TODO(), user))

	dbUser, err := db.GetUserByEmail(context.TODO(), "test@gmail.com")
	assert.Nil(t, err)

	dbUser2, err := db.GetUserById(context.TODO(), dbUser.Id)
	assert.Nil(t, err)
	assert.Equal(t, dbUser, dbUser2)
}

func TestMongoGetUserByEmail(t *testing.T) {
	db := getTestMongoDatabase()
	defer deleteTestMongoDatabase(db)

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

func TestMongoAuthenticateUser(t *testing.T) {
	db := getTestMongoDatabase()
	defer deleteTestMongoDatabase(db)

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

func TestMongoGetAllUsers(t *testing.T) {

}

func TestMongoInsertGame(t *testing.T) {

}

func TestMongoGetUserGame(t *testing.T) {

}

func TestMongoUpdateGame(t *testing.T) {

}
