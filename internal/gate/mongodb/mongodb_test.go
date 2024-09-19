package mongodb

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"testing"

	"log/slog"
)

var (
	mongoConnectWrapper = mongo.Connect
	initWrapper         = Init
)

func TestInit(t *testing.T) {
	tests := []struct {
		name         string
		url          string
		wantErr      bool
		connectError error
		pingError    error
	}{
		{
			name:         "success",
			url:          "mongodb://localhost:27017",
			wantErr:      false,
			connectError: nil,
			pingError:    nil,
		},
		{
			name:         "invalid url",
			url:          "invalid_url",
			wantErr:      true,
			connectError: errors.New("invalid URI: must have a scheme"),
			pingError:    nil,
		},
		{
			name:         "ping error",
			url:          "mongodb://localhost:27018",
			wantErr:      true,
			connectError: nil,
			pingError:    errors.New("ping error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			log := slog.Default()

			originalMongoConnectWrapper := mongoConnectWrapper
			mongoConnectWrapper = func(ctx context.Context, opts ...*options.ClientOptions) (*mongo.Client, error) {
				return &mongo.Client{}, tt.connectError
			}
			defer func() { mongoConnectWrapper = originalMongoConnectWrapper }()

			originalInitWrapper := initWrapper
			initWrapper = func(ctx context.Context, log *slog.Logger, url string) (Handler, error) {
				client, err := mongoConnectWrapper(ctx, options.Client().ApplyURI(url))
				if err != nil {
					return Handler{}, err
				}
				if tt.pingError != nil {
					return Handler{}, tt.pingError
				}
				return Handler{client: client}, nil
			}
			defer func() { initWrapper = originalInitWrapper }()

			_, err := initWrapper(ctx, log, tt.url)
			if (err != nil) != tt.wantErr {
				t.Errorf("Init() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
