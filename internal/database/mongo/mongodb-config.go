package mongo

import (
	"context"
	"fmt"
	"github.com/julioisaac/daxxer-api/internal/database"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type mongoConfig struct{}

type Database struct {
	Mongo *mongo.Client
}

var DB *Database

func NewMongoConfig() database.DBConfig {
	return &mongoConfig{}
}

func (m *mongoConfig) Init() {
	DB = &Database{
		Mongo: m.getConnect(),
	}
}

func (m *mongoConfig) getConnect() *mongo.Client {
	URI := "mongodb://mongo:27017"
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, Err := mongo.Connect(ctx, options.Client().ApplyURI(URI).SetMaxPoolSize(20))
	if Err != nil {
		fmt.Println(Err)
	}
	return client
}