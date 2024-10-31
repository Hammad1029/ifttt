package redisInfra

import (
	"context"
	"encoding/json"
	"fmt"
	"ifttt/manager/common"
	"ifttt/manager/domain/auth"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisTokenRepository struct {
	*RedisBaseRepository
}

func NewRedisTokenRepository(base *RedisBaseRepository) *RedisTokenRepository {
	return &RedisTokenRepository{RedisBaseRepository: base}
}

func (r *RedisTokenRepository) StoreTokenPair(email string, tokens *auth.TokenPair, ctx context.Context) error {
	now := time.Now()
	atMarshalled, err := json.Marshal(tokens.Access)
	if err != nil {
		return fmt.Errorf("method *RedisTokenRepository.StoreTokenPair: error in marshalling accessToken: %s", err)
	}
	rtMarshalled, err := json.Marshal(tokens.Refresh)
	if err != nil {
		return fmt.Errorf("method *RedisTokenRepository.StoreTokenPair: error in marshalling refreshToken: %s", err)
	}

	if err := r.client.Set(
		ctx, fmt.Sprintf("%s:%s", common.AccessTokenKey, email), atMarshalled, time.Unix(tokens.Access.Expiry, 0).Sub(now),
	).Err(); err != nil {
		return err
	}
	if err := r.client.Set(
		ctx, fmt.Sprintf("%s:%s", common.RefreshTokenKey, email), rtMarshalled, time.Unix(tokens.Refresh.Expiry, 0).Sub(now),
	).Err(); err != nil {
		return err
	}

	return nil
}

func (r *RedisTokenRepository) GetTokenPair(email string, ctx context.Context) (*auth.TokenPair, error) {
	var pair auth.TokenPair

	if accessToken, err := r.client.Get(ctx, fmt.Sprintf("%s:%s", common.AccessTokenKey, email)).Result(); err != redis.Nil {
		if err != nil {
			return nil, err
		} else if err := json.Unmarshal([]byte(accessToken), &pair.Access); err != nil {
			return nil, err
		}

	}
	if refreshToken, err := r.client.Get(ctx, fmt.Sprintf("%s:%s", common.RefreshTokenKey, email)).Result(); err != redis.Nil {
		if err != nil {
			return nil, err
		} else if err := json.Unmarshal([]byte(refreshToken), &pair.Refresh); err != nil {
			return nil, err
		}
	}

	return &pair, nil
}

func (r *RedisTokenRepository) DeleteTokenPair(email string, ctx context.Context) error {
	if err := r.client.Del(ctx,
		fmt.Sprintf("%s:%s", common.AccessTokenKey, email), fmt.Sprintf("%s:%s", common.RefreshTokenKey, email),
	).Err(); err != nil {
		return err
	}
	return nil
}
