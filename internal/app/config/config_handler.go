package config

import (
	"github.com/labstack/echo/v4"
	"meli/internal/request"
	"net/http"
)

type ConfigCacheHandler interface {
	Fetch(c echo.Context) error
	Update(ctx echo.Context) error
}

type ConfigCacheHandle struct {
	configService ConfigServicer
}

func NewConfigCacheHandle(service ConfigServicer) ConfigCacheHandler {
	return ConfigCacheHandle{
		configService: service,
	}
}

func (c ConfigCacheHandle) Fetch(ctx echo.Context) error {
	val, err := c.configService.Fetch()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(http.StatusOK, val)
}

func (c ConfigCacheHandle) Update(ctx echo.Context) error {
	req := new(request.UpdateStorage)
	if err := ctx.Bind(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	if err := ctx.Validate(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	err := c.configService.Update(req.Storage)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(http.StatusOK, "storage updated")
}
