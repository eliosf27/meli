package item

import (
	"github.com/labstack/echo/v4"
	"meli/internal/request"
	"net/http"
)

type ConfigHandler interface {
	Fetch(c echo.Context) error
	Save(ctx echo.Context) error
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
	req := new(request.FindConfig)
	if err := ctx.Bind(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	if err := ctx.Validate(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return ctx.JSON(http.StatusOK, req)
}

func (c ConfigHandle) Save(ctx echo.Context) error {
	req := new(request.SaveConfig)
	if err := ctx.Bind(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	if err := ctx.Validate(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return ctx.JSON(http.StatusOK, req)
}
