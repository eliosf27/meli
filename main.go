package main

import (
	"meli/cmd/httpserver"
	"meli/internal/container"
)

func main() {
	group := container.Build()

	server := httpserver.NewServer()
	server.Middleware()
	server.Routes(group)
	server.Start()
}
