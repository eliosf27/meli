package main

import (
	"meli/cmd/http"
	"meli/internal/container"
)

func main() {
	controllerGroup := container.Build()

	server := http.NewServer()
	server.Middleware()
	server.Routes(controllerGroup)
	server.Start()
}
