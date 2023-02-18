package database

import (
	"context"
	"errors"
	"fmt"

	"github.com/sina-am/chess/core"
	"github.com/sina-am/chess/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type mongoDatabase struct {
	client       *mongo.Client
	databaseName string
}

func NewMongoDatabase(ctx context.Context, uri string, dbName string) (*mongoDatabase, error) {
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}

	insertMongoIndexes(ctx, client.Database(dbName))
	return &mongoDatabase{
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

func (db *mongoDatabase) getUserCollection() *mongo.Collection {
	return db.client.Database(db.databaseName).Collection("users")
}

func (db *mongoDatabase) findUser(ctx context.Context, filter any) (*types.User, error) {
	collection := db.getUserCollection()
	document := collection.FindOne(ctx, filter)
	user := &types.User{}

	if err := document.Decode(user); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, ErrNoRecord
		}
		return nil, err
	}

	return user, nil
}

func (db *mongoDatabase) UpdateUser(ctx context.Context, user *types.User) error {
	collection := db.getUserCollection()
	_, err := collection.UpdateOne(ctx, bson.M{"_id": user.Id}, bson.M{"$set": user})
	return err
}

func (db *mongoDatabase) GetAllUsers(ctx context.Context) ([]*types.User, error) {
	collection := db.getUserCollection()
	cur, err := collection.Find(ctx, bson.D{})
	if err != nil {
		return nil, err
	}

	users := []*types.User{}
	for cur.Next(ctx) {
		user := &types.User{}
		if err := cur.Decode(user); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

func (db *mongoDatabase) InsertUser(ctx context.Context, user *types.User) error {
	collection := db.getUserCollection()
	if _, err := db.GetUserByEmail(ctx, user.Email); err == nil {
		return fmt.Errorf("user with email %s already exist", user.Email)
	}
	_, err := collection.InsertOne(ctx, user)

	return err
}

func (db *mongoDatabase) GetUserById(ctx context.Context, id primitive.ObjectID) (*types.User, error) {
	return db.findUser(
		ctx,
		bson.D{{Key: "_id", Value: id}},
	)
}

func (db *mongoDatabase) GetUserByEmail(ctx context.Context, email string) (*types.User, error) {
	return db.findUser(
		ctx,
		bson.D{{Key: "email", Value: email}},
	)
}

func (db *mongoDatabase) InsertGame(ctx context.Context, game *types.Game) error {
	collection := db.getUserCollection()
	game.Id = primitive.NewObjectID()
	_, err := collection.UpdateMany(
		ctx,
		bson.M{"$or": bson.A{bson.M{"_id": game.Players[0].UserId}, bson.M{"_id": game.Players[1].UserId}}},
		bson.M{"$push": bson.M{"games": game}},
	)
	return err
}

func (db *mongoDatabase) GetUserGame(ctx context.Context, userId string, gameId string) (*types.Game, error) {
	return nil, nil
}

func (db *mongoDatabase) UpdateGame(ctx context.Context, game *types.Game) error {
	collection := db.getUserCollection()
	_, err := collection.UpdateMany(
		ctx,
		bson.M{"$or": bson.A{
			bson.M{"_id": game.Players[0].UserId},
			bson.M{"_id": game.Players[1].UserId}},
			"games._id": game.Id},
		bson.M{"$set": bson.M{"games.$": game}},
	)
	return err
}

func (db *mongoDatabase) AuthenticateUser(ctx context.Context, email string, plainPassword string) (*types.User, error) {
	user, err := db.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, ErrNoRecord) {
			return nil, ErrAuthentication
		}
		return nil, err
	}

	if err := core.VerifyPassword(plainPassword, user.Password); err != nil {
		return nil, ErrAuthentication
	}
	return user, nil
}
