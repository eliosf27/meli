package container

import (
	log "github.com/sirupsen/logrus"
	"meli/internal/api"
	"meli/internal/app/item"
	"meli/internal/app/status"
	pg "meli/internal/postgres"
	rd "meli/internal/redis"
	config "meli/pkg/config"
)

type Dependencies struct {
	StatusController status.StatusController
	ItemController   item.ItemController
	Config           config.Config
}

// Build build all the project dependencies
func Build() Dependencies {
	configs := config.NewConfig()

	// storage
	postgres := pg.NewPostgres(configs)
	redis := rd.NewRedis(configs)
	log.Info("redis: ", redis)

	// httpserver
	httpClient := api.NewHttpClient(configs)

	// repositories
	itemRepository := item.NewItemRepository(postgres)

	// services
	itemService := item.NewItemService(httpClient, itemRepository)

	// controllers
	dependencies := Dependencies{}
	dependencies.StatusController = status.NewStatusController(configs)
	dependencies.ItemController = item.NewItemController(configs, itemService)
	dependencies.Config = configs

	return dependencies
}
