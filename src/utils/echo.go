package utils

import (
	"github.com/labstack/echo/v5"
)

func InitEchoApp() *echo.Echo {
	e := echo.New()
	return e
}
