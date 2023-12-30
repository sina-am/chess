package users

import (
	"context"
	"errors"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	ErrNoRecord       = errors.New("no record found")
	ErrAuthentication = errors.New("email or password is not correct")
)

type Storage interface {
	GetAllUsers(ctx context.Context) ([]*User, error)
	InsertUser(ctx context.Context, user *User) error
	UpdateUser(ctx context.Context, user *User) error
	GetUserById(ctx context.Context, id primitive.ObjectID) (*User, error)
	GetUserByEmail(ctx context.Context, email string) (*User, error)
	AuthenticateUser(ctx context.Context, email string, plainPassword string) (*User, error)
}

type mongoStorage struct {
	client       *mongo.Client
	databaseName string
}

func NewMongoStorage(ctx context.Context, client *mongo.Client, dbName string) (*mongoStorage, error) {
	insertMongoIndexes(ctx, client.Database(dbName))
	return &mongoStorage{
		client:       client,
		databaseName: dbName,
	}, nil
}

func insertMongoIndexes(ctx context.Context, database *mongo.Database) {
	collection := database.Collection("users")
	indexModel := mongo.IndexModel{
		Keys:    bson.D{{Key: "email", Value: 1}},
		Options: options.Index().SetUnique(true),
	}
	collection.Indexes().CreateOne(ctx, indexModel)

}

func (db *mongoStorage) getUserCollection() *mongo.Collection {
	return db.client.Database(db.databaseName).Collection("users")
}

func (db *mongoStorage) findUser(ctx context.Context, filter any) (*User, error) {
	collection := db.getUserCollection()
	document := collection.FindOne(ctx, filter)
	user := &User{}

	if err := document.Decode(user); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, ErrNoRecord
		}
		return nil, err
	}

	return user, nil
}

func (db *mongoStorage) UpdateUser(ctx context.Context, user *User) error {
	collection := db.getUserCollection()
	_, err := collection.UpdateOne(ctx, bson.M{"_id": user.Id}, bson.M{"$set": user})
	return err
}

func (db *mongoStorage) GetAllUsers(ctx context.Context) ([]*User, error) {
	collection := db.getUserCollection()
	cur, err := collection.Find(ctx, bson.D{})
	if err != nil {
		return nil, err
	}

	users := []*User{}
	for cur.Next(ctx) {
		user := &User{}
		if err := cur.Decode(user); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

func (db *mongoStorage) InsertUser(ctx context.Context, user *User) error {
	collection := db.getUserCollection()
	if _, err := db.GetUserByEmail(ctx, user.Email); err == nil {
		return fmt.Errorf("user with email %s already exist", user.Email)
	}
	_, err := collection.InsertOne(ctx, user)

	return err
}

func (db *mongoStorage) GetUserById(ctx context.Context, id primitive.ObjectID) (*User, error) {
	return db.findUser(
		ctx,
		bson.D{{Key: "_id", Value: id}},
	)
}

func (db *mongoStorage) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	return db.findUser(
		ctx,
		bson.D{{Key: "email", Value: email}},
	)
}

// func (db *mongoStorage) InsertGame(ctx context.Context, game *Game) error {
// 	collection := db.getUserCollection()
// 	game.Id = primitive.NewObjectID()
// 	_, err := collection.UpdateMany(
// 		ctx,
// 		bson.M{"$or": bson.A{bson.M{"_id": game.Players[0].UserId}, bson.M{"_id": game.Players[1].UserId}}},
// 		bson.M{"$push": bson.M{"games": game}},
// 	)
// 	return err
// }

// func (db *mongoStorage) GetUserGame(ctx context.Context, userId string, gameId string) (*Game, error) {
// 	return nil, nil
// }

// func (db *mongoStorage) UpdateGame(ctx context.Context, game *Game) error {
// 	collection := db.getUserCollection()
// 	_, err := collection.UpdateMany(
// 		ctx,
// 		bson.M{"$or": bson.A{
// 			bson.M{"_id": game.Players[0].UserId},
// 			bson.M{"_id": game.Players[1].UserId}},
// 			"games._id": game.Id},
// 		bson.M{"$set": bson.M{"games.$": game}},
// 	)
// 	return err
// }

func (db *mongoStorage) AuthenticateUser(ctx context.Context, email string, plainPassword string) (*User, error) {
	user, err := db.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, ErrNoRecord) {
			return nil, ErrAuthentication
		}
		return nil, err
	}

	if err := VerifyPassword(plainPassword, user.Password); err != nil {
		return nil, ErrAuthentication
	}
	return user, nil
}

type memoryStorage struct {
	users []*User
	games []*Game
}

func NewMemoryStorage() *memoryStorage {
	return &memoryStorage{
		users: make([]*User, 0),
		games: make([]*Game, 0),
	}
}

func (db *memoryStorage) UpdateUser(ctx context.Context, user *User) error {
	for i := range db.users {
		if db.users[i].Id == user.Id {
			db.users[i] = user
			return nil
		}
	}
	return ErrNoRecord
}

func (db *memoryStorage) GetAllUsers(ctx context.Context) ([]*User, error) {
	return db.users, nil
}

func (db *memoryStorage) InsertUser(ctx context.Context, user *User) error {
	user.Id = primitive.NewObjectID()
	db.users = append(db.users, user)
	return nil
}

func (db *memoryStorage) GetUserById(ctx context.Context, id primitive.ObjectID) (*User, error) {
	for _, user := range db.users {
		if user.Id == id {
			return user, nil
		}
	}
	return nil, ErrNoRecord
}

func (db *memoryStorage) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	for _, user := range db.users {
		if user.Email == email {
			return user, nil
		}
	}
	return nil, ErrNoRecord
}

func (db *memoryStorage) AuthenticateUser(ctx context.Context, email string, plainPassword string) (*User, error) {
	user, err := db.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, ErrNoRecord) {
			return nil, ErrAuthentication
		}
		return nil, err
	}

	if err := VerifyPassword(plainPassword, user.Password); err != nil {
		return nil, ErrAuthentication
	}
	return user, nil
}
