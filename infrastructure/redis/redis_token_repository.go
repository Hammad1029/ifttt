package redisInfra

import (
	"context"
	"encoding/json"
	"fmt"
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

const (
	accessTokenKey  = "user-session-at"
	refreshTokenKey = "user-session-rt"
)

func (r *RedisTokenRepository) StoreTokenPair(email string, tokens *auth.TokenPair, ctx context.Context) error {
	now := time.Now()
	atMarshalled, err := json.Marshal(tokens.AccessToken)
	if err != nil {
		return fmt.Errorf("method *RedisTokenRepository.StoreTokenPair: error in marshalling accessToken: %s", err)
	}
	rtMarshalled, err := json.Marshal(tokens.AccessToken)
	if err != nil {
		return fmt.Errorf("method *RedisTokenRepository.StoreTokenPair: error in marshalling refreshToken: %s", err)
	}

	if err := r.client.Set(
		ctx, fmt.Sprintf("%s:%s", accessTokenKey, email), atMarshalled, time.Unix(tokens.AccessToken.Expiry, 0).Sub(now),
	).Err(); err != nil {
		return fmt.Errorf("method *RedisTokenRepository.StoreTokenPair: error in setting access token: %s", err)
	}
	if err := r.client.Set(
		ctx, fmt.Sprintf("%s:%s", refreshTokenKey, email), rtMarshalled, time.Unix(tokens.RefreshToken.Expiry, 0).Sub(now),
	).Err(); err != nil {
		return fmt.Errorf("method *RedisTokenRepository.StoreTokenPair: error in setting refresh token: %s", err)
	}

	return nil
}

func (r *RedisTokenRepository) GetTokenPair(email string, ctx context.Context) (*auth.TokenPair, error) {
	var pair auth.TokenPair

	if accessToken, err := r.client.Get(ctx, fmt.Sprintf("%s:%s", accessTokenKey, email)).Result(); err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		return nil, fmt.Errorf("method *RedisTokenRepository.GetTokenPair: error in getting access token: %s", err)
	} else {
		if err := json.Unmarshal([]byte(accessToken), &pair.AccessToken); err != nil {
			return nil,
				fmt.Errorf("method *RedisTokenRepository.GetTokenPair: error in unmarshalling access token: %s", err)
		}
	}

	if refreshToken, err := r.client.Get(ctx, fmt.Sprintf("%s:%s", refreshTokenKey, email)).Result(); err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		return nil, fmt.Errorf("method *RedisTokenRepository.GetTokenPair: error in getting refresh token: %s", err)
	} else {
		if err := json.Unmarshal([]byte(refreshToken), &pair.RefreshToken); err != nil {
			return nil,
				fmt.Errorf("method *RedisTokenRepository.GetTokenPair: error in unmarshalling refresh token: %s", err)
		}
	}

	return &pair, nil
}

func (r *RedisTokenRepository) DeleteTokenPair(email string, ctx context.Context) error {
	if err := r.client.Del(ctx,
		fmt.Sprintf("%s:%s", accessTokenKey, email), fmt.Sprintf("%s:%s", refreshTokenKey, email),
	).Err(); err != nil {
		return fmt.Errorf("method *RedisTokenRepository.DeleteTokenPair: error in deleting token pair: %s", err)
	}
	return nil
}
