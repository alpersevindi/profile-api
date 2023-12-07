package models

type User struct {
	ID      string  `json:"id" dynamodbav:"ID"`
	Name    string  `json:"name" dynamodbav:"Name"`
	Surname string  `json:"surname" dynamodbav:"Surname"`
	Events  []Event `json:"events" dynamodbav:"Events,omitempty"`
}
