package entities

// ItemMetric the type of the storage
type ItemMetric struct {
	ResponsesTime []float64
	StatusCode    map[int]int64
}
