package main

import (
	"meli/internal/container"
)

func main() {
	dependencies := container.Build()

	server := NewServer(dependencies)
	server.Middleware()
	server.Validator()
	server.Routes()
	server.Start()
}
