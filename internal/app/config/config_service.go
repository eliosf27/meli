package config

import (
	"errors"
	log "github.com/sirupsen/logrus"
	"meli/internal/redis"
	"time"
)

const (
	CacheStorageKey    = "meli:cache_storage"
	RedisStorage       = "redis"
	PostgresStorage    = "postgres"
	CacheStorageKeyTtl = time.Minute * 60
)

type ConfigServicer interface {
	Fetch() (string, error)
	Update(storage string) error
}

type ConfigService struct {
	redis redis.Redis
}

func NewConfigService(redis redis.Redis) ConfigServicer {
	return &ConfigService{redis: redis}
}

func (s *ConfigService) Fetch() (string, error) {
	val, err := s.redis.Get(CacheStorageKey)
	if !s.redis.Exist(err) {

		return "", errors.New("configuracion not found")
	}

	if err != nil {
		log.Errorf("ConfigService.Fetch | error fetching storage cache [%s] - %+v", CacheStorageKey, err)

		return "", err
	}

	log.Errorf("ConfigService.Fetch | fetch storage: %s", string(val))

	return string(val), nil
}

func (s *ConfigService) Update(storage string) error {
	_, err := s.redis.Set(CacheStorageKey, []byte(storage), CacheStorageKeyTtl)
	if err != nil {
		log.Errorf("ConfigService.Update |error saving storage cache [%s] - %+v", CacheStorageKey, err)

		return err
	}

	log.Errorf("ConfigService.Update | update storage: %s", storage)

	return nil
}
