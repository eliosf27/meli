package metric

import (
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"meli/internal/entities"
	"meli/internal/redis"
	configPkg "meli/pkg/config"
	"meli/pkg/testcontainers"
	timePkg "meli/pkg/time"
)

var _ = Describe("MetricService", func() {
	var ctrl *gomock.Controller
	var container testcontainers.RedisContainer

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
		container = testcontainers.NewRedisContainer()
	})
	AfterEach(func() {
		ctrl.Finish()
		container.Down()
	})

	Context("calling the metric service", func() {
		When("update the metrics", func() {
			It("should return a valid response", func() {
				connection := container.Up()

				configs := configPkg.NewConfig()
				configs.Redis.RedisHost = connection.Host
				configs.Redis.RedisPort = connection.Port

				redis := redis.NewRedis(configs)

				itemMetrics := entities.ItemMetric{
					Type:         entities.ExternalApi,
					ResponseTime: 10,
					StatusCode:   200,
					Time:         timePkg.Now(),
				}
				metricsService := NewMetricService(redis)
				err := metricsService.UpdateMetric(itemMetrics)

				Expect(err).To(BeNil())
			})
		})

		When("update external metrics", func() {
			It("should return a valid response", func() {
				connection := container.Up()

				configs := configPkg.NewConfig()
				configs.Redis.RedisHost = connection.Host
				configs.Redis.RedisPort = connection.Port

				redis := redis.NewRedis(configs)

				itemMetrics := entities.ItemMetric{
					Type:         entities.ExternalApi,
					ResponseTime: 10,
					StatusCode:   200,
					Time:         timePkg.Now(),
				}
				metricsService := NewMetricService(redis)
				err := metricsService.UpdateMetric(itemMetrics)
				Expect(err).To(BeNil())

				result := metricsService.FetchMetrics()

				Expect(err).To(BeNil())
				Expect(len(result)).To(Equal(1))
				Expect(result[0].TotalRequests).To(Equal(int64(0)))
				Expect(result[0].TotalCountApiCalls).To(Equal(int64(1)))
				Expect(result[0].AvgResponseTime).To(Equal(0.0))
				Expect(result[0].AvgResponseTimeApiCall).To(Equal(float64(itemMetrics.ResponseTime)))
			})
		})

		When("update local metrics", func() {
			It("should return a valid response", func() {
				connection := container.Up()

				configs := configPkg.NewConfig()
				configs.Redis.RedisHost = connection.Host
				configs.Redis.RedisPort = connection.Port

				redis := redis.NewRedis(configs)

				itemMetrics := entities.ItemMetric{
					Type:         entities.LocalApi,
					ResponseTime: 10,
					StatusCode:   200,
					Time:         timePkg.Now(),
				}
				metricsService := NewMetricService(redis)
				err := metricsService.UpdateMetric(itemMetrics)
				Expect(err).To(BeNil())

				result := metricsService.FetchMetrics()

				Expect(err).To(BeNil())
				Expect(len(result)).To(Equal(1))
				Expect(result[0].TotalCountApiCalls).To(Equal(int64(0)))
				Expect(result[0].TotalRequests).To(Equal(int64(1)))
				Expect(result[0].AvgResponseTimeApiCall).To(Equal(0.0))
				Expect(result[0].AvgResponseTime).To(Equal(float64(itemMetrics.ResponseTime)))
			})
		})

		When("update local and metrics at the same time", func() {
			It("should return a valid response", func() {
				connection := container.Up()

				configs := configPkg.NewConfig()
				configs.Redis.RedisHost = connection.Host
				configs.Redis.RedisPort = connection.Port

				redis := redis.NewRedis(configs)
				metricsService := NewMetricService(redis)

				localMetrics := entities.ItemMetric{
					Type:         entities.LocalApi,
					ResponseTime: 5,
					StatusCode:   500,
					Time:         timePkg.Now(),
				}

				err := metricsService.UpdateMetric(localMetrics)
				Expect(err).To(BeNil())

				externalMetrics := entities.ItemMetric{
					Type:         entities.ExternalApi,
					ResponseTime: 50,
					StatusCode:   200,
					Time:         timePkg.Now(),
				}

				err = metricsService.UpdateMetric(externalMetrics)
				Expect(err).To(BeNil())

				result := metricsService.FetchMetrics()

				Expect(err).To(BeNil())
				Expect(len(result)).To(Equal(1))
				Expect(result[0].TotalCountApiCalls).To(Equal(int64(1)))
				Expect(result[0].TotalRequests).To(Equal(int64(1)))
				Expect(result[0].AvgResponseTimeApiCall).To(Equal(float64(externalMetrics.ResponseTime)))
				Expect(result[0].AvgResponseTime).To(Equal(float64(localMetrics.ResponseTime)))
			})
		})

		When("update two local metrics at the same time", func() {
			It("should return a valid response", func() {
				connection := container.Up()

				configs := configPkg.NewConfig()
				configs.Redis.RedisHost = connection.Host
				configs.Redis.RedisPort = connection.Port

				redis := redis.NewRedis(configs)
				metricsService := NewMetricService(redis)

				localMetrics1 := entities.ItemMetric{
					Type:         entities.LocalApi,
					ResponseTime: 5,
					StatusCode:   500,
					Time:         timePkg.Now(),
				}

				err := metricsService.UpdateMetric(localMetrics1)
				Expect(err).To(BeNil())

				localMetrics2 := entities.ItemMetric{
					Type:         entities.LocalApi,
					ResponseTime: 50,
					StatusCode:   200,
					Time:         timePkg.Now(),
				}

				err = metricsService.UpdateMetric(localMetrics2)
				Expect(err).To(BeNil())

				result := metricsService.FetchMetrics()

				avgLocal := (float64(localMetrics1.ResponseTime) + float64(localMetrics2.ResponseTime)) / 2
				Expect(err).To(BeNil())
				Expect(len(result)).To(Equal(1))
				Expect(result[0].TotalCountApiCalls).To(Equal(int64(0)))
				Expect(result[0].TotalRequests).To(Equal(int64(2)))
				Expect(result[0].AvgResponseTimeApiCall).To(Equal(float64(0)))
				Expect(result[0].AvgResponseTime).To(Equal(avgLocal))
			})
		})

		When("update two external metrics at the same time", func() {
			It("should return a valid response", func() {
				connection := container.Up()

				configs := configPkg.NewConfig()
				configs.Redis.RedisHost = connection.Host
				configs.Redis.RedisPort = connection.Port

				redis := redis.NewRedis(configs)
				metricsService := NewMetricService(redis)

				externalMetrics1 := entities.ItemMetric{
					Type:         entities.ExternalApi,
					ResponseTime: 5,
					StatusCode:   500,
					Time:         timePkg.Now(),
				}

				err := metricsService.UpdateMetric(externalMetrics1)
				Expect(err).To(BeNil())

				externalMetrics2 := entities.ItemMetric{
					Type:         entities.ExternalApi,
					ResponseTime: 50,
					StatusCode:   200,
					Time:         timePkg.Now(),
				}

				err = metricsService.UpdateMetric(externalMetrics2)
				Expect(err).To(BeNil())

				result := metricsService.FetchMetrics()

				avgExternal := (float64(externalMetrics1.ResponseTime) + float64(externalMetrics2.ResponseTime)) / 2
				Expect(err).To(BeNil())
				Expect(len(result)).To(Equal(1))
				Expect(result[0].TotalRequests).To(Equal(int64(0)))
				Expect(result[0].TotalCountApiCalls).To(Equal(int64(2)))
				Expect(result[0].AvgResponseTime).To(Equal(float64(0)))
				Expect(result[0].AvgResponseTimeApiCall).To(Equal(avgExternal))
				Expect(len(result[0].InfoRequests)).To(Equal(2))
				Expect(len(result[0].InfoRequests)).To(Equal(2))
				Expect(result[0].InfoRequests[0].StatusCode).To(Equal(int64(200)))
				Expect(result[0].InfoRequests[0].Count).To(Equal(int64(1)))
				Expect(result[0].InfoRequests[1].StatusCode).To(Equal(int64(500)))
				Expect(result[0].InfoRequests[1].Count).To(Equal(int64(1)))
			})
		})
	})
})
