package item

import (
	log "github.com/sirupsen/logrus"
	config "meli/internal/app/config"
	"meli/internal/entities"
	"meli/internal/postgres"
)

const (
	InsertItemQuery         = "insert_item"
	InsertItemChildrenQuery = "insert_item_children"
	SelectItemsChildren     = "select_items_children"
)

type ItemPostgresCache struct {
	postgres postgres.Postgres
}

func NewItemPostgresCache(postgres postgres.Postgres) Cacher {
	return ItemPostgresCache{postgres: postgres}
}

func (r ItemPostgresCache) MustApply(strategy string) bool {
	return strategy == PostgresCache
}

func (r ItemPostgresCache) Get(id string) (entities.Item, error) {
	item, err := r.getItem(id)
	if err != nil {
		log.Errorf("ItemPostgresCache.Get | Error getting children from DB: %+v", err)

		return entities.Item{}, err
	}

	return item, nil
}

func (r ItemPostgresCache) Save(item entities.Item) error {
	err := r.postgres.Execute(
		InsertItemQuery,
		item.ItemId, item.Title, item.CategoryId,
		item.Price, item.StartTime, item.StopTime,
	)
	if err != nil {
		log.Errorf(
			"ItemPostgresCache.Execute Item | item : %v"+
				"Error executing insert: %+v", item, err,
		)

		return err
	}

	for _, child := range item.Children {
		err := r.postgres.Execute(
			InsertItemChildrenQuery,
			child.ItemId, child.StopTime,
		)
		if err != nil {
			log.Errorf(
				"ItemPostgresCache.Execute Item child | child : %v"+
					"Error executing insert: %+v", child, err,
			)

			return err
		}
	}

	return nil
}

func (r ItemPostgresCache) Name() string {

	return config.PostgresStorage
}

func (r ItemPostgresCache) getItem(id string) (entities.Item, error) {
	rows, err := r.postgres.Query(SelectItemsChildren, id)
	if err != nil {

		return entities.Item{}, err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			log.Errorf("Error closing rows: %+v", err)
		}
	}()

	itemResult := entities.Item{}
	children := make([]entities.ItemChildren, 0)
	for rows.Next() {
		item := entities.Item{}
		itemChildren := entities.ItemChildren{}
		err := rows.Scan(
			&item.ItemId, &item.Title, &item.CategoryId,
			&item.Price, &item.StartTime, &item.StopTime,
			&itemChildren.ItemId, &itemChildren.StopTime,
		)

		if err != nil {

			return entities.Item{}, err
		} else {
			if itemResult.IsZero() {
				itemResult = item
			}

			if !itemChildren.IsZero() {
				children = append(children, itemChildren)
			}
		}
	}

	itemResult.Children = children

	return itemResult, nil
}
