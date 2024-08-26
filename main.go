package main

import (
	"net/http"
	"seonator-api/Methods"
	"github.com/labstack/echo/v4"
)

func main() {
	app := echo.New()
	app.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "seonator online")
	})
	app.GET("/*", Methods.GetSite)
	app.Start(":4000")
}
