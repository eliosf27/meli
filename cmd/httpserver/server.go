package httpserver

import (
	"github.com/labstack/echo/v4"
)

type Server struct {
	e *echo.Echo
}

func NewServer() *Server {
	e := echo.New()
	s := &Server{
		e: e,
	}
	return s
}

func (s *Server) Start() {
	s.e.Logger.Fatal(s.e.Start(":8000"))
}
