package redisInfra

import "github.com/redis/go-redis/v9"

type RedisBaseRepository struct {
	client *redis.Client
}

func NewRedisBaseRepository(client *redis.Client) *RedisBaseRepository {
	if client == nil {
		panic("missing redis client")
	}
	return &RedisBaseRepository{client: client}
}
