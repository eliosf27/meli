package item

import (
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"meli/internal/entities"
	"meli/internal/redis"
	"time"
)

type ItemRedisCache struct {
	redis redis.Redis
}

func NewItemRedisCache(redis redis.Redis) Cacher {
	return ItemRedisCache{redis: redis}
}

func (r ItemRedisCache) MustApply(strategy string) bool {
	return strategy == RedisCache
}

func (r ItemRedisCache) Get(id string) (entities.Item, error) {
	key := r.buildKey(id)
	val, err := r.redis.Get(key)
	if err != nil {
		log.Errorf("ItemRedisCache.Get | error fetching item cache [%s] - %+v", key, err)

		return entities.Item{}, err
	}

	item := entities.Item{}
	if err := json.Unmarshal(val, &item); err != nil {
		log.Errorf("ItemRedisCache.Get | error unmarshal: %+v", err)

		return entities.Item{}, err
	}

	return item, nil
}

func (r ItemRedisCache) Save(item entities.Item) error {
	key := r.buildKey(item.ItemId)
	itemRaw, err := json.Marshal(&item)
	if err != nil {
		log.Errorf(
			"ItemRedisCache.Save | error marshalling item for saving cache [%s] - %+v", key, err,
		)

		return err
	}

	_, err = r.redis.Set(key, itemRaw, time.Minute*10)
	if err != nil {
		log.Errorf("ItemRedisCache.Save |error saving item cache [%s] - %+v", key, err)

		return err
	}

	return nil
}

func (r ItemRedisCache) buildKey(id string) string {

	return fmt.Sprintf("meli:cache_item_%s", id)
}
