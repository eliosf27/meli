package main

import (
	log "github.com/sirupsen/logrus"
	"meli/cmd/http"
	"meli/internal/container"
	"meli/pkg/testcontainers"
)

func main() {
	controllerGroup := container.Build()

	c := testcontainers.PostgresContainer{}
	connection := c.Up()
	log.Info(connection.Host, connection.Port, connection.User, connection.Password)

	server := http.NewServer()
	server.Middleware()
	server.Routes(controllerGroup)
	server.Start()
}
