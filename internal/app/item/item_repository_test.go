package item_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	log "github.com/sirupsen/logrus"
	"meli/app/item"
	mocks "meli/internal/mocks"
	pg "meli/internal/postgres"
	config "meli/pkg/config"
	"meli/pkg/testcontainers"
	"testing"
)

func TestItem22(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Item Suite")
}

var _ = Describe("ItemRepository", func() {
	var container testcontainers.PostgresContainer
	BeforeEach(func() {
		container = testcontainers.NewPostgresContainer()
	})
	AfterEach(func() {
		_ = container.Down()
	})

	Context("calling the item repository", func() {
		When("the item is saving", func() {
			It("should return a valid operation", func() {
				connection := container.Up()

				configs := config.NewConfig()
				configs.Postgres.ItemsConnection = connection.ConnectionString

				postgres := pg.NewPostgres(configs)
				postgres.RunMigrations()

				repository := item.NewItemRepository(postgres)
				log.Info(repository)

				err := repository.Save(mocks.MockItem("xxx"))

				Expect(err).To(BeNil())
			})
		})

		When("the item is fetching", func() {
			It("should return a valid item", func() {
				itemId := "xxx"
				connection := container.Up()

				configs := config.NewConfig()
				configs.Postgres.ItemsConnection = connection.ConnectionString

				postgres := pg.NewPostgres(configs)
				postgres.RunMigrations()

				repository := item.NewItemRepository(postgres)
				log.Info(repository)

				err := repository.Save(mocks.MockItem(itemId))
				Expect(err).To(BeNil())

				item, err := repository.Get(itemId)
				Expect(err).To(BeNil())
				Expect(item.ItemId).To(Equal(itemId))
				Expect(len(item.Children)).To(Equal(0))
			})
		})

		When("the item is fetching and has children", func() {
			It("should return a valid item and a valid item children", func() {
				itemId := "xxx"
				connection := container.Up()

				configs := config.NewConfig()
				configs.Postgres.ItemsConnection = connection.ConnectionString

				postgres := pg.NewPostgres(configs)
				postgres.RunMigrations()

				repository := item.NewItemRepository(postgres)
				log.Info(repository)

				mockItem := mocks.MockItem(itemId)
				mockItem.Children = mocks.MockItemChildren(itemId)
				err := repository.Save(mockItem)
				Expect(err).To(BeNil())

				item, err := repository.Get(itemId)
				Expect(err).To(BeNil())
				Expect(item.ItemId).To(Equal(itemId))
				Expect(len(item.Children)).To(Equal(2))
			})
		})
	})
})
