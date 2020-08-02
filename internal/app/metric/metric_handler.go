package metric

import (
	"github.com/labstack/echo/v4"
	config "meli/pkg/config"
	"net/http"
)

type MetricHandler interface {
	Health(c echo.Context) error
}

type MetricHandle struct {
	config        config.Config
	metricService MetricService
}

func NewMetricHandler(Config config.Config, ItemService MetricService) MetricHandler {
	return MetricHandle{
		config:        Config,
		metricService: ItemService,
	}
}

func (c MetricHandle) Health(ctx echo.Context) error {

	return ctx.JSON(http.StatusOK, c.metricService.FetchMetrics())
}
