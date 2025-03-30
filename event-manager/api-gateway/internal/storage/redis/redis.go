package redis

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/pkg/errors"
)

func NewClient(redisURL string) (*redis.Client, error) {
	opt, err := redis.ParseURL(redisURL)
	if err != nil {
		return nil, errors.WithMessage(err, "parse url")
	}

	client := redis.NewClient(opt)
	ctx := context.Background()

	if _, err := client.Ping(ctx).Result(); err != nil {
		return nil, errors.WithMessage(err, "ping redis")
	}

	return client, nil
}
