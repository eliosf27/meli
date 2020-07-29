package http

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"meli/internal/container"
)

type Server struct {
	server *echo.Echo
}

func NewServer() *Server {
	return &Server{
		server: echo.New(),
	}
}

// Middleware run the server
func (s *Server) Start(dependencies container.Dependencies) {
	s.server.Logger.Fatal(s.server.Start(fmt.Sprintf(":%s", dependencies.Config.Port)))
}
