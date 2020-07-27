package http

import (
	"github.com/labstack/echo/v4"
)

type Server struct {
	server *echo.Echo
}

func NewServer() *Server {
	return &Server{
		server: echo.New(),
	}
}

func (s *Server) Start() {
	s.server.Logger.Fatal(s.server.Start(":8000"))
}
