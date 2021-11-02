package repository

import "context"

type DBRepository interface {
	FindAll(ctx context.Context, Skip, Limit, sort int, query interface{}, objType interface{}) []interface{}
	FindOne(ctx context.Context, query string, response interface{}) error
	Insert(ctx context.Context, value interface{}) error
	Upsert(ctx context.Context, selector interface{}, update interface{}) error
	DeleteOne(ctx context.Context, key string, value interface{}) (int64, error)
}