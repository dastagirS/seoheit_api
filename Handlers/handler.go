package Handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func RouteHandlers(c echo.Context) error {


	return c.JSON(http.StatusOK, "success")
}
