package handlers

import (
	"context"
	"time"

	"net/http"
	"profile-api/database"
	"profile-api/models"

	"github.com/labstack/echo/v4"
)

func GetStatisticsBetweenRange(c echo.Context) error {
	var statistics models.Statistics
	if err := c.Bind(&statistics); err != nil {
		return err
	}

	conn, err := database.GetClickHouseConnection()
	if err != nil {
		return err
	}

	defer conn.Close()

	query := `
	SELECT
		SUM(price) AS total_price
	FROM
		user_events
	WHERE
		type = ?
		AND timestamp BETWEEN ? AND ?;
	`

	row := conn.QueryRow(
		context.Background(),
		query,
		statistics.Type,
		ConvertEpochTime(statistics.StartDate),
		ConvertEpochTime(statistics.EndDate),
	)
	var (
		totalValue float64
	)
	if err := row.Scan(&totalValue); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, totalValue)
}

func ConvertEpochTime(date int) string {
	epochTime := int64(date)
	timestamp := time.Unix(epochTime, 0)
	formattedTimestamp := timestamp.Format("2006-01-02T15:04:05-07:00")

	return formattedTimestamp
}
