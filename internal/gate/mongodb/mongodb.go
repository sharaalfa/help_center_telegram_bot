package mongodb

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log/slog"
)

type Handler struct {
	client *mongo.Client
}

func Init(ctx context.Context, log *slog.Logger, url string) (Handler, error) {
	clientOptions := options.Client().ApplyURI(url)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return Handler{}, err
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		return Handler{}, err
	}

	log.Info("Connected to MongoDB")
	return Handler{client: client}, nil
}
