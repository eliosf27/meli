package item

import (
	"database/sql"
	log "github.com/sirupsen/logrus"
	config "meli/internal/app/config"
	"meli/internal/entities"
	"meli/internal/postgres"
	"meli/pkg/queries"
)

const (
	InsertItem          = "insert_item"
	InsertItemChildren  = "insert_item_children"
	SelectItemsChildren = "select_items_children"
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
	client := r.postgres.Client

	tx, err := client.Begin()
	if err != nil {
		log.Errorf("ItemPostgresCache.SaveChildren | Error connecting to DB: %+v", err)

		return err
	}

	err = r.save(tx, item)
	if err != nil {
		log.Errorf(
			"ItemPostgresCache.save | item : %v"+
				"Error executing insert: %+v", item, err,
		)

		return err
	}

	for _, child := range item.Children {
		err := r.saveChildren(tx, child)
		if err != nil {
			log.Errorf(
				"ItemPostgresCache.saveChildren | child : %v"+
					"Error executing insert: %+v", child, err,
			)

			return err
		}
	}

	if err := tx.Commit(); err != nil {
		log.Errorf("ItemPostgresCache.SaveChildren | Error committing transaction: %s", err.Error())

		return err
	}

	return nil
}

func (r ItemPostgresCache) Name() string {

	return config.PostgresStorage
}

func (r ItemPostgresCache) save(tx *sql.Tx, item entities.Item) error {
	query, err := queries.ReadQuery(InsertItem)
	if err != nil {
		return err
	}

	_, err = tx.Exec(query,
		item.ItemId, item.Title, item.CategoryId,
		item.Price, item.StartTime, item.StopTime,
	)
	if err != nil {
		tx.Rollback()

		return err
	}

	return nil
}

func (r ItemPostgresCache) saveChildren(tx *sql.Tx, itemChildren entities.ItemChildren) error {
	query, err := queries.ReadQuery(InsertItemChildren)
	if err != nil {

		return err
	}
	_, err = tx.Exec(query,
		itemChildren.ItemId, itemChildren.StopTime,
	)
	if err != nil {
		tx.Rollback()

		return err
	}

	return nil
}

func (r ItemPostgresCache) getItem(id string) (entities.Item, error) {
	client := r.postgres.Client
	query, err := queries.ReadQuery(SelectItemsChildren)
	if err != nil {
		log.Errorf("ItemPostgresCache.SaveChildren | Error reading query: %+v", err)

		return entities.Item{}, err
	}

	rows, err := client.Query(query, id)
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
