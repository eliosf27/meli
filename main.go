package main

import (
	"meli/cmd/http"
	"meli/internal/container"
)

func main() {
	dependencies := container.Build()

	server := http.NewServer()
	server.Middleware()
	server.Routes(dependencies)
	server.Start(dependencies)
}
