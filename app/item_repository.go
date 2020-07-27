package app

import config "meli/pkg/config"

type ItemRepository interface {
	Get(Id string) (Item, error)
	Save(item Item) error
}

type ItemRepository2 struct{}

func NewItemRepository(Config config.Config) ItemRepository {
	return ItemRepository2{}
}

func (ItemRepository2) Get(Id string) (Item, error) {
	panic("implement me")
}

func (ItemRepository2) Save(item Item) error {
	panic("implement me")
}
