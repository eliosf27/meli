package container

import (
	appConfig "meli/internal/app/config"
	appItem "meli/internal/app/item"
	appMetric "meli/internal/app/metric"
	appStatus "meli/internal/app/status"
	"meli/internal/http"
	pg "meli/internal/postgres"
	"meli/internal/queue"
	rd "meli/internal/redis"
	config "meli/pkg/config"
)

type Dependencies struct {
	StatusHandler      appStatus.StatusHandler
	ItemHandler        appItem.ItemHandler
	MetricHandler      appMetric.MetricHandler
	ConfigCacheHandler appConfig.ConfigCacheHandler
	Config             config.Config
	Queue              *queue.ItemQueue
	QueueConsumers     []queue.Consumer
}

// Build build all the project dependencies
func Build() Dependencies {
	configs := config.NewConfig()

	// storage
	postgres := pg.NewPostgres(configs)
	redis := rd.NewRedis(configs)

	// queues
	itemQueue := queue.NewItemQueue()
	metricsService := appMetric.NewMetricService(redis)
	itemConsumer := queue.NewItemConsumer(&itemQueue, metricsService)

	// httpClient
	httpClient := http.NewHttpClient(configs, &itemQueue)

	// http services
	itemHttpService := http.NewItemHttpService(&httpClient)

	// other services
	configService := appConfig.NewConfigService(redis)

	// caches
	itemPostgresCache := appItem.NewItemPostgresCache(postgres)
	itemRedisCache := appItem.NewItemRedisCache(redis)
	itemCache := appItem.NewItemCache(itemPostgresCache, itemRedisCache)

	// services
	itemService := appItem.NewItemService(itemHttpService, itemCache, configService)
	metricService := appMetric.NewMetricService(redis)

	// server dependencies
	dependencies := Dependencies{}
	dependencies.StatusHandler = appStatus.NewStatusHandler(configs)
	dependencies.ItemHandler = appItem.NewItemHandle(configs, itemService)
	dependencies.MetricHandler = appMetric.NewMetricHandler(configs, metricService)
	dependencies.ConfigCacheHandler = appConfig.NewConfigCacheHandle(configService)
	dependencies.Config = configs
	dependencies.Queue = &itemQueue
	dependencies.QueueConsumers = append(dependencies.QueueConsumers, itemConsumer)

	return dependencies
}
