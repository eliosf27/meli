package redis

import (
	"fmt"
	"github.com/go-redis/redis"
	log "github.com/sirupsen/logrus"
	config "meli/pkg/config"
	"time"
)

type Redis struct {
	client *redis.Client
	config config.Config
}

func NewRedis(config config.Config) Redis {
	client := buildClient(config)
	if client == nil {
		panic(fmt.Sprintf("Error connecting to redis server"))
	}

	status := client.Ping()
	if status.Err() != nil {
		log.Info("Redis => Monitoring | Cannot connect to redis server | Error => ", status.Err())
	} else {
		log.Info("Redis => Monitoring | Connected successfully")
	}

	return Redis{
		client: client,
		config: config,
	}
}

func buildClient(config config.Config) *redis.Client {
	var url = fmt.Sprintf("%s:%s", config.Redis.RedisHost, config.Redis.RedisPort)
	var options = &redis.Options{
		Addr:     url,
		Password: "", // no password set
		DB:       0,  // use default DB
		OnConnect: func(*redis.Conn) error {
			return nil
		},
		DialTimeout:  10 * time.Second,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		PoolSize:     10,
		PoolTimeout:  30 * time.Second,
	}

	return redis.NewClient(options)
}

func (r *Redis) Set(key string, data []byte, duration time.Duration) ([]byte, error) {
	status := r.client.Set(key, data, duration)
	if status.Err() != nil {
		return []byte{}, status.Err()
	}

	return []byte(status.Val()), nil
}

func (r *Redis) Get(key string) ([]byte, error) {
	status := r.client.Get(key)
	if status.Err() != nil {
		return []byte{}, status.Err()
	}

	return []byte(status.Val()), nil
}

func (r *Redis) HSet(key string, field string, data []byte) (bool, error) {
	status := r.client.HSet(key, field, data)
	if status.Err() != nil {
		return false, status.Err()
	}

	return status.Val(), nil
}

func (r *Redis) HGet(key string, field string) ([]byte, error) {
	status := r.client.HGet(key, field)
	if status.Err() != nil {
		return []byte{}, status.Err()
	}

	return []byte(status.Val()), nil
}

func (r *Redis) HGetAll(key string) (map[string]string, error) {
	status := r.client.HGetAll(key)
	if status.Err() != nil {
		return map[string]string{}, status.Err()
	}

	return status.Val(), nil
}
