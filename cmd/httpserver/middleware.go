package main

import (
	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
	"strings"
)

var (
	excludedGzipPaths = []string{"docs", "metrics"}
)

// Middleware build the middlewares of the server
func (s *Server) Middleware() {
	s.server.Use(echoMiddleware.LoggerWithConfig(echoMiddleware.LoggerConfig{
		CustomTimeFormat: "2006-01-02T15:04:05.1483386-00:00",
		Format:           "[${time_custom}][INFO] [method=${method}] [uri=${uri}] [status=${status}] [origin=${header:X-Application-ID}]\n",
	}))
	s.server.Use(echoMiddleware.Recover())
	s.server.Use(echoMiddleware.GzipWithConfig(echoMiddleware.GzipConfig{
		Skipper: func(c echo.Context) bool {
			for _, path := range excludedGzipPaths {
				if strings.Contains(c.Request().URL.Path, path) {
					return true
				}
			}
			return false
		},
	}))
}
