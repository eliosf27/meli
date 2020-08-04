package postgres

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"meli/internal/entities"
	config "meli/pkg/config"
	"meli/pkg/testcontainers"
)

var _ = Describe("Postgres", func() {
	var postgresContainer testcontainers.PostgresContainer
	BeforeEach(func() {
		postgresContainer = testcontainers.NewPostgresContainer()
	})
	AfterEach(func() {
		_ = postgresContainer.Down()
	})

	Context("using postgres database", func() {
		When("fetching data from postgres database", func() {
			It("should return a valid value", func() {
				connection := postgresContainer.Up()

				configs := config.NewConfig()
				configs.Postgres.ItemsConnection = connection.ConnectionString

				postgres := NewPostgres(configs)
				postgres.RunMigrations()

				itemId := "xxx"
				item := entities.Item{ItemId: itemId}
				err := postgres.Execute(
					"insert_item",
					item.ItemId, item.Title, item.CategoryId,
					item.Price, item.StartTime, item.StopTime,
				)
				Expect(err).To(BeNil())

				rows, err := postgres.Query("select_items_children", itemId)
				Expect(err).To(BeNil())
				Expect(rows.Next()).To(Equal(true))
			})
		})
	})

	Context("using postgres database", func() {
		When("saving data in postgres", func() {
			It("should return a valid operation", func() {
				connection := postgresContainer.Up()

				configs := config.NewConfig()
				configs.Postgres.ItemsConnection = connection.ConnectionString

				postgres := NewPostgres(configs)
				postgres.RunMigrations()

				itemId := "xxx"
				item := entities.Item{ItemId: itemId}
				err := postgres.Execute(
					"insert_item",
					item.ItemId, item.Title, item.CategoryId,
					item.Price, item.StartTime, item.StopTime,
				)
				Expect(err).To(BeNil())
			})
		})
	})
})
