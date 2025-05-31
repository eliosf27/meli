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
	settings "meli/pkg/config" // Renamed to avoid conflict with our es.Config
	es "meli/pkg/elasticsearch"
	"log" // For the Elasticsearch client logger
	"os"  // For the Elasticsearch client logger
)

type Dependencies struct {
	StatusHandler       appStatus.StatusHandler
	ItemHandler         appItem.ItemHandler
	MetricHandler       appMetric.MetricHandler
	ConfigCacheHandler  appConfig.ConfigCacheHandler
	Config              settings.Config
	Queue               *queue.ItemQueue
	QueueConsumers      []queue.Consumer
	ElasticsearchClient *es.Client // Added Elasticsearch client
}

// Build build all the project dependencies
func Build() Dependencies {
	configs := settings.NewConfig()

	// Elasticsearch Client
	esLogger := log.New(os.Stdout, "elasticsearch: ", log.LstdFlags)
	esCfg := es.Config{
		Addresses: configs.Elasticsearch.Addresses,
		Username:  configs.Elasticsearch.Username,
		Password:  configs.Elasticsearch.Password,
	}
	esClient, err := es.NewClient(esCfg, esLogger)
	if err != nil {
		// Decide on error handling: panic, log and exit, or return error
		// For now, let's panic, similar to how config loading errors are handled.
		panic("Failed to create Elasticsearch client: " + err.Error())
	}
	// TODO: Optionally, add a Ping here to verify connection on startup
	// if err := esClient.Ping(); err != nil {
	// 	panic("Failed to ping Elasticsearch: " + err.Error())
	// }


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
	dependencies.ElasticsearchClient = esClient // Added Elasticsearch client

	return dependencies
}
