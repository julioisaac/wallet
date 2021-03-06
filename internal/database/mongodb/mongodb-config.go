package mongodb

import (
	"context"
	"github.com/julioisaac/daxxer-api/internal/database"
	"github.com/julioisaac/daxxer-api/internal/logs"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
	"os"
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
	URI := os.Getenv("MONGODB_URI")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(URI).SetMaxPoolSize(20))
	if err != nil {
		logs.Instance.Log.Error(context.Background(), "error trying to create mongodb connection", zap.Error(err))
	}
	logs.Instance.Log.Debug(context.Background(), "mongodb connection created")
	return client
}