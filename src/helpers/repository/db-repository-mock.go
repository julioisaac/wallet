package repository

import "github.com/stretchr/testify/mock"

type MockDBRepository struct {
	mock.Mock
}

func (m *MockDBRepository) FindAll(Skip, Limit, sort int, query interface{}, objType interface{}) []interface{} {
	ret := m.Called(Skip, Limit, sort, query, objType)
	return ret
}

func (m *MockDBRepository) FindOne(query string, response interface{}) error {
	ret := m.Called(query, response)
	return ret.Error(0)
}

func (m *MockDBRepository) Insert(value interface{}) error {
	ret := m.Called(value)
	return ret.Error(0)}

func (m *MockDBRepository) Upsert(selector interface{}, update interface{}) error {
	ret := m.Called(selector, update)
	return ret.Error(0)
}

func (m *MockDBRepository) DeleteOne(key string, value interface{}) (int64, error) {
	ret := m.Called(key, value)
	return ret.Get(0).(int64), ret.Error(1)
}
