package item

import (
	log "github.com/sirupsen/logrus"
	config "meli/internal/app/config"
	"meli/internal/entities"
	"meli/internal/http"
)

type ItemServicer interface {
	FetchItemById(id string) entities.Item
}

type ItemService struct {
	itemHttpService http.ItemHttpService
	itemCache       ItemCacher
	configService   config.ConfigServicer
}

func NewItemService(httpClient http.ItemHttpService, itemCache ItemCacher, configService config.ConfigServicer) ItemServicer {
	return &ItemService{itemCache: itemCache, itemHttpService: httpClient, configService: configService}
}

func (s *ItemService) FetchItemById(id string) entities.Item {
	storage := s.Storage()
	s.itemCache.SetStrategy(storage)

	item, _ := s.itemCache.Get(id)
	if item.IsZero() {
		item = s.itemHttpService.Get(id)
		err := s.itemCache.Save(item)
		if err != nil {
			log.Errorf("error saving config | error: %+v", err)

			return item
		}
	}

	return item
}

func (s *ItemService) Storage() Cacher {
	storage, err := s.configService.Fetch()
	if err != nil {
		log.Errorf("ItemService.FetchItemById | using default storage | error: %+v", err)

		storage = config.PostgresStorage
	}

	return s.itemCache.GetStrategy(storage)
}
