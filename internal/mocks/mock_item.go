package mock_item

import "meli/app/entities"

func MockItem(id string) entities.Item {
	return entities.Item{
		ItemId:     id,
		Title:      "Google Pixel 32gb",
		CategoryId: "MLU1055",
		Price:      27,
		StartTime:  "2019-03-02T20:31:02.000Z",
		StopTime:   "2019-10-25T23:28:35.000Z",
		Children:   nil,
	}
}

func MockItemChildren(id string) []entities.ItemChildren {
	return []entities.ItemChildren{
		{ItemId: id, StopTime: "2019-03-02T20:31:02.000Z"},
		{ItemId: id, StopTime: "2019-03-02T20:31:02.000Z"},
	}
}
