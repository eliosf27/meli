package api

import (
	"fmt"
	"github.com/dghubble/sling"
	log "github.com/sirupsen/logrus"
	"meli/app"
	"net/http"
)

const (
	itemsBasePath = "/items/"
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

func (s ItemService) getItem(id string) (app.Item, error) {
	resp := app.Item{
		ItemId: id,
	}
	res, err := s.sling.Get(id).ReceiveSuccess(&resp)
	if err != nil || res == nil || res.StatusCode != http.StatusOK {
		log.Error("Error retrieving toolkits: ", err)

		return resp, err
	}

	return resp, nil
}

func (s ItemService) getItemChildren(id string) ([]app.ItemChildren, error) {
	var resp []app.ItemChildren
	path := fmt.Sprintf(childrenPath, id)
	res, err := s.sling.New().Get(path).ReceiveSuccess(&resp)
	if err != nil || res == nil || res.StatusCode != http.StatusOK {
		log.Error("Error retrieving toolkits: ", err)

		return resp, err
	}

	return resp, nil
}

func (s ItemService) getItemChildrenAsync(id string) <-chan []app.ItemChildren {
	future := make(chan []app.ItemChildren)

	go func() {
		if children, err := s.getItemChildren(id); err != nil {
			future <- []app.ItemChildren{}
		} else {
			future <- children
		}
	}()

	return future
}

func (s ItemService) getItemAsync(id string) <-chan app.Item {
	future := make(chan app.Item)

	go func() {
		if item, err := s.getItem(id); err != nil {
			future <- app.Item{}
		} else {
			future <- item
		}
	}()

	return future
}

func (s ItemService) Get(id string) app.Item {
	itemAsync := s.getItemAsync(id)
	childAsync := s.getItemChildrenAsync(id)

	item := <-itemAsync
	child := <-childAsync

	item.Children = child

	return item
}
