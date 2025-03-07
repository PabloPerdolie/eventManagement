package repository

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"time"
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
		return false, err
	}
	return result > 0, nil
}

func (r *Client) StoreResetToken(ctx context.Context, userID uuid.UUID, token string, expiry time.Duration) error {
	return r.redisDb.Set(ctx, "reset:"+token, userID.String(), expiry).Err()
}

func (r *Client) GetUserIDByResetToken(ctx context.Context, token string) (uuid.UUID, error) {
	result, err := r.redisDb.Get(ctx, "reset:"+token).Result()
	if err != nil {
		return uuid.Nil, err
	}

	return uuid.Parse(result)
}

func (r *Client) DeleteResetToken(ctx context.Context, token string) error {
	return r.redisDb.Del(ctx, "reset:"+token).Err()
}

func (r *Client) DeleteRefreshToken(ctx context.Context, token string) error {
	return r.redisDb.Del(ctx, "blacklist:"+token).Err()
}
