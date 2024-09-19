package redis

import (
	"context"
	"github.com/go-redis/redis/v8"
	"log/slog"
)

type Handler struct {
	client InMemory
}

func Init(ctx context.Context, log *slog.Logger, url string) (Handler, *redis.Client, error) {
	opt, err := redis.ParseURL(url)
	if err != nil {
		return Handler{}, nil, err
	}

	client := redis.NewClient(opt)
	_, err = client.Ping(ctx).Result()
	if err != nil {
		return Handler{}, nil, err
	}

	log.Info("Connected to Redis")
	return Handler{client: &redisClientWrapper{client: client}}, client, err
}
