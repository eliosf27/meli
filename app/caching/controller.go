package caching

import (
	"errors"
	"github.com/labstack/echo/v4"
	"meli/kit/strings"
	config "meli/pkg/config"
	"net/http"
)

type (
	Controller struct {
		Config      config.Config
		ItemService ItemService
	}
)

func NewItemController(Config config.Config, ItemService ItemService) ItemController {
	return Controller{
		Config:      Config,
		ItemService: ItemService,
	}
}

func (c Controller) Get(ctx echo.Context) error {
	id := ctx.QueryParam("item_id")
	if strings.IsEmpty(id) {
		return echo.NewHTTPError(http.StatusBadRequest, errors.New("invalid item id"))
	}

	return ctx.JSON(http.StatusOK, c.ItemService.FetchItemByID(id))
}

type ItemController interface {
	Get(c echo.Context) error
}
