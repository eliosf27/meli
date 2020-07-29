package testcontainers

import (
	"context"
	"fmt"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

const (
	port     = "5432"
	user     = "docker"
	password = "docker"
	database = "items_test"
)

type PostgresConnection struct {
	Host             string
	Port             string
	User             string
	Password         string
	ConnectionString string
}

type PostgresContainer struct {
	Container testcontainers.Container
}

func NewPostgresContainer() PostgresContainer {
	return PostgresContainer{
		Container: nil,
	}
}

// Up create a new postgres container with the connections params
func (c PostgresContainer) Up() PostgresConnection {
	containerRequest := testcontainers.ContainerRequest{
		Image:        "postgres",
		ExposedPorts: []string{fmt.Sprintf("%s/tcp", port)},
		WaitingFor:   wait.ForListeningPort(port),
		Env: map[string]string{
			"POSTGRES_USER":     user,
			"POSTGRES_PASSWORD": password,
			"POSTGRES_DB":       database,
		},
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
	containerPort, _ := container.MappedPort(context.Background(), port)
	connectionString := "postgres://%s:%s@%s:%s/%s"

	return PostgresConnection{
		Host:             containerHost,
		Port:             containerPort.Port(),
		User:             user,
		Password:         password,
		ConnectionString: fmt.Sprintf(connectionString, user, password, containerHost, containerPort.Port(), database),
	}
}

// Down destroy the postgres container
func (c PostgresContainer) Down() error {
	if c.Container != nil {
		err := c.Container.Terminate(context.Background())
		if err != nil {
			return err
		}
	}

	return nil
}
