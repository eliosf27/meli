package http

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"meli/internal/entities"
)

const (
	itemsBasePath = "items/"
	childrenPath  = "%s/children"
)

type (
	ItemHttpService struct {
		httpClient HttpClient
	}
)

func NewItemHttpService(sling *HttpClient) ItemHttpService {
	return ItemHttpService{
		httpClient: sling.Path(itemsBasePath),
	}
}

// getItem fetch the item from the http in a sync way
func (s ItemHttpService) getItem(id string) (entities.Item, error) {
	resp := entities.Item{
		ItemId: id,
	}
	res, err := s.httpClient.Get(id, &resp)
	if err != nil || res == nil {
		log.Error("Error retrieving item: ", err)

		return resp, err
	}

	return resp, nil
}

// getItemChildren fetch the children item from the http in a sync way
func (s ItemHttpService) getItemChildren(id string) ([]entities.ItemChildren, error) {
	var resp []entities.ItemChildren
	path := fmt.Sprintf(childrenPath, id)
	res, err := s.httpClient.Get(path, &resp)
	if err != nil || res == nil {
		log.Error("Error retrieving children item: ", err)

		return resp, err
	}

	return resp, nil
}

// getItemChildren fetch the children item from the http in a async way
func (s ItemHttpService) getItemChildrenAsync(id string) <-chan []entities.ItemChildren {
	future := make(chan []entities.ItemChildren)

	go func() {
		if children, err := s.getItemChildren(id); err != nil {
			future <- []entities.ItemChildren{}
		} else {
			future <- children
		}
	}()

	return future
}

// getItem fetch the item from the http in a async way
func (s ItemHttpService) getItemAsync(id string) <-chan entities.Item {
	future := make(chan entities.Item)

	go func() {
		if item, err := s.getItem(id); err != nil {
			future <- entities.Item{}
		} else {
			future <- item
		}
	}()

	return future
}

// Get return the item, fetching data from the http and building the item
func (s ItemHttpService) Get(id string) entities.Item {
	itemAsync := s.getItemAsync(id)
	childAsync := s.getItemChildrenAsync(id)

	item := <-itemAsync
	child := <-childAsync

	item.Children = child

	return item
}

// GetItemPath return the item path
func (s ItemHttpService) GetItemPath(id string) string {
	return itemsBasePath + id
}

// GetItemChildrenPath return the item children path
func (s ItemHttpService) GetItemChildrenPath(id string) string {
	return itemsBasePath + fmt.Sprintf(childrenPath, id)
}
