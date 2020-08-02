package item

import (
	log "github.com/sirupsen/logrus"
	"meli/internal/entities"
	"meli/internal/http"
)

type ItemService interface {
	FetchItemById(id string) entities.Item
}

type service struct {
	itemHttpService http.ItemHttpService
	itemCache       ItemCacher
}

func NewItemService(httpClient http.ItemHttpService, itemCache ItemCacher) ItemService {
	return &service{itemCache: itemCache, itemHttpService: httpClient}
}

func (s *service) FetchItemById(id string) entities.Item {
	item, err := s.itemCache.Get(id)
	if err != nil {
		log.Errorf("error fetching item | error: %+v", err)

		return entities.Item{}
	}

	if item.IsZero() {
		item = s.itemHttpService.Get(id)
		err = s.itemCache.Save(item)
		if err != nil {
			log.Errorf("error saving item | error: %+v", err)

			return item
		}
	}

	return item
}
