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
	"meli/pkg/queries"
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
func (p *Postgres) RunMigrations() {
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

func (p *Postgres) Execute(sql string, args ...interface{}) error {
	query, err := queries.ReadQuery(sql)
	if err != nil {

		return err
	}

	_, err = p.Client.Exec(query, args)
	if err != nil {

		return err
	}

	return nil
}

type Rows sql.Rows

func (p *Postgres) Query(sql string, params ...interface{}) (*sql.Rows, error) {
	query, err := queries.ReadQuery(sql)
	if err != nil {
		log.Errorf("ItemRepository.SaveChildren | Error reading query: %+v", err)

		return nil, err
	}

	rows, err := p.Client.Query(query, params)
	if err != nil {
		return nil, err
	}

	defer func() {
		if err := rows.Close(); err != nil {
			log.Errorf("Error closing rows: %+v", err)
		}
	}()

	//for rows.Next() {
	//	err := rows.Scan(values)
	//	if err != nil {
	//
	//	}
	//}

	return rows, nil
}
