package main

import (
	"meli/internal/container"
)

func main() {
	dependencies := container.Build()

	server := NewServer(dependencies)
	server.Middleware()
	server.Routes()
	server.Start()
}
