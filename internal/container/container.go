package container

import (
	log "github.com/sirupsen/logrus"
	"meli/app/status"
	pg "meli/internal/postgres"
	rd "meli/internal/redis"
	config "meli/pkg/config"
)

type ControllerGroup struct {
	StatusController status.StatusController
}

func Build() ControllerGroup {
	configs := config.NewConfig()

	postgres := pg.NewPostgres(configs)
	log.Info("postgres: ", postgres)

	redis := rd.NewRedis(configs)
	log.Info("redis: ", redis)

	group := ControllerGroup{}
	group.StatusController = status.NewStatusController(configs)

	return group
}
