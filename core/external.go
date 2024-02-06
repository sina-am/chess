package core

import "go.mongodb.org/mongo-driver/mongo"

type ExternalServices struct {
	MongoClient *mongo.Client
}
