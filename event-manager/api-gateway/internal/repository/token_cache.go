package repository

import (
	"context"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/pkg/errors"
)

type Client struct {
	redisDb *redis.Client
}

func NewClient(client *redis.Client) Client {
	return Client{
		redisDb: client,
	}
}

func (r *Client) BlacklistToken(ctx context.Context, token string, expiry time.Duration) error {
	return r.redisDb.Set(ctx, "blacklist:"+token, true, expiry).Err()
}

func (r *Client) IsTokenBlacklisted(ctx context.Context, token string) (bool, error) {
	result, err := r.redisDb.Exists(ctx, "blacklist:"+token).Result()
	if err != nil {
		return false, errors.WithMessage(err, "get token blacklist")
	}
	return result > 0, nil
}

func (r *Client) StoreResetToken(ctx context.Context, userId int, token string, expiry time.Duration) error {
	return r.redisDb.Set(ctx, "reset:"+token, userId, expiry).Err()
}

func (r *Client) GetUserIdByResetToken(ctx context.Context, token string) (int, error) {
	result, err := r.redisDb.Get(ctx, "reset:"+token).Result()
	if err != nil {
		return 0, errors.WithMessage(err, "GetUserIdByResetToken")
	}

	res, err := strconv.Atoi(result)
	if err != nil {
		return 0, errors.WithMessage(err, "parse token")
	}

	return res, nil
}

func (r *Client) DeleteResetToken(ctx context.Context, token string) error {
	return r.redisDb.Del(ctx, "reset:"+token).Err()
}

func (r *Client) DeleteRefreshToken(ctx context.Context, token string) error {
	return r.redisDb.Del(ctx, "blacklist:"+token).Err()
}
