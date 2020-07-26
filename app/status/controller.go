package status

import (
	"github.com/labstack/echo/v4"
	config "meli/pkg/config"
	"net/http"
)

type Status struct {
	Version string `json:"version"`
	Name    string `json:"name"`
}

type (
	Controller struct {
		Config config.Config
	}
)

func NewStatusController(Config config.Config) StatusController {
	return Controller{
		Config: Config,
	}
}

func (c Controller) Status(ctx echo.Context) error {
	return ctx.JSON(http.StatusOK, Status{
		Version: c.Config.ProjectVersion,
		Name:    c.Config.ProjectName,
	})
}
