package redis

import (
	"context"
	"github.com/go-redis/redis/v8"
	"time"
)

type InMemory interface {
	Set(ctx context.Context, chatId int64, key string, value interface{})
	Get(ctx context.Context, chatId int64, key string) string
}

func (h *Handler) Set(ctx context.Context, chatId int64, key string, value interface{}) {
	ctxTimeout, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()
	h.client.Set(ctxTimeout, chatId, key, value)
}

func (h *Handler) Get(ctx context.Context, chatId int64, key string) string {
	ctxTimeout, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()
	return h.client.Get(ctxTimeout, chatId, key)
}

type redisClientWrapper struct {
	client *redis.Client
}

func (r redisClientWrapper) Set(ctx context.Context, chatId int64, key string, value interface{}) {
	ctxTimeout, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()
	r.client.Set(ctxTimeout, key, value, 0)
}

func (r redisClientWrapper) Get(ctx context.Context, chatId int64, key string) string {
	//TODO implement me
	ctxTimeout, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()
	res, _ := r.client.Get(ctxTimeout, key).Result()
	return res
}
