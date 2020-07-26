package dal

import (
	"database/sql"
	"fmt"
	_ "github.com/Kount/pq-timeouts"
)

type Postgres struct {
	client  *sql.DB
	address string
}

func NewPostgres(address string) Postgres {
	client, err := buildClient(address)
	if err != nil {
		panic(fmt.Sprintf("Error connecting to DB: %s", err.Error()))
	}

	return Postgres{
		client:  client,
		address: address,
	}
}

func (p *Postgres) Client() Postgres {
	client, err := buildClient(p.address)
	if err != nil {
		panic(fmt.Sprintf("Error connecting to DB: %s", err.Error()))
	}

	return Postgres{
		client:  client,
		address: p.address,
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
