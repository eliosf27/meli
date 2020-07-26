package status

import "github.com/labstack/echo/v4"

type StatusController interface {
	Status(c echo.Context) error
}
