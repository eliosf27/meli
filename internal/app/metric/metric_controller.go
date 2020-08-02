package metric

import (
	"github.com/labstack/echo/v4"
	config "meli/pkg/config"
	"net/http"
)

type (
	Controller struct {
		Config        config.Config
		MetricService MetricService
	}
)

func NewMetricController(Config config.Config, ItemService MetricService) MetricController {
	return Controller{
		Config:        Config,
		MetricService: ItemService,
	}
}

func (c Controller) Health(ctx echo.Context) error {

	return ctx.JSON(http.StatusOK, c.MetricService.FetchMetrics())
}

type MetricController interface {
	Health(c echo.Context) error
}
