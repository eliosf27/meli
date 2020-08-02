package item

import (
	"errors"
	"github.com/labstack/echo/v4"
	config "meli/pkg/config"
	"meli/pkg/strings"
	"net/http"
)

type ItemHandler interface {
	Get(c echo.Context) error
}

type ItemHandle struct {
	config      config.Config
	itemService ItemService
}

func NewItemHandle(Config config.Config, ItemService ItemService) ItemHandler {
	return ItemHandle{
		config:      Config,
		itemService: ItemService,
	}
}

func (c ItemHandle) Get(ctx echo.Context) error {
	id := ctx.Param("item_id")
	if strings.IsEmpty(id) {
		return echo.NewHTTPError(http.StatusBadRequest, errors.New("invalid item id"))
	}

	return ctx.JSON(http.StatusOK, c.itemService.FetchItemById(id))
}
