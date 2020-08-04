package testcontainers

import (
	"context"
	"fmt"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

const (
	redisPort = "6379"
)

type RedisConnection struct {
	Host             string
	Port             string
	ConnectionString string
}

type RedisContainer struct {
	Container testcontainers.Container
}

func NewRedisContainer() RedisContainer {
	return RedisContainer{
		Container: nil,
	}
}

// Up create a new redis container with the connections params
func (c RedisContainer) Up() RedisConnection {
	containerRequest := testcontainers.ContainerRequest{
		Image:        "redis",
		ExposedPorts: []string{fmt.Sprintf("%s/tcp", redisPort)},
		WaitingFor:   wait.ForListeningPort(redisPort),
	}
	container, err := testcontainers.GenericContainer(context.Background(), testcontainers.GenericContainerRequest{
		ContainerRequest: containerRequest,
		Started:          true,
	})
	if err != nil {
		panic(err)
	}

	c.Container = container

	containerHost, _ := container.Host(context.Background())
	containerPort, _ := container.MappedPort(context.Background(), redisPort)

	return RedisConnection{
		Host: containerHost,
		Port: containerPort.Port(),
	}
}

// Down destroy the redis container
func (c RedisContainer) Down() error {
	if c.Container != nil {
		err := c.Container.Terminate(context.Background())
		if err != nil {
			return err
		}
	}

	return nil
}
