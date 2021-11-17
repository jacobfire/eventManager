package redisConfig

import (
	"github.com/go-redis/redis"
	"log"
)

type RedisStorage struct {
	config *Config
	Client *redis.Client
}

func New(config *Config) *RedisStorage {
	return &RedisStorage{
		config: config,
	}
}

func (rs *RedisStorage) InitRedis() error {
	client := redis.NewClient(&redis.Options{
		Addr: rs.config.Host + rs.config.Port,
		Password: rs.config.Password,
		DB: rs.config.DB,
	})
	rs.Client = client
	if err := client.Ping().Err(); err != nil {
		log.Println(err)
		return err
	}

	return nil
}
