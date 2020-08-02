package container

import (
	log "github.com/sirupsen/logrus"
	"meli/internal/app/item"
	"meli/internal/app/status"
	"meli/internal/http"
	pg "meli/internal/postgres"
	"meli/internal/queue"
	rd "meli/internal/redis"
	config "meli/pkg/config"
)

type Dependencies struct {
	StatusController status.StatusController
	ItemController   item.ItemController
	Config           config.Config
	Queue            *queue.ItemQueue
}

// Build build all the project dependencies
func Build() Dependencies {
	configs := config.NewConfig()

	// storage
	postgres := pg.NewPostgres(configs)
	redis := rd.NewRedis(configs)
	log.Info("redis: ", redis)

	// queues
	itemQueue := queue.NewItemQueue()

	// httpClient
	httpClient := http.NewHttpClient(configs, &itemQueue)

	// http services
	itemHttpService := http.NewItemHttpService(&httpClient)

	// repositories
	itemRepository := item.NewItemRepository(postgres)

	// services
	itemService := item.NewItemService(itemHttpService, itemRepository)

	// server dependencies
	dependencies := Dependencies{}
	dependencies.StatusController = status.NewStatusController(configs)
	dependencies.ItemController = item.NewItemController(configs, itemService)
	dependencies.Config = configs
	dependencies.Queue = &itemQueue

	return dependencies
}
