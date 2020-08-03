package config

import (
	"github.com/labstack/echo/v4"
	"meli/internal/request"
	"net/http"
)

type ConfigHandler interface {
	Fetch(c echo.Context) error
	Update(ctx echo.Context) error
}

type ConfigHandle struct {
	configService ConfigServicer
}

func NewConfigHandle(service ConfigServicer) ConfigHandler {
	return ConfigHandle{
		configService: service,
	}
}

func (c ConfigHandle) Fetch(ctx echo.Context) error {
	val, err := c.configService.Fetch()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(http.StatusOK, val)
}

func (c ConfigHandle) Update(ctx echo.Context) error {
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
