package postgres

import (
	"database/sql"
	"fmt"
	_ "github.com/Kount/pq-timeouts"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	log "github.com/sirupsen/logrus"
	config "meli/pkg/config"
	"meli/pkg/files"
	"path/filepath"
	"time"
)

type Postgres struct {
	Client *sql.DB
}

func NewPostgres(config config.Config) Postgres {
	client, err := buildClient(config.Postgres.ItemsConnection)
	if err != nil {
		panic(fmt.Sprintf("Error connecting to postgres server: %s", err.Error()))
	}

	err = client.Ping()
	if err != nil {
		log.Info("PostgresDB => Monitoring | Cannot connect to postgres server | Error => ", err)
	} else {
		log.Info("PostgresDB => Monitoring | Connected successfully")
	}

	return Postgres{
		Client: client,
	}
}

// buildClient build a new postgres client with its configurations
func buildClient(connection string) (*sql.DB, error) {
	url := "%s?sslmode=disable&connect_timeout=%d&read_timeout=%d&write_timeout=%d"

	connectionString := fmt.Sprintf(url,
		connection,
		5,
		5000,
		5000,
	)

	db, err := sql.Open("pq-timeouts", connectionString)
	if err != nil {
		return db, err
	}

	db.SetMaxIdleConns(5)
	db.SetMaxOpenConns(5)
	db.SetConnMaxLifetime(time.Minute * 2)

	return db, nil
}

// RunMigrations execute the postgres migration of the project
func (p Postgres) RunMigrations() {
	projectPath := files.GetProjectPath()
	filePath := filepath.Join("file://", projectPath, "internal/postgres/migrations")

	driver, err := postgres.WithInstance(p.Client, &postgres.Config{})
	if err != nil {
		panic(fmt.Sprintf("RunMigrations | Error creating DB instance: %s", err.Error()))
	}

	m, err := migrate.NewWithDatabaseInstance(filePath, "postgres", driver)
	if err != nil {
		panic(fmt.Sprintf("RunMigrations | Error creating migration instance: %s", err.Error()))
	}

	if err = m.Up(); err != nil {
		log.Warning("RunMigrations | Error running migration: ", err)
	} else {
		log.Info("RunMigrations | Database migrations processed...")
	}
}
