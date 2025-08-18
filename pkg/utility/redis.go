package utility

import (
	shared "provider-report-api/internal/modules/shared/dtos"
	"provider-report-api/pkg/vault"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

// Remove duplicate struct - use shared.UserActionAccessRights instead

type RedisService struct {
	client *redis.Client
}

func (s *RedisService) GetSearchPreference(subModule, username string) (interface{}, error) {
	key := fmt.Sprintf("save:search:%s:%s", subModule, username)
	fmt.Println("key", key)
	value, err := s.client.Get(context.Background(), key).Result()
	if err != nil {
	 return nil, err
	}
   
	var reqSearch interface{}
	err = json.Unmarshal([]byte(value), &reqSearch)
	if err != nil {
	 return nil, err
	}
   
	return reqSearch, nil
}

func (s *RedisService) SetSearchPreference(subModule, username string, reqSearch interface{}) error {
	key := fmt.Sprintf("save:search:%s:%s", subModule, username)
	value, err := json.Marshal(reqSearch)
	if err != nil {
	 return err
	}
   
	err = s.client.Set(context.Background(), key, value, 0).Err()
   
	return err
}
	

func NewRedisService() *RedisService {
	creds := vault.GetRedisSecret()
	if creds == nil {
		log.Fatalf("Secret data type is not expected.")
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     creds.RedisUrl,
		Password: "", // no password set
		DB:       0,
	})

	// Ping the Redis server to check if the connection is established
	pong, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		log.Fatal("Failed to connect to Redis:", err)
	}
	fmt.Println("Connected to Redis:", pong)

	return &RedisService{
		client: rdb,
	}
}

func (s *RedisService) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	err := s.client.Set(ctx, key, value, expiration).Err()
	return err
}

func (s *RedisService) Get(ctx context.Context, key string) (string, error) {
	value, err := s.client.Get(ctx, key).Result()
	if err != nil {
		return "", err
	}
	return value, nil
}

func (s *RedisService) SetPermissions(userId string, userRoleId string, moduleName string, permissions interface{}) error {
	key := fmt.Sprintf("permissions:%s:%s:%s", userId, userRoleId, moduleName)
	value, err := json.Marshal(permissions)
	if err != nil {
		return err
	}

	err = s.client.Set(context.Background(), key, value, 0).Err()
	return err
}

func (s *RedisService) GetPermissions(userId string, userRoleId string, moduleName string) (*[]shared.UserActionAccessRights, error) {
	key := fmt.Sprintf("permissions:%s:%s:%s", userId, userRoleId, moduleName)
	fmt.Println(key)

	value, err := s.client.Get(context.Background(), key).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, err // Key not found
		}
		return nil, err
	}

	var permissions *[]shared.UserActionAccessRights
	err = json.Unmarshal([]byte(value), &permissions)
	if err != nil {
		return nil, err
	}
	return permissions, nil
}

func (s *RedisService) FetchKeysByPattern(pattern string) ([]string, error) {
	var keys []string
	var cursor uint64
	var err error

	for {
		var scanResult []string
		scanResult, cursor, err = s.client.Scan(context.Background(), cursor, pattern, 100).Result()
		if err != nil {
			return nil, err
		}

		// Append keys to the result
		keys = append(keys, scanResult...)
		if cursor == 0 {
			break
		}
	}

	return keys, nil
}

func (s *RedisService) InvalidateUserPermissions(userId string) error {
	pattern := fmt.Sprintf("permissions:%s:*", userId)
	keys, err := s.FetchKeysByPattern(pattern)
	if err != nil {
		return err
	}

	if len(keys) > 0 {
		err := s.client.Del(context.Background(), keys...).Err()
		if err != nil {
			return err
		}
	}

	return nil
}