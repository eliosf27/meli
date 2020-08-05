package queue

import (
	"meli/internal/app/metric"
	"meli/internal/entities"
)

type Consumer interface {
	Listen()
}

type ItemConsumer struct {
	queue         *ItemQueue
	metricService metric.MetricServicer
}

func NewItemConsumer(queue *ItemQueue, metricService metric.MetricServicer) ItemConsumer {
	return ItemConsumer{
		queue:         queue,
		metricService: metricService,
	}
}

func (c ItemConsumer) Listen() {
	c.queue.Listen(func(item entities.ItemMetric) error {

		return c.metricService.UpdateMetric(item)
	})
}
