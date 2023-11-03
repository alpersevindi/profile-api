package handlers

import (
	"net/http"
	"profile-api/models"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

var users []models.User

func GetUserList(c echo.Context) error {
	return c.JSON(http.StatusOK, users)
}

func GetUser(c echo.Context) error {
	userID := c.Param("id")
	for _, user := range users {
		if user.ID == userID {
			return c.JSON(http.StatusOK, user)
		}
	}
	return c.String(http.StatusNotFound, "User not found")
}

func CreateUser(c echo.Context) error {
	var user models.User
	if err := c.Bind(&user); err != nil {
		return err
	}
	user.ID = uuid.New().String()
	users = append(users, user)
	return c.JSON(http.StatusCreated, user)
}

func UpdateUser(c echo.Context) error {
	userID := c.Param("id")
	var updatedUser models.User
	if err := c.Bind(&updatedUser); err != nil {
		return err
	}
	for i, user := range users {
		if user.ID == userID {
			updatedUser.ID = user.ID
			users[i] = updatedUser
			return c.JSON(http.StatusOK, updatedUser)
		}
	}
	return c.String(http.StatusNotFound, "User not found")
}

func DeleteUser(c echo.Context) error {
	userID := c.Param("id")
	for i, user := range users {
		if user.ID == userID {
			users = append(users[:i], users[i+1:]...)
			return c.NoContent(http.StatusNoContent)
		}
	}
	return c.String(http.StatusNotFound, "User not found")
}

func CreateEvent(c echo.Context) error {
	userID := c.Param("id")
	var event models.Event
	if err := c.Bind(&event); err != nil {
		return err
	}
	event.ID = uuid.New().String()
	event.Timestamp = time.Now().Format(time.RFC3339)
	for i, user := range users {
		if user.ID == userID {
			user.Events = append(user.Events, event)
			users[i] = user
			return c.JSON(http.StatusCreated, event)
		}
	}
	return c.String(http.StatusNotFound, "User not found")
}

func UpdateEvent(c echo.Context) error {
	userID := c.Param("id")
	eventID := c.Param("eventID")
	var updatedEvent models.Event
	if err := c.Bind(&updatedEvent); err != nil {
		return err
	}
	for i, user := range users {
		if user.ID == userID {
			for j, event := range user.Events {
				if event.ID == eventID {
					updatedEvent.ID = event.ID
					updatedEvent.Timestamp = time.Now().Format(time.RFC3339)
					user.Events[j] = updatedEvent
					users[i] = user
					return c.JSON(http.StatusOK, updatedEvent)
				}
			}
		}
	}
	return c.String(http.StatusNotFound, "User or event not found")
}

func DeleteEvent(c echo.Context) error {
	userID := c.Param("id")
	eventID := c.Param("eventID")
	for i, user := range users {
		if user.ID == userID {
			for j, event := range user.Events {
				if event.ID == eventID {
					user.Events = append(user.Events[:j], user.Events[j+1:]...)
					users[i] = user
					return c.NoContent(http.StatusNoContent)
				}
			}
		}
	}
	return c.String(http.StatusNotFound, "User or event not found")
}
