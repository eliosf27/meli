package main

import (
	"github.com/labstack/echo/v4"
	"meli/internal/entities"
	"time"
)

// Routes build the routes of the server
func (s *Server) Routes() {
	s.server.GET("/", s.dependencies.StatusHandler.Status)
	s.server.GET("/health", s.dependencies.MetricHandler.Health)

	items := s.server.Group("/items")
	items.GET("/:item_id", s.Tracking(s.dependencies.ItemHandler.Get))

	cache := s.server.Group("/config/cache")
	cache.GET("/", s.dependencies.ConfigCacheHandler.Fetch)
	cache.PUT("/", s.dependencies.ConfigCacheHandler.Update)
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
