package status

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

type (
	Controller struct{}
)

func NewStatusController() StatusController {
	return Controller{}
}

func (Controller) Status(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, World!")
}
