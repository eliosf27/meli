package main

import (
	"github.com/labstack/echo/v4"
	"meli/internal/app/entities"
	"time"
)

// Routes build the routes of the server
func (s *Server) Routes() {
	s.server.GET("/", s.dependencies.StatusController.Status)
	s.server.GET("/health", s.dependencies.MetricController.Health)

	items := s.server.Group("/items")
	items.GET("/:item_id", s.Tracking(s.dependencies.ItemController.Get))
}

// Tracking track and save the local requests
func (s *Server) Tracking(callback func(ctx echo.Context) error) func(ctx echo.Context) error {
	return func(ctx echo.Context) error {
		start := time.Now()
		err := callback(ctx)
		elapsed := time.Since(start)

		s.dependencies.Queue.Enqueue(entities.NewLocalMetric(elapsed.Milliseconds()))

		return err
	}
}
