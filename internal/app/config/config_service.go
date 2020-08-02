package item

import (
	"meli/internal/entities"
	"meli/internal/redis"
)

type ConfigServicer interface {
	Get(id string) entities.Item
	Save(id string) error
}

type ConfigService struct {
	redis redis.Redis
}

func NewConfigService(redis redis.Redis) ConfigServicer {
	return &ConfigService{redis: redis}
}

func (s *ConfigService) Get(key string) entities.Item {

	return entities.Item{}
}

func (s *ConfigService) Save(id string) error {

	return nil
}
