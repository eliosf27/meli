package queue

import (
	log "github.com/sirupsen/logrus"
	"sync"
	"time"
)

const (
	LocalApi    = "local_api"
	ExternalApi = "external_api"
)

// Item the type of the queue
type Item struct {
	Type         string
	ResponseTime int64
	StatusCode   int
}

func (s *Item) IsZero() bool {
	return s.Type == ""
}

// ItemQueue the queue of Items
type ItemQueue struct {
	items []Item
	lock  sync.RWMutex
}

// NewItemQueue creates a new ItemQueue
func NewItemQueue() ItemQueue {
	return ItemQueue{
		items: []Item{},
		lock:  sync.RWMutex{},
	}
}

// Enqueue adds an Item to the end of the queue
func (s *ItemQueue) Enqueue(t Item) {
	s.lock.Lock()
	s.items = append(s.items, t)
	s.lock.Unlock()
}

// Dequeue removes an Item from the start of the queue
func (s *ItemQueue) Dequeue() *Item {
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
func (s *ItemQueue) Listen(callback func(item Item) error) {
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
