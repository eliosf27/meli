package main

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"meli/internal/container"
)

type Server struct {
	server       *echo.Echo
	dependencies container.Dependencies
}

func NewServer(dependencies container.Dependencies) *Server {
	return &Server{
		server:       echo.New(),
		dependencies: dependencies,
	}
}

// Middleware run the server
func (s *Server) Start() {
	s.server.Logger.Fatal(s.server.Start(fmt.Sprintf(":%s", s.dependencies.Config.Port)))
}
