package models

type Event struct {
	ID        string  `json:"id" dynamodbav:"ID"`
	Type      string  `json:"type" dynamodbav:"Type"`
	Timestamp string  `json:"timestamp" dynamodbav:"Timestamp"`
	Product   Product `json:"product" dynamodbav:"Product"`
}

type Product struct {
	Name     string  `json:"name" dynamodbav:"Name"`
	Price    float64 `json:"price" dynamodbav:"Price"`
	Quantity int     `json:"quantity" dynamodbav:"Quantity"`
}
