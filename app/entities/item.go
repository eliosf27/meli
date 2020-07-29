package entities

type Item struct {
	ItemId     string         `json:"item_id"`
	Title      string         `json:"title"`
	CategoryId string         `json:"category_id"`
	Price      float64        `json:"price"`
	StartTime  string         `json:"start_time"`
	StopTime   string         `json:"stop_time"`
	Children   []ItemChildren `json:"children"`
}

// IsZero checks if the item is empty
func (i *Item) IsZero() bool {
	return i.ItemId == ""
}

type ItemChildren struct {
	ItemId   *string `json:"parent_item_id"`
	StopTime *string `json:"stop_time"`
}

// IsZero checks if the children item is empty
func (i *ItemChildren) IsZero() bool {
	return i.ItemId == nil || *i.ItemId == ""
}
