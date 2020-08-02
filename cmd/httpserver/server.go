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

// Start run the server
func (s *Server) Start() {
	s.StartQueueConsumers()
	s.server.Logger.Fatal(s.server.Start(fmt.Sprintf(":%s", s.dependencies.Config.Port)))
}

// StartQueueConsumers run the queue consumers
func (s *Server) StartQueueConsumers() {
	for _, consumer := range s.dependencies.QueueConsumers {
		go consumer.Listen()
	}
}
