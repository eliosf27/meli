package app

import (
	"database/sql"
	log "github.com/sirupsen/logrus"
	"meli/internal/postgres"
	"meli/kit/queries"
)

const (
	InsertItem          = "insert_item"
	InsertItemChildren  = "insert_item_children"
	SelectItems         = "select_items"
	SelectItemsChildren = "select_items_children"
)

type ItemRepository interface {
	Get(Id string) (Item, error)
	Save(item Item) error
}

type ItemRepo struct {
	Postgres postgres.Postgres
}

func NewItemRepository(postgres postgres.Postgres) ItemRepository {
	return ItemRepo{Postgres: postgres}
}

func (r ItemRepo) Get(id string) (Item, error) {
	item, err := r.get(id)
	if err != nil {
		log.Errorf("ItemRepository.Get | Error getting item from DB: %+v", err)

		return Item{}, err
	}
	children, err := r.getChildren(id)
	if err != nil {
		log.Errorf("ItemRepository.Get | Error getting children from DB: %+v", err)

		return item, err
	}

	item.Children = children

	return item, nil
}

func (r ItemRepo) Save(item Item) error {
	client := r.Postgres.Client

	tx, err := client.Begin()
	if err != nil {
		log.Errorf("ItemRepository.SaveChildren | Error connecting to DB: %+v", err)

		return err
	}

	err = r.save(client, item)
	if err != nil {
		log.Errorf(
			"ItemRepository.SaveChildren | item : %v"+
				"Error executing insert: %+v", item, err,
		)

		_ = tx.Rollback()
	}

	for _, child := range item.Children {
		err := r.saveChildren(client, child)
		if err != nil {
			log.Errorf(
				"ItemRepository.SaveChildren | item : %v"+
					"Error executing insert: %+v", child, err,
			)

			_ = tx.Rollback()
		}
	}

	if err := tx.Commit(); err != nil {
		log.Errorf("ItemRepository.SaveChildren | Error committing transaction: %s", err.Error())
		_ = tx.Rollback()

		return err
	}

	return nil
}

func (r ItemRepo) save(client *sql.DB, item Item) error {
	query, err := queries.ReadQuery(InsertItem)
	if err != nil {
		log.Errorf("ItemRepository.Save | Error reading query: %+v", err)

		return err
	}

	_, err = client.Exec(query,
		item.ItemId, item.Title, item.CategoryId,
		item.Price, item.StartTime, item.StopTime,
	)
	if err != nil {
		log.Errorf(
			"ItemRepository.Save | item : %v"+
				"Error executing insert: %+v", item, err,
		)

		return err
	}

	return nil
}

func (r ItemRepo) saveChildren(client *sql.DB, itemChildren ItemChildren) error {
	query, err := queries.ReadQuery(InsertItemChildren)
	if err != nil {
		log.Errorf("ItemRepository.SaveChildren | Error reading query: %+v", err)

		return err
	}
	_, err = client.Exec(query,
		itemChildren.ItemId, itemChildren.StopTime,
	)
	if err != nil {
		return err
	}

	return nil
}

func (r ItemRepo) get(id string) (Item, error) {
	client := r.Postgres.Client

	query, err := queries.ReadQuery(SelectItems)
	if err != nil {
		log.Errorf("ItemRepository.SaveChildren | Error reading query: %+v", err)

		return Item{}, err
	}

	rows, err := client.Query(query, id)
	if err != nil {
		log.Errorf("Error executing query: %+v", err)

		return Item{}, err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			log.Errorf("Error closing rows: %+v", err)
		}
	}()

	item := Item{}
	for rows.Next() {
		if err := rows.Scan(&item.ItemId, &item.Title, &item.CategoryId, &item.Price, &item.StartTime, &item.StopTime); err != nil {
			log.Errorf("Error parsing response: %+v", err)

			return Item{}, err
		}
	}

	return item, nil
}

func (r ItemRepo) getChildren(id string) ([]ItemChildren, error) {
	client := r.Postgres.Client

	query, err := queries.ReadQuery(SelectItemsChildren)
	if err != nil {
		log.Errorf("ItemRepository.SaveChildren | Error reading query: %+v", err)

		return []ItemChildren{}, err
	}

	rows, err := client.Query(query, id)
	if err != nil {
		return []ItemChildren{}, err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			log.Errorf("Error closing rows: %+v", err)
		}
	}()

	var children []ItemChildren
	for rows.Next() {
		item := ItemChildren{}
		if err := rows.Scan(&item.ItemId, &item.StopTime); err != nil {
			return []ItemChildren{}, err
		} else {
			children = append(children, item)
		}
	}

	return children, nil
}
