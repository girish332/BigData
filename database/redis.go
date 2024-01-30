package database

import (
	"github.com/gin-gonic/gin"
	goredis "github.com/redis/go-redis/v9"
	"time"
)

type RedisRepo struct {
	client goredis.Client
}

func NewRedisRepo(address, password string) *RedisRepo {
	return &RedisRepo{
		client: *goredis.NewClient(&goredis.Options{
			Addr:     address,
			Password: password,
		}),
	}
}

func (repo *RedisRepo) Get(ctx *gin.Context, key string) (string, error) {
	val, err := repo.client.Get(ctx, key).Result()
	if err != nil {
		return "", err
	}
	return val, nil
}

func (repo *RedisRepo) Set(ctx *gin.Context, key, value string) error {
	_, err := repo.client.Set(ctx, key, value, 5*time.Hour).Result()
	if err != nil {
		return err
	}
	return nil
}

func (repo *RedisRepo) Ping(ctx *gin.Context) error {
	_, err := repo.client.Ping(ctx).Result()
	if err != nil {
		return err
	}
	return nil
}

func (repo *RedisRepo) Delete(ctx *gin.Context, key string) error {
	_, err := repo.client.Del(ctx, key).Result()
	if err != nil {
		return err
	}
	return nil
}

func (repo *RedisRepo) Keys(c *gin.Context, pattern string) ([]string, error) {
	keys, err := repo.client.Keys(c, pattern).Result()
	if err != nil {
		return nil, err
	}
	return keys, nil
}
