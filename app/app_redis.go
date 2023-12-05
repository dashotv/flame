package app

import (
	"context"

	"github.com/go-redis/redis/v8"
	"github.com/pkg/errors"
)

var cache *redis.Client

func setupRedis() error {
	log.Infof("connecting redis: %s", cfg.RedisURL())
	cache = redis.NewClient(&redis.Options{
		Addr: cfg.RedisURL(),
		DB:   cfg.RedisDatabase,
	})
	status := cache.Ping(context.Background())
	if status.Err() != nil {
		return errors.Errorf("failed to connect to redis: %s", status.Err())
	}
	return nil
}
