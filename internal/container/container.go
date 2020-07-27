package container

import (
	log "github.com/sirupsen/logrus"
	"meli/app"
	"meli/app/caching"
	"meli/app/status"
	"meli/internal/api"
	pg "meli/internal/postgres"
	rd "meli/internal/redis"
	config "meli/pkg/config"
)

type ControllerGroup struct {
	StatusController status.StatusController
	ItemController   caching.ItemController
}

func Build() ControllerGroup {
	configs := config.NewConfig()

	// storage
	postgres := pg.NewPostgres(configs)
	redis := rd.NewRedis(configs)
	log.Info("redis: ", redis)

	// http
	httpClient := api.NewHttpClient(configs)

	// repositories
	itemRepository := app.NewItemRepository(postgres)

	// services
	itemService := caching.NewItemService(httpClient, itemRepository)

	// controllers
	group := ControllerGroup{}
	group.StatusController = status.NewStatusController(configs)
	group.ItemController = caching.NewItemController(configs, itemService)

	return group
}
