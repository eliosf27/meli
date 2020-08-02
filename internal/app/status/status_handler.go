package status

import (
	"github.com/labstack/echo/v4"
	config "meli/pkg/config"
	"net/http"
)

type StatusHandler interface {
	Status(c echo.Context) error
}

type Response struct {
	Version string `json:"version"`
	Name    string `json:"name"`
}

type StatusHandle struct {
	config config.Config
}

func NewStatusHandler(Config config.Config) StatusHandler {
	return StatusHandle{
		config: Config,
	}
}

func (c StatusHandle) Status(ctx echo.Context) error {
	return ctx.JSON(http.StatusOK, Response{
		Version: c.config.ProjectVersion,
		Name:    c.config.ProjectName,
	})
}
