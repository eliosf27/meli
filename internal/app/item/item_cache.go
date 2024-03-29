package item

import (
	"errors"
	"meli/internal/entities"
)

const (
	RedisCache    = "redis"
	PostgresCache = "postgres"
)

type ItemCacher interface {
	Get(Id string) (entities.Item, error)
	Save(item entities.Item) error
	SetStrategy(strategy Cacher)
	GetStrategy(storage string) Cacher
}

type Cacher interface {
	Get(Id string) (entities.Item, error)
	Save(item entities.Item) error
	MustApply(strategy string) bool
	Name() string
}

type ItemCache struct {
	strategy      Cacher
	redisCache    Cacher
	postgresCache Cacher
}

func NewItemCache(redisCache Cacher, postgresCache Cacher) *ItemCache {
	return &ItemCache{redisCache: redisCache, postgresCache: postgresCache}
}

func (c *ItemCache) Get(id string) (entities.Item, error) {
	if c.strategy != nil {
		return c.strategy.Get(id)
	}

	return entities.Item{}, errors.New("strategy not configured")
}

func (c *ItemCache) Save(item entities.Item) error {
	if c.strategy != nil {
		return c.strategy.Save(item)
	}

	return errors.New("strategy not configured")
}

func (c *ItemCache) GetStrategy(storage string) Cacher {
	for _, strategy := range []Cacher{c.redisCache, c.postgresCache} {
		if strategy.MustApply(storage) {
			return strategy
		}
	}

	return c.postgresCache
}

func (c *ItemCache) SetStrategy(strategy Cacher) {
	c.strategy = strategy
}
