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

func (r *MongoRepository) Insert(ctx context.Context, value interface{}) error {
	client := mongodb.DB.Mongo
	collection := client.Database(r.database).Collection(r.collection)
	_, err := collection.InsertOne(ctx, value)
	if err != nil {
		logs.Instance.Log.Error(ctx, "error inserting in mongodb", zap.Error(err))
	}
	logs.Instance.Log.Debug(ctx, "successfully inserted into mongodb")
	return nil
}

func (r *MongoRepository) Upsert(ctx context.Context, selector interface{}, update interface{}) error {
	client := mongodb.DB.Mongo
	opts := options.Update().SetUpsert(true)
	collection := client.Database(r.database).Collection(r.collection)
	_, err := collection.UpdateOne(ctx, selector, update, opts)
	if err != nil {
		logs.Instance.Log.Error(ctx, "error upserting in mongodb", zap.Error(err))
		return err
	}
	logs.Instance.Log.Debug(ctx, "successfully upserted into mongodb")
	return nil
}

func (r *MongoRepository) DeleteOne(ctx context.Context, key string, value interface{}) (int64, error) {
	client := mongodb.DB.Mongo
	collection := client.Database(r.database).Collection(r.collection)
	filter := bson.D{{key, value}}
	count, err := collection.DeleteOne(ctx, filter, nil)
	if err != nil {
		logs.Instance.Log.Error(ctx, "error deleting in mongodb", zap.Error(err))
		return 0, err
	}
	logs.Instance.Log.Debug(ctx, "successfully deleted from mongodb")
	return count.DeletedCount, err
}

func (r *MongoRepository) FindAll(ctx context.Context, Skip, Limit, sort int, query interface{}, objType interface{}) []interface{} {
	var responses = make([]interface{}, 0)
	objectType := reflect.TypeOf(objType).Elem()

	client := mongodb.DB.Mongo
	collection := client.Database(r.database).Collection(r.collection)
	SORT := bson.D{{"_id", sort}}
	findOptions := options.Find().SetSort(SORT).SetLimit(int64(Limit)).SetSkip(int64(Skip))
	sortCursor, err := collection.Find(ctx, query, findOptions)
	if err != nil {
		logs.Instance.Log.Error(ctx, "error finding all in mongodb", zap.Error(err))
		return responses
	}

	defer sortCursor.Close(ctx)
	for sortCursor.Next(ctx) {
		result := reflect.New(objectType).Interface()
		err1 := sortCursor.Decode(result)
		if err1 != nil {
			logs.Instance.Log.Error(ctx, "error decoding data from mongodb", zap.Error(err1))
			return nil
		}
		responses = append(responses, result)
	}
	logs.Instance.Log.Debug(ctx, "successfully found all from mongodb")
	return responses
}

func (r *MongoRepository) FindOne(ctx context.Context, query string, response interface{}) error {
	client := mongodb.DB.Mongo
	collection, err := client.Database(r.database).Collection(r.collection).Clone()
	if err != nil {
		return err
	}
	var filter interface{}
	err = bson.UnmarshalExtJSON([]byte(query), true, &filter)
	if err != nil {
		logs.Instance.Log.Error(ctx, "error converting query to mongodb", zap.Error(err))
		return err
	}
	if err = collection.FindOne(ctx, filter).Decode(response); err != nil {
		logs.Instance.Log.Error(ctx, "error finding from mongodb", zap.Error(err))
		return err
	}
	logs.Instance.Log.Debug(ctx, "successfully found one from mongodb")
	return nil
}