package mongodb

import (
	"context"
	"github.com/julioisaac/daxxer-api/src/helpers/repository"
	"github.com/julioisaac/daxxer-api/storage/mongodb"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"reflect"
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
	client := mongodb.DB.Mongo
	collection := client.Database(r.database).Collection(r.collection)
	filter := bson.D{{key, value}}
	count, err := collection.DeleteOne(context.TODO(), filter, nil)
	if err != nil {
		log.Println(err)
	}
	return count.DeletedCount, err
}

func (r *repo) FindAll(Skip, Limit, sort int, query interface{}, objType interface{}) []interface{} {
	var responses = make([]interface{}, 0)
	objectType := reflect.TypeOf(objType).Elem()

	client := mongodb.DB.Mongo
	collection := client.Database(r.database).Collection(r.collection)
	SORT := bson.D{{"_id", sort}}
	findOptions := options.Find().SetSort(SORT).SetLimit(int64(Limit)).SetSkip(int64(Skip))
	sortCursor, err := collection.Find(context.TODO(), query, findOptions)
	if err != nil {
		return responses
	}

	defer sortCursor.Close(context.TODO())
	for sortCursor.Next(context.TODO()) {
		result := reflect.New(objectType).Interface()
		err := sortCursor.Decode(result)
		if err != nil {
			log.Println(err)
			return nil
		}
		responses = append(responses, result)
	}
	return responses
}

func (r *repo) FindOne(query string, response interface{}) error {
	client := mongodb.DB.Mongo
	collection, _ := client.Database(r.database).Collection(r.collection).Clone()
	var filter interface{}
	err := bson.UnmarshalExtJSON([]byte(query), true, &filter)
	if err != nil {
		return err
	}
	if err := collection.FindOne(context.TODO(), filter).Decode(response); err != nil {
		return err
	}
	return nil
}