package mongodb

import (
	"context"
	"github.com/julioisaac/daxxer-api/internal/database/mongodb"
	"github.com/julioisaac/daxxer-api/internal/logs"
	"github.com/julioisaac/daxxer-api/src/helpers/repository"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
	"reflect"
)

type MongoRepository struct {
	database   string
	collection string
}

func NewMongodbRepository(database, collection string) repository.DBRepository {
	return &MongoRepository{database, collection}
}

func (r *MongoRepository) Insert(value interface{}) error {
	client := mongodb.DB.Mongo
	collection := client.Database(r.database).Collection(r.collection)
	_, err := collection.InsertOne(context.TODO(), value)
	if err != nil {
		logs.Instance.Log.Error(context.Background(), "error inserting in mongodb", zap.Error(err))
	}
	logs.Instance.Log.Debug(context.Background(), "successfully inserted into mongodb")
	return nil
}

func (r *MongoRepository) Upsert(selector interface{}, update interface{}) error {
	client := mongodb.DB.Mongo
	opts := options.Update().SetUpsert(true)
	collection := client.Database(r.database).Collection(r.collection)
	_, err := collection.UpdateOne(context.TODO(), selector, update, opts)
	if err != nil {
		logs.Instance.Log.Error(context.Background(), "error upserting in mongodb", zap.Error(err))
		return err
	}
	logs.Instance.Log.Debug(context.Background(), "successfully upserted into mongodb")
	return nil
}

func (r *MongoRepository) DeleteOne(key string, value interface{}) (int64, error) {
	client := mongodb.DB.Mongo
	collection := client.Database(r.database).Collection(r.collection)
	filter := bson.D{{key, value}}
	count, err := collection.DeleteOne(context.TODO(), filter, nil)
	if err != nil {
		logs.Instance.Log.Error(context.Background(), "error deleting in mongodb", zap.Error(err))
		return 0, err
	}
	logs.Instance.Log.Debug(context.Background(), "successfully deleted from mongodb")
	return count.DeletedCount, err
}

func (r *MongoRepository) FindAll(Skip, Limit, sort int, query interface{}, objType interface{}) []interface{} {
	var responses = make([]interface{}, 0)
	objectType := reflect.TypeOf(objType).Elem()

	client := mongodb.DB.Mongo
	collection := client.Database(r.database).Collection(r.collection)
	SORT := bson.D{{"_id", sort}}
	findOptions := options.Find().SetSort(SORT).SetLimit(int64(Limit)).SetSkip(int64(Skip))
	sortCursor, err := collection.Find(context.TODO(), query, findOptions)
	if err != nil {
		logs.Instance.Log.Error(context.Background(), "error finding all in mongodb", zap.Error(err))
		return responses
	}

	defer sortCursor.Close(context.TODO())
	for sortCursor.Next(context.TODO()) {
		result := reflect.New(objectType).Interface()
		err1 := sortCursor.Decode(result)
		if err1 != nil {
			logs.Instance.Log.Error(context.Background(), "error decoding data from mongodb", zap.Error(err1))
			return nil
		}
		responses = append(responses, result)
	}
	logs.Instance.Log.Debug(context.Background(), "successfully found all from mongodb")
	return responses
}

func (r *MongoRepository) FindOne(query string, response interface{}) error {
	client := mongodb.DB.Mongo
	collection, _ := client.Database(r.database).Collection(r.collection).Clone()
	var filter interface{}
	err := bson.UnmarshalExtJSON([]byte(query), true, &filter)
	if err != nil {
		logs.Instance.Log.Error(context.Background(), "error converting query to mongodb", zap.Error(err))
		return err
	}
	if err1 := collection.FindOne(context.TODO(), filter).Decode(response); err != nil {
		logs.Instance.Log.Error(context.Background(), "error finding from mongodb", zap.Error(err1))
		return err1
	}
	logs.Instance.Log.Debug(context.Background(), "successfully found one from mongodb")
	return nil
}