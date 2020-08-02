package main

import "github.com/labstack/echo/v4"

// Middleware build the routes of the server
func (s *Server) Routes() {
	s.server.GET("/", s.dependencies.StatusController.Status)

	items := s.server.Group("/items")
	items.GET("/:item_id", s.dependencies.ItemController.Get)
}

func C() func(ctx echo.Context) error {

	return nil
}
