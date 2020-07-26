package container

import (
	log "github.com/sirupsen/logrus"
	"meli/app/status"
	dal "meli/internal/postgres"
)

type ControllerGroup struct {
	StatusController status.StatusController
}

func Build() ControllerGroup {
	postgres := dal.NewPostgres("")
	log.Info("postgres: ", postgres.Client)

	group := ControllerGroup{}
	group.StatusController = status.NewStatusController()

	return group
}
