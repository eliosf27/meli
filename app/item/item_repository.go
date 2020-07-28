package item

import (
	"database/sql"
	log "github.com/sirupsen/logrus"
	"meli/app/entities"
	"meli/internal/postgres"
	"meli/pkg/queries"
)

const (
	InsertItem          = "insert_item"
	InsertItemChildren  = "insert_item_children"
	SelectItemsChildren = "select_items_children"
)

type ItemRepository interface {
	Get(Id string) (entities.Item, error)
	Save(item entities.Item) error
}

type ItemRepo struct {
	Postgres postgres.Postgres
}

func NewItemRepository(postgres postgres.Postgres) ItemRepository {
	return ItemRepo{Postgres: postgres}
}

func (r ItemRepo) Get(id string) (entities.Item, error) {
	item, err := r.getItem(id)
	if err != nil {
		log.Errorf("ItemRepository.Get | Error getting children from DB: %+v", err)

		return entities.Item{}, err
	}

	return item, nil
}

func (r ItemRepo) Save(item entities.Item) error {
	client := r.Postgres.Client

	tx, err := client.Begin()
	if err != nil {
		log.Errorf("ItemRepository.SaveChildren | Error connecting to DB: %+v", err)

		return err
	}

	err = r.save(tx, item)
	if err != nil {
		log.Errorf(
			"ItemRepository.save | item : %v"+
				"Error executing insert: %+v", item, err,
		)
	}

	for _, child := range item.Children {
		err := r.saveChildren(tx, child)
		if err != nil {
			log.Errorf(
				"ItemRepository.saveChildren | child : %v"+
					"Error executing insert: %+v", child, err,
			)
		}
	}

	if err := tx.Commit(); err != nil {
		log.Errorf("ItemRepository.SaveChildren | Error committing transaction: %s", err.Error())

		return err
	}

	return nil
}

func (r ItemRepo) save(tx *sql.Tx, item entities.Item) error {
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

func (r ItemRepo) saveChildren(tx *sql.Tx, itemChildren entities.ItemChildren) error {
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

func (r ItemRepo) getItem(id string) (entities.Item, error) {
	client := r.Postgres.Client
	query, err := queries.ReadQuery(SelectItemsChildren)
	if err != nil {
		log.Errorf("ItemRepository.SaveChildren | Error reading query: %+v", err)

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
			children = append(children, itemChildren)
		}
	}

	itemResult.Children = children

	return itemResult, nil
}
