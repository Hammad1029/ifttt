package redisInfra

import (
	"encoding/json"
	"fmt"
	"ifttt/manager/domain/token"
	"time"

	"github.com/go-redis/redis"
)

type RedisTokenRepository struct {
	*RedisBaseRepository
}

func NewRedisTokenRepository(base *RedisBaseRepository) *RedisTokenRepository {
	return &RedisTokenRepository{RedisBaseRepository: base}
}

func (r *RedisTokenRepository) StoreTokenPair(email string, tokens *token.TokenPair) error {
	now := time.Now()
	atMarshalled, err := json.Marshal(tokens.AccessToken)
	if err != nil {
		return fmt.Errorf("method *RedisTokenRepository.StoreTokenPair: erorr in marshalling accessToken: %s", err)
	}
	rtMarshalled, err := json.Marshal(tokens.AccessToken)
	if err != nil {
		return fmt.Errorf("method *RedisTokenRepository.StoreTokenPair: erorr in marshalling refreshToken: %s", err)
	}

	r.client.Set(fmt.Sprintf("%s-at", email), atMarshalled, now.Sub(time.Unix(tokens.AccessToken.Expiry, 0)))
	r.client.Set(fmt.Sprintf("%s-rt", email), rtMarshalled, now.Sub(time.Unix(tokens.RefreshToken.Expiry, 0)))

	return nil
}

func (r *RedisTokenRepository) GetTokenPair(email string) (*token.TokenPair, error) {
	var pair token.TokenPair

	if accessToken, err := r.client.Get(fmt.Sprintf("%s-at", email)).Result(); err != nil {
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

	if refreshToken, err := r.client.Get(fmt.Sprintf("%s-rt", email)).Result(); err != nil {
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

func (r *RedisTokenRepository) DeleteTokenPair(email string) error {
	if err := r.client.Del(
		fmt.Sprintf("%s-at", email), fmt.Sprintf("%s-rt", email),
	).Err(); err != nil {
		return fmt.Errorf("method *RedisTokenRepository.DeleteTokenPair: error in deleting token pair: %s", err)
	}
	return nil
}
