package postgres

import (
	"database/sql"
	"fmt"
	_ "github.com/Kount/pq-timeouts"
	log "github.com/sirupsen/logrus"
	config "meli/pkg/config"
)

type Postgres struct {
	Client *sql.DB
}

func NewPostgres(config config.Config) Postgres {
	client, err := buildClient(config.Postgres.ItemsConnection)
	if err != nil {
		panic(fmt.Sprintf("Error connecting to DB: %s", err.Error()))
	} else {
		log.Info("Connected to postgres server")
	}

	return Postgres{
		Client: client,
	}
}

func buildClient(connection string) (*sql.DB, error) {
	url := "%s?sslmode=disable&connect_timeout=%d&read_timeout=%d&write_timeout=%d"

	connectionString := fmt.Sprintf(url,
		connection,
		10,
		10,
		10,
	)

	db, err := sql.Open("pq-timeouts", connectionString)
	if err != nil {
		return db, err
	}

	db.SetMaxIdleConns(10)
	db.SetMaxOpenConns(10)
	db.SetConnMaxLifetime(10)

	return db, nil
}
