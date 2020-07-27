package http

import (
	"meli/internal/container"
)

func (s *Server) Routes(group container.ControllerGroup) {
	s.server.GET("/", group.StatusController.Status)

	items := s.server.Group("/items")
	items.GET("/:item_id", group.ItemController.Get)
}
