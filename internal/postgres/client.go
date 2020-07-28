package postgres

import (
	"database/sql"
	"fmt"
	_ "github.com/Kount/pq-timeouts"
	log "github.com/sirupsen/logrus"
	config "meli/pkg/config"
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
