package testcontainers

import (
	"context"
	"fmt"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

const (
	postgresPort     = "5432"
	postgresUser     = "docker"
	postgresPassword = "docker"
	postgresDatabase = "items_test"
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
		ExposedPorts: []string{fmt.Sprintf("%s/tcp", postgresPort)},
		WaitingFor:   wait.ForListeningPort(postgresPort),
		Env: map[string]string{
			"POSTGRES_USER":     postgresUser,
			"POSTGRES_PASSWORD": postgresPassword,
			"POSTGRES_DB":       postgresDatabase,
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
	containerPort, _ := container.MappedPort(context.Background(), postgresPort)
	connectionString := "postgres://%s:%s@%s:%s/%s"

	return PostgresConnection{
		Host:     containerHost,
		Port:     containerPort.Port(),
		User:     postgresUser,
		Password: postgresPassword,
		ConnectionString: fmt.Sprintf(
			connectionString, postgresUser, postgresPassword,
			containerHost, containerPort.Port(), postgresDatabase,
		),
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
