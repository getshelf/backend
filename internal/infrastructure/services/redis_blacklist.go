package services

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/getshelf/backend/internal/config"
)

type RedisBlacklist struct {
	client *redis.Client
}

func NewRedisBlacklist(cfg config.Config) RedisBlacklist {
	return RedisBlacklist{client: redis.NewClient(&redis.Options{Addr: cfg.Redis.Address()})}
}

func (blacklist RedisBlacklist) Revoke(ctx context.Context, tokenID string, ttl time.Duration) error {
	if ttl <= 0 {
		return nil
	}
	return blacklist.client.Set(ctx, "auth:revoked:"+tokenID, "1", ttl).Err()
}

func (blacklist RedisBlacklist) IsRevoked(ctx context.Context, tokenID string) (bool, error) {
	result, err := blacklist.client.Exists(ctx, "auth:revoked:"+tokenID).Result()
	return result > 0, err
}
