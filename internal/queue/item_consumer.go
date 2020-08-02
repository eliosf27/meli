package queue

import (
	log "github.com/sirupsen/logrus"
	"meli/internal/app/entities"
)

type Consumer interface {
	Listen()
}

type ItemConsumer struct {
	queue *ItemQueue
}

func NewItemConsumer(queue *ItemQueue) ItemConsumer {
	return ItemConsumer{
		queue: queue,
	}
}

func (c ItemConsumer) Listen() {
	acc := entities.ItemMetric{
		ResponsesTime: []int64{},
		StatusCode:    map[int]int64{},
	}
	c.queue.Listen(func(item Item) error {
		acc.ResponsesTime = append(acc.ResponsesTime, item.ResponseTime)
		acc.StatusCode[item.StatusCode] += 1
		log.Info("item: ", item)
		log.Info("acc.ResponsesTime: ", acc.ResponsesTime)
		log.Info("acc.StatusCod: ", acc.StatusCode)
		return nil
	})
}
