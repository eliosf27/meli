package mock_item

import "meli/internal/entities"

// MockItem create a fake entity Item
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

// MockItemChildren create a fake entity ItemChildren
func MockItemChildren(id string) []entities.ItemChildren {
	stopTime := "2019-03-02T20:31:02.000Z"
	return []entities.ItemChildren{
		{ItemId: &id, StopTime: &stopTime},
		{ItemId: &id, StopTime: &stopTime},
	}
}
