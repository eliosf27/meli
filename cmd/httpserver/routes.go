package main

import (
	"github.com/labstack/echo/v4"
	"meli/internal/queue"
	"time"
)

// Routes build the routes of the server
func (s *Server) Routes() {
	s.server.GET("/", s.dependencies.StatusController.Status)

	items := s.server.Group("/items")
	items.GET("/:item_id", s.Tracking(s.dependencies.ItemController.Get))
}

// Tracking track and save the local requests
func (s *Server) Tracking(callback func(ctx echo.Context) error) func(ctx echo.Context) error {
	return func(ctx echo.Context) error {
		start := time.Now()
		err := callback(ctx)
		elapsed := time.Since(start)

		s.dependencies.Queue.Enqueue(queue.Item{
			Type:         queue.LocalApi,
			ResponseTime: elapsed.Milliseconds(),
		})

		return err
	}
}
