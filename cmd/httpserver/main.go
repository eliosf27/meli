package main

import (
	"meli/internal/container"
	"meli/internal/queue"
)

func main() {
	dependencies := container.Build()

	server := NewServer(dependencies)
	server.Middleware()
	server.Routes()
	// server.Start()

	s := queue.NewItemQueue()
	s.Enqueue(queue.Item{
		Type:         queue.LocalApi,
		ResponseTime: 220,
		StatusCode:   200,
	})
	s.Enqueue(queue.Item{
		Type:         queue.LocalApi,
		ResponseTime: 789,
		StatusCode:   200,
	})
	s.Enqueue(queue.Item{
		Type:         queue.LocalApi,
		ResponseTime: 456,
		StatusCode:   200,
	})

	//acc := entities.ItemMetric{
	//	ResponsesTime: []float64{},
	//	StatusCode:    map[int]int64{},
	//}
	//s.Listen(func(item queue.Item) error {
	//	acc.ResponsesTime = append(acc.ResponsesTime, item.ResponseTime)
	//	acc.StatusCode[item.StatusCode] += 1
	//	log.Info("item: ", item)
	//	log.Info("acc.ResponsesTime: ", acc.ResponsesTime)
	//	log.Info("acc.StatusCod: ", acc.StatusCode)
	//	return nil
	//})

	server.Start()
}
