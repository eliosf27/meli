package item_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	itemApp "meli/internal/app/item"
	itemMocks "meli/internal/mocks/item"
	pg "meli/internal/postgres"
	"meli/internal/redis"
	configPkg "meli/pkg/config"
	"meli/pkg/testcontainers"
)

var _ = Describe("ItemCache", func() {
	var redisContainer testcontainers.RedisContainer
	var postgresContainer testcontainers.PostgresContainer
	BeforeEach(func() {
		redisContainer = testcontainers.NewRedisContainer()
		postgresContainer = testcontainers.NewPostgresContainer()
	})
	AfterEach(func() {
		_ = redisContainer.Down()
		_ = postgresContainer.Down()
	})

	Context("calling the item cache strategy", func() {
		When("the no strategy is configured", func() {
			It("should use a default strategy and return valid values from postgres cache", func() {
				connection := postgresContainer.Up()

				configs := configPkg.NewConfig()
				configs.Postgres.ItemsConnection = connection.ConnectionString

				postgres := pg.NewPostgres(configs)
				postgres.RunMigrations()

				redis := redis.NewRedis(configs)

				itemId := "xxx"
				item := itemMocks.MockItem(itemId)
				itemRedisCache := itemApp.NewItemRedisCache(redis)
				itemPostgresCache := itemApp.NewItemPostgresCache(postgres)
				itemCache := itemApp.NewItemCache(itemRedisCache, itemPostgresCache)

				itemCache.SetStrategy(itemCache.GetStrategy("xxx"))
				err := itemCache.Save(item)
				Expect(err).To(BeNil())

				itemExpected, err := itemCache.Get(itemId)

				Expect(err).To(BeNil())
				Expect(itemExpected.ItemId).To(Equal(itemId))
			})
		})

		When("the postgres strategy is configured", func() {
			It("should return valid values from cache", func() {
				connection := postgresContainer.Up()

				configs := configPkg.NewConfig()
				configs.Postgres.ItemsConnection = connection.ConnectionString

				postgres := pg.NewPostgres(configs)
				postgres.RunMigrations()

				itemId := "xxx"
				item := itemMocks.MockItem(itemId)
				itemPostgresCache := itemApp.NewItemPostgresCache(postgres)
				itemCache := itemApp.NewItemCache(nil, itemPostgresCache)

				itemCache.SetStrategy(itemPostgresCache)
				err := itemCache.Save(item)
				Expect(err).To(BeNil())

				itemExpected, err := itemCache.Get(itemId)

				Expect(err).To(BeNil())
				Expect(itemExpected.ItemId).To(Equal(itemId))
			})
		})

		When("the redis strategy is configured", func() {
			It("should return valid values from cache", func() {
				connection := redisContainer.Up()

				configs := configPkg.NewConfig()
				configs.Redis.RedisHost = connection.Host
				configs.Redis.RedisPort = connection.Port

				redis := redis.NewRedis(configs)

				itemId := "xxx"
				item := itemMocks.MockItem(itemId)
				itemRedisCache := itemApp.NewItemRedisCache(redis)
				itemCache := itemApp.NewItemCache(nil, itemRedisCache)

				itemCache.SetStrategy(itemRedisCache)
				err := itemCache.Save(item)
				Expect(err).To(BeNil())

				itemExpected, err := itemCache.Get(itemId)

				Expect(err).To(BeNil())
				Expect(itemExpected.ItemId).To(Equal(itemId))
			})
		})
	})
})
