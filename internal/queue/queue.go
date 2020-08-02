package queue

import (
	log "github.com/sirupsen/logrus"
	"meli/internal/entities"
	"sync"
	"time"
)

// ItemQueue the queue of Items
type ItemQueue struct {
	items []entities.ItemMetric
	lock  sync.RWMutex
}

// NewItemQueue creates a new ItemQueue
func NewItemQueue() ItemQueue {
	return ItemQueue{
		items: []entities.ItemMetric{},
		lock:  sync.RWMutex{},
	}
}

// Enqueue adds an ItemMetric to the end of the queue
func (s *ItemQueue) Enqueue(t entities.ItemMetric) {
	s.lock.Lock()
	s.items = append(s.items, t)
	s.lock.Unlock()
}

// Dequeue removes an ItemMetric from the start of the queue
func (s *ItemQueue) Dequeue() *entities.ItemMetric {
	s.lock.Lock()
	item := s.items[0]
	s.items = s.items[1:len(s.items)]
	s.lock.Unlock()
	return &item
}

// IsEmpty returns true if the queue is empty
func (s *ItemQueue) IsEmpty() bool {
	return len(s.items) == 0
}

// Listen check and process an item every n seconds
func (s *ItemQueue) Listen(callback func(item entities.ItemMetric) error) {
	for {
		time.Sleep(1 * time.Second)

		if s.IsEmpty() {
			continue
		}

		val := s.Dequeue()
		if val.IsZero() {
			continue
		}

		log.Infof(
			"reading item with type: %s, status_code: %d, response_time: %d",
			val.Type, val.StatusCode, val.ResponseTime,
		)

		err := callback(*val)
		if err != nil {
			log.Errorf(
				"error trying to executing the item with type: %s, status_code: %d, response_time: %d",
				val.Type, val.StatusCode, val.ResponseTime,
			)
			s.Enqueue(*val)
		}
	}
}
