package config

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
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

	Context("calling the config service", func() {
		When("the storage is fetching", func() {
			It("should return a valid operation", func() {
				connection := container.Up()

				configs := config.NewConfig()
				configs.Redis.RedisHost = connection.Host
				configs.Redis.RedisPort = connection.Port

				storage := "redis"
				redis := redis.NewRedis(configs)
				configService := NewConfigService(redis)

				err := configService.Update(storage)
				Expect(err).To(BeNil())

				result, err := configService.Fetch()
				Expect(err).To(BeNil())
				Expect(storage).To(Equal(result))
			})
		})

		When("the storage is fetching", func() {
			It("should return an error", func() {
				connection := container.Up()

				configs := config.NewConfig()
				configs.Redis.RedisHost = connection.Host
				configs.Redis.RedisPort = connection.Port

				redis := redis.NewRedis(configs)
				configService := NewConfigService(redis)

				result, err := configService.Fetch()
				Expect(err).To(Not(BeNil()))
				Expect(err.Error()).To(Equal("configuracion not found"))
				Expect(result).To(Equal(""))
			})
		})

		When("the storage is updated to redis storage", func() {
			It("should return a valid operation", func() {
				connection := container.Up()

				configs := config.NewConfig()
				configs.Redis.RedisHost = connection.Host
				configs.Redis.RedisPort = connection.Port

				storage := "redis"
				redis := redis.NewRedis(configs)
				configService := NewConfigService(redis)
				err := configService.Update(storage)
				Expect(err).To(BeNil())
			})
		})

		When("the storage is updated to postgres storage", func() {
			It("should return a valid operation", func() {
				connection := container.Up()

				configs := config.NewConfig()
				configs.Redis.RedisHost = connection.Host
				configs.Redis.RedisPort = connection.Port

				storage := "postgres"
				redis := redis.NewRedis(configs)
				configService := NewConfigService(redis)
				err := configService.Update(storage)
				Expect(err).To(BeNil())
			})
		})
	})
})
