package item

import (
	log "github.com/sirupsen/logrus"
	"meli/internal/app/entities"
	"meli/internal/http"
)

type ItemService interface {
	FetchItemById(id string) entities.Item
}

type service struct {
	itemHttpService http.ItemHttpService
	itemRepository  ItemRepository
}

func NewItemService(httpClient http.ItemHttpService, items ItemRepository) ItemService {
	return &service{itemRepository: items, itemHttpService: httpClient}
}

func (s *service) FetchItemById(id string) entities.Item {
	item, err := s.itemRepository.Get(id)
	if err != nil {
		log.Errorf("error fetching item | error: %+v", err)

		return entities.Item{}
	}

	if item.IsZero() {
		item = s.itemHttpService.Get(id)
		err = s.itemRepository.Save(item)
		if err != nil {
			log.Errorf("error saving item | error: %+v", err)

			return item
		}
	}

	return item
}
