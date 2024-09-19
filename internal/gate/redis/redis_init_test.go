package redis

import (
	"context"
	"errors"
	"github.com/go-redis/redis/v8"
	"github.com/go-redis/redismock/v8"
	"log/slog"
	"testing"
)

var (
	newClientWrapper = redis.NewClient
	initWrapper      = Init
	parseURLWrapper  = redis.ParseURL
)

func TestInit(t *testing.T) {
	ctx := context.Background()
	log := slog.Default()

	db, mock := redismock.NewClientMock()

	originalNewClientWrapper := newClientWrapper
	newClientWrapper = func(opt *redis.Options) *redis.Client {
		return db
	}
	defer func() { newClientWrapper = originalNewClientWrapper }()

	originalInitWrapper := initWrapper
	initWrapper = func(ctx context.Context, log *slog.Logger, url string) (Handler, *redis.Client, error) {
		if url == "invalid_url" {
			return Handler{}, nil, errors.New("invalid redis URL scheme: invalid_url")
		}
		return Handler{client: &redisClientWrapper{client: newClientWrapper(&redis.Options{Addr: url})}}, nil, nil
	}
	defer func() { initWrapper = originalInitWrapper }()

	tests := []struct {
		name            string
		url             string
		expectedHandler Handler
		expectedErr     error
		pingExpectation func(mock redismock.ClientMock)
		parseURLError   error
	}{
		{
			name:            "success to connect to redis",
			url:             "redis://localhost:6379",
			expectedHandler: Handler{client: &redisClientWrapper{client: newClientWrapper(&redis.Options{Addr: "localhost:6379"})}},
			expectedErr:     nil,
			pingExpectation: func(mock redismock.ClientMock) {
				mock.ExpectPing().SetVal("PONG")
			},
			parseURLError: nil,
		},
		{
			name:            "failed to parse redis url",
			url:             "invalid_url",
			expectedHandler: Handler{},
			expectedErr:     errors.New("invalid redis URL scheme: invalid_url"),
			pingExpectation: func(mock redismock.ClientMock) {
				// No ping expectation for this test case
			},
			parseURLError: errors.New("invalid redis URL scheme: invalid_url"),
		},
		{
			name:            "failed to ping redis",
			url:             "redis://localhost:6379",
			expectedHandler: Handler{},
			expectedErr:     errors.New("failed to ping redis"),
			pingExpectation: func(mock redismock.ClientMock) {
				// No ping expectation for this test case
			},
			parseURLError: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler, _, err := initWrapper(ctx, log, tt.url)
			originalParseURLWrapper := parseURLWrapper
			parseURLWrapper = func(url string) (*redis.Options, error) {
				return nil, tt.parseURLError
			}
			defer func() { parseURLWrapper = originalParseURLWrapper }()
			if err != nil {
				if tt.expectedErr == nil || err.Error() != tt.expectedErr.Error() {
					t.Errorf("Init() error=%v, wantErr %v", err, tt.expectedErr)
					return
				}
			} else {
				if tt.expectedErr != nil && tt.parseURLError != nil {
					t.Errorf("Init() expected error=%v, but got none", tt.expectedErr)
					return
				}
			}

			if handler.client != tt.expectedHandler.client && tt.parseURLError != nil {
				t.Errorf("expected client to be %v, but got %v", tt.expectedHandler.client, handler.client)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}
