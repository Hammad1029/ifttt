package store

import (
	"fmt"
	redisInfra "ifttt/manager/infrastructure/redis"
	"strconv"

	"github.com/go-redis/redis"
	"github.com/mitchellh/mapstructure"
)

const redisCache = "redis"

type RedisStore struct {
	client *redis.Client
	config redisConfig
}

type redisConfig struct {
	Db       string `json:"db" mapstructure:"db"`
	Host     string `json:"host" mapstructure:"host"`
	Port     string `json:"port" mapstructure:"port"`
	Password string `json:"password" mapstructure:"password"`
	DbIndex  string `json:"dbIndex" mapstructure:"dbIndex"`
}

func (r *RedisStore) init(config map[string]any) error {
	if err := mapstructure.Decode(config, &r.config); err != nil {
		return fmt.Errorf("method: *RedisStore.Init: could not decode redis configuration from env: %s", err)
	}
	dbIndex, _ := strconv.Atoi(r.config.DbIndex)
	r.client = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", r.config.Host, r.config.Port),
		Password: r.config.Password,
		DB:       dbIndex,
	})
	return nil
}

func (r *RedisStore) createCacheStore() *CacheStore {
	redisBase := redisInfra.NewRedisBaseRepository(r.client)
	return &CacheStore{
		Store:     r,
		TokenRepo: redisInfra.NewRedisTokenRepository(redisBase),
	}
}
