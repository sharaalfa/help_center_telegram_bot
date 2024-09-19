package redis

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

type MockInMemory struct {
	mock.Mock
}

func (m *MockInMemory) Set(ctx context.Context, chatId int64, key string, value interface{}) {
	m.Called(ctx, chatId, key, value)
}

func (m *MockInMemory) Get(ctx context.Context, chatId int64, key string) string {
	args := m.Called(ctx, chatId, key)
	return args.String(0)
}

func TestSet(t *testing.T) {
	mockClient := new(MockInMemory)

	handler := &Handler{client: mockClient}

	ctx := context.Background()

	mockClient.On("Set", mock.Anything, int64(123), "testKey", "testValue").Return()

	handler.Set(ctx, 123, "testKey", "testValue")

	mockClient.AssertExpectations(t)
}

func TestGet(t *testing.T) {
	mockClient := new(MockInMemory)

	handler := &Handler{client: mockClient}

	ctx := context.Background()

	expectedValue := "testValue"
	mockClient.On("Get", mock.Anything, int64(123), "testKey").Return(expectedValue)

	result := handler.Get(ctx, 123, "testKey")

	assert.Equal(t, expectedValue, result)

	mockClient.AssertExpectations(t)
}
