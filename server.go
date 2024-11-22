package main

import (
	"github.com/introdevio/receipt_processor/api"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	e := echo.New()
	e.Use(middleware.CORS())
	api.NewReceiptsApi(e)
	e.Logger.Fatal(e.Start(":8080"))
}
