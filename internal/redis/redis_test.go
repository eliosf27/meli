package redis

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	config "meli/pkg/config"
	"meli/pkg/testcontainers"
	"time"
)

var _ = Describe("Redis", func() {
	var redisContainer testcontainers.RedisContainer
	BeforeEach(func() {
		redisContainer = testcontainers.NewRedisContainer()
	})
	AfterEach(func() {
		_ = redisContainer.Down()
	})

	Context("using redis cache", func() {
		When("saving data with a set operation", func() {
			It("should return a valid operation", func() {
				connection := redisContainer.Up()

				configs := config.NewConfig()
				configs.Redis.RedisHost = connection.Host
				configs.Redis.RedisPort = connection.Port

				redis := NewRedis(configs)

				key := "test_key"
				value := []byte("test_value")
				_, err := redis.Set(key, value, time.Second*10)
				Expect(err).To(BeNil())
			})
		})

		When("fetching data from a get operation", func() {
			It("should return a valid value", func() {
				connection := redisContainer.Up()

				configs := config.NewConfig()
				configs.Redis.RedisHost = connection.Host
				configs.Redis.RedisPort = connection.Port

				redis := NewRedis(configs)

				key := "test_key"
				value := []byte("test_value")
				_, err := redis.Set(key, value, time.Second*10)
				Expect(err).To(BeNil())

				result, err := redis.Get(key)
				Expect(err).To(BeNil())
				Expect(result).To(Equal(value))
			})
		})
	})
})
