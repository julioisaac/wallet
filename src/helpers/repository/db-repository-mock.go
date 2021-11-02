package repository

import (
	"context"
	"github.com/stretchr/testify/mock"
)

type MockDBRepository struct {
	mock.Mock
}

func (m *MockDBRepository) FindAll(ctx context.Context, Skip, Limit, sort int, query interface{}, objType interface{}) []interface{} {
	ret := m.Called(ctx, Skip, Limit, sort, query, objType)
	return ret.Get(0).([]interface{})
}

func (m *MockDBRepository) FindOne(ctx context.Context, query string, response interface{}) error {
	ret := m.Called(ctx, query, response)
	return ret.Error(0)
}

func (m *MockDBRepository) Insert(ctx context.Context, value interface{}) error {
	ret := m.Called(ctx, value)
	return ret.Error(0)}

func (m *MockDBRepository) Upsert(ctx context.Context, selector interface{}, update interface{}) error {
	ret := m.Called(ctx, selector, update)
	return ret.Error(0)
}

func (m *MockDBRepository) DeleteOne(ctx context.Context, key string, value interface{}) (int64, error) {
	ret := m.Called(ctx, key, value)
	return ret.Get(0).(int64), ret.Error(1)
}
