package routers

import (
	"profile-api/handlers"

	"github.com/labstack/echo/v4"
)

func InitRoutes(e *echo.Echo) {
	e.GET("/users", handlers.GetUserList)
	e.GET("/users/:id", handlers.GetUser)
	e.POST("/users", handlers.CreateUser)
	e.PUT("/users/:id", handlers.UpdateUser)
	e.DELETE("/users/:id", handlers.DeleteUser)

	e.POST("/users/:id/events", handlers.CreateEvent)

	e.GET("/statistics", handlers.GetStatisticsBetweenRange)
}
