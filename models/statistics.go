package models

type Statistics struct {
	StartDate int    `json:"start_date"`
	EndDate   int    `json:"end_date"`
	Type      string `json:"type"`
}
