package caching

import (
	log "github.com/sirupsen/logrus"
	"meli/api"
	"meli/app"
)

type ItemService interface {
	FetchItemByID(id string) app.Item
}

type service struct {
	httpClient     api.HttpClient
	itemRepository app.ItemRepository
}

func NewItemService(httpClient api.HttpClient, items app.ItemRepository) ItemService {
	return &service{itemRepository: items, httpClient: httpClient}
}

func (s *service) FetchItemByID(id string) app.Item {
	item, err := s.itemRepository.Get(id)
	if err != nil {
		log.Errorf("error fetching item | error: %+v", err)

		return app.Item{}
	}

	if item.IsZero() {
		item = s.httpClient.ItemService.Build(id)
		err = s.itemRepository.Save(item)
		if err != nil {
			log.Errorf("error saving item | error: %+v", err)

			return app.Item{}
		}
	}

	return item
}
