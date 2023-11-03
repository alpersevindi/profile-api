package models

type Product struct {
	Name     string  `json:"name"`
	Price    float64 `json:"price"`
	Quantity int     `json:"quantity"`
}

type Event struct {
	ID        string  `json:"id"`
	Type      string  `json:"type"`
	Timestamp string  `json:"timestamp"`
	Product   Product `json:"product"`
}
