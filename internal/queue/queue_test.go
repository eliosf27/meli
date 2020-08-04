package queue

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"meli/internal/entities"
	"testing"
	"time"
)

func TestItem2222(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Item Suite")
}

var _ = Describe("Queue", func() {
	Context("using a queue", func() {
		When("enqueue a metric", func() {
			It("should dequeue a valid metric", func() {
				itemQueue := NewItemQueue()

				itemMetric := entities.ItemMetric{
					Type:         "xxx",
					ResponseTime: 210,
					StatusCode:   200,
					Time:         time.Now(),
				}
				itemQueue.Enqueue(itemMetric)

				itemMetricExpected := itemQueue.Dequeue()
				Expect(itemMetricExpected.Type).To(Equal(itemMetric.Type))
			})
		})

		When("enqueue a metric", func() {
			It("should not be empty", func() {
				itemQueue := NewItemQueue()

				itemMetric := entities.ItemMetric{
					Type:         "xxx",
					ResponseTime: 210,
					StatusCode:   200,
					Time:         time.Now(),
				}
				itemQueue.Enqueue(itemMetric)

				Expect(itemQueue.IsEmpty()).To(Equal(false))
			})
		})

		When("there is not metrics in the queue", func() {
			It("should be empty", func() {
				itemQueue := NewItemQueue()

				Expect(itemQueue.IsEmpty()).To(Equal(true))
			})
		})

		When("a consumer is running", func() {
			It("should read a new queue message", func() {
				itemQueue := NewItemQueue()
				itemMetric := entities.ItemMetric{
					Type:         "xxx",
					ResponseTime: 210,
					StatusCode:   200,
					Time:         time.Now(),
				}

				go itemQueue.Listen(func(item entities.ItemMetric) error {
					Expect(item.Type).To(Equal(itemMetric.Type))

					return nil
				})

				itemQueue.Enqueue(itemMetric)

				itemQueue.Stop()
			})
		})
	})
})
