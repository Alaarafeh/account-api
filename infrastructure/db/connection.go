package db

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	userCollection        = "users"
)

type Connection struct {
	client    *mongo.Client
	defaultDB string
}

func NewConnection(config *Config) (*Connection, error) {
	client, err := mongo.NewClient(options.Client().ApplyURI(config.DBURL))
	if err != nil {
		return nil, err
	}

	dbConnectionTimeout := 10 * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), dbConnectionTimeout)
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		return nil, err
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, err
	}

	return &Connection{client: client, defaultDB: config.DBName}, nil
}

func (c *Connection) GetUsers() *mongo.Collection {
	return c.selectDefaultDB().Collection(userCollection)
}

func (c *Connection) DropDB() error {
	return c.selectDefaultDB().Drop(context.TODO())
}

func (c *Connection) selectDefaultDB() *mongo.Database {
	return c.client.Database(c.defaultDB)
}
