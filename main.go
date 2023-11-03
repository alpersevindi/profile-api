package main

import (
	"profile-api/routers"

	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()
	routers.InitRoutes(e)
	e.Logger.Fatal(e.Start(":8080"))
}
