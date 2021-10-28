package mongodb

import (
	"context"
	"github.com/julioisaac/daxxer-api/src/helpers/repository"
	"github.com/julioisaac/daxxer-api/storage/mongodb"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type repo struct {
	database   string
	collection string
}

func NewMongodbRepository(database, collection string) repository.DBRepository {
	return &repo{database, collection}
}

func (r *repo) Insert(value interface{}) error {
	client := mongodb.DB.Mongo
	collection := client.Database(r.database).Collection(r.collection)
	_, err := collection.InsertOne(context.TODO(), value)
	if err != nil {
		return err
	}
	return nil
}

func (r *repo) Upsert(selector interface{}, update interface{}) error {
	client := mongodb.DB.Mongo
	opts := options.Update().SetUpsert(true)
	collection := client.Database(r.database).Collection(r.collection)
	_, err := collection.UpdateOne(context.TODO(), selector, update, opts)
	if err != nil {
		return err
	}
	return nil
}

func (r *repo) DeleteOne(key string, value interface{}) (int64, error) {
	panic("implement me")
}

func (r *repo) FindAll(Skip, Limit, sort int, query interface{}, objType interface{}) []interface{} {
	panic("implement me")
}

func (r *repo) FindOne(query string, response interface{}) error {
	panic("implement me")
}