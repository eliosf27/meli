package item

import "meli/internal/entities"

const (
	RedisCache    = "redis"
	PostgresCache = "postgres"
)

type ItemCacher interface {
	Get(Id string) (entities.Item, error)
	Save(item entities.Item) error
}

type Cacher interface {
	Get(Id string) (entities.Item, error)
	Save(item entities.Item) error
	MustApply(strategy string) bool
}

type ItemCache struct {
	strategy   string
	strategies []Cacher
}

func NewItemCache(strategies []Cacher) *ItemCache {
	return &ItemCache{strategies: strategies}
}

func (c *ItemCache) Get(id string) (entities.Item, error) {
	strategy := c.getStrategy()

	return strategy.Get(id)
}

func (c *ItemCache) Save(item entities.Item) error {
	strategy := c.getStrategy()

	return strategy.Save(item)
}

func (c *ItemCache) getStrategies() []Cacher {

	return c.strategies
}

func (c *ItemCache) getStrategy() Cacher {
	for _, strategy := range c.getStrategies() {
		if strategy.MustApply(c.strategy) {

			return strategy
		}
	}
	return nil
}

func (c *ItemCache) SetStrategy(strategy string) {
	c.strategy = strategy
}
