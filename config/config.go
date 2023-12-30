package config

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Database struct {
	Uri      string
	Username string
	Password string
	Name     string
	Timeout  time.Duration
}

type Config struct {
	Debug     bool
	SecretKey string
	Database  Database

	dbClient *mongo.Client
}

func (c *Config) HasDatabase() bool {
	return c.dbClient != nil
}

func (c *Config) NewDatabaseClient() error {
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(c.Database.Uri).SetTimeout(c.Database.Timeout))
	if err != nil {
		return err
	}

	if err := client.Ping(context.Background(), nil); err != nil {
		return fmt.Errorf("database error: %s", err.Error())
	}
	c.dbClient = client
	return nil
}

func (c *Config) GetDatabaseClient() *mongo.Client {
	if c.HasDatabase() {
		return c.dbClient
	}
	return nil
}
