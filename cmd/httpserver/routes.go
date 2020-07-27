package httpserver

import (
	"meli/internal/container"
)

func (s *Server) Routes(group container.ControllerGroup) {
	s.e.GET("/", group.StatusController.Status)

	items := s.e.Group("/items")
	items.GET("/:item_id", group.ItemController.Get)
}
