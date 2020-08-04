package item_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"meli/internal/app/item"
	mocks "meli/internal/mocks"
	"meli/internal/redis"
	config "meli/pkg/config"
	"meli/pkg/testcontainers"
)

var _ = Describe("ItemRedisCache", func() {
	var container testcontainers.RedisContainer
	BeforeEach(func() {
		container = testcontainers.NewRedisContainer()
	})
	AfterEach(func() {
		_ = container.Down()
	})

	Context("calling the item redis cache", func() {
		When("the item is saving", func() {
			It("should return a valid operation", func() {
				connection := container.Up()

				configs := config.NewConfig()
				configs.Redis.RedisHost = connection.Host
				configs.Redis.RedisPort = connection.Port

				redis := redis.NewRedis(configs)
				repository := item.NewItemRedisCache(redis)

				err := repository.Save(mocks.MockItem("xxx"))

				Expect(err).To(BeNil())
			})
		})

		When("the item is fetching", func() {
			It("should return a valid item", func() {
				itemId := "xxx"
				connection := container.Up()

				configs := config.NewConfig()
				configs.Redis.RedisHost = connection.Host
				configs.Redis.RedisPort = connection.Port

				redis := redis.NewRedis(configs)
				repository := item.NewItemRedisCache(redis)

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
				configs.Redis.RedisHost = connection.Host
				configs.Redis.RedisPort = connection.Port

				redis := redis.NewRedis(configs)
				repository := item.NewItemRedisCache(redis)

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
