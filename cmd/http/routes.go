package http

import (
	"meli/internal/container"
)

// Middleware build the routes of the server
func (s *Server) Routes(group container.Dependencies) {
	s.server.GET("/", group.StatusController.Status)

	items := s.server.Group("/items")
	items.GET("/:item_id", group.ItemController.Get)
}
