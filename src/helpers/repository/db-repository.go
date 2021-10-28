package repository

type DBRepository interface {
	FindAll(Skip, Limit, sort int, query interface{}, objType interface{}) []interface{}
	FindOne(query string, response interface{}) error
	Insert(value interface{}) error
	Upsert(selector interface{}, update interface{}) error
	DeleteOne(key string, value interface{}) (int64, error)
}