package redisInfra

import "github.com/go-redis/redis"

type RedisBaseRepository struct {
	client *redis.Client
}

func NewRedisBaseRepository(client *redis.Client) *RedisBaseRepository {
	if client == nil {
		panic("missing redis client")
	}
	return &RedisBaseRepository{client: client}
}
