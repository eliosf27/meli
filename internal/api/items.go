package api

import (
	"fmt"
	"github.com/dghubble/sling"
	log "github.com/sirupsen/logrus"
	"meli/internal/app/entities"
	"net/http"
)

const (
	itemsBasePath = "items/"
	childrenPath  = "%s/children"
)

type (
	ItemService struct {
		sling *sling.Sling
	}
)

func NewItemService(sling *sling.Sling) *ItemService {
	return &ItemService{
		sling: sling.Path(itemsBasePath),
	}
}

// getItem fetch the item from the api in a sync way
func (s ItemService) getItem(id string) (entities.Item, error) {
	resp := entities.Item{
		ItemId: id,
	}
	res, err := s.sling.Get(id).ReceiveSuccess(&resp)
	if err != nil || res == nil || res.StatusCode != http.StatusOK {
		log.Error("Error retrieving item: ", err)

		return resp, err
	}

	return resp, nil
}

// getItemChildren fetch the children item from the api in a sync way
func (s ItemService) getItemChildren(id string) ([]entities.ItemChildren, error) {
	var resp []entities.ItemChildren
	path := fmt.Sprintf(childrenPath, id)
	res, err := s.sling.New().Get(path).ReceiveSuccess(&resp)
	if err != nil || res == nil || res.StatusCode != http.StatusOK {
		log.Error("Error retrieving children item: ", err)

		return resp, err
	}

	return resp, nil
}

// getItemChildren fetch the children item from the api in a async way
func (s ItemService) getItemChildrenAsync(id string) <-chan []entities.ItemChildren {
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

// getItem fetch the item from the api in a async way
func (s ItemService) getItemAsync(id string) <-chan entities.Item {
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

// Get return the item, fetching data from the api and building the item
func (s ItemService) Get(id string) entities.Item {
	itemAsync := s.getItemAsync(id)
	childAsync := s.getItemChildrenAsync(id)

	item := <-itemAsync
	child := <-childAsync

	item.Children = child

	return item
}

// GetItemPath return the item path
func (s ItemService) GetItemPath(id string) string {
	return itemsBasePath + id
}

// GetItemChildrenPath return the item children path
func (s ItemService) GetItemChildrenPath(id string) string {
	return itemsBasePath + fmt.Sprintf(childrenPath, id)
}
