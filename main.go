package main

import (
	"profile-api/database"
	"profile-api/routers"

	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()

	database.InitDynamoDB()
	routers.InitRoutes(e)
	e.Logger.Fatal(e.Start(":8080"))
}
