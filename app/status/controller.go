package status

import (
	"github.com/labstack/echo/v4"
	config "meli/pkg/config"
	"net/http"
)

type Response struct {
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
	return ctx.JSON(http.StatusOK, Response{
		Version: c.Config.ProjectVersion,
		Name:    c.Config.ProjectName,
	})
}

type StatusController interface {
	Status(c echo.Context) error
}
