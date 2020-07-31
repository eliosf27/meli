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
