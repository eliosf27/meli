package httpserver

import (
	"meli/internal/container"
)

func (s *Server) Routes(group container.ControllerGroup) {
	s.e.GET("/", group.StatusController.Status)
}
