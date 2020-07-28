package item

import (
	log "github.com/sirupsen/logrus"
	"meli/app/entities"
	"meli/internal/api"
)

type ItemService interface {
	FetchItemById(id string) entities.Item
}

type service struct {
	httpClient     api.HttpClient
	itemRepository ItemRepository
}

func NewItemService(httpClient api.HttpClient, items ItemRepository) ItemService {
	return &service{itemRepository: items, httpClient: httpClient}
}

func (s *service) FetchItemById(id string) entities.Item {
	item, err := s.itemRepository.Get(id)
	if err != nil {
		log.Errorf("error fetching item | error: %+v", err)

		return entities.Item{}
	}

	if item.IsZero() {
		item = s.httpClient.ItemService.Get(id)
		err = s.itemRepository.Save(item)
		if err != nil {
			log.Errorf("error saving item | error: %+v", err)

			return item
		}
	}

	return item
}
