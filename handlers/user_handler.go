package handlers

import (
	"net/http"
	"profile-api/database"
	"profile-api/models"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func GetUserList(c echo.Context) error {
	input := &dynamodb.ScanInput{
		TableName: aws.String("Users"),
	}

	result, err := database.DynamoDBClient.Scan(input)
	if err != nil {
		return err
	}

	var users []models.User
	err = dynamodbattribute.UnmarshalListOfMaps(result.Items, &users)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, users)
}

func GetUser(c echo.Context) error {
	userID := c.Param("id")

	input := &dynamodb.GetItemInput{
		TableName: aws.String("Users"),
		Key: map[string]*dynamodb.AttributeValue{
			"ID": {
				S: aws.String(userID),
			},
		},
	}

	result, err := database.DynamoDBClient.GetItem(input)
	if err != nil {
		return err
	}

	var user models.User
	if err := dynamodbattribute.UnmarshalMap(result.Item, &user); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, user)
}

func CreateUser(c echo.Context) error {
	var user models.User
	if err := c.Bind(&user); err != nil {
		return err
	}

	user.ID = uuid.New().String()

	av, err := dynamodbattribute.MarshalMap(user)
	if err != nil {
		return err
	}

	input := &dynamodb.PutItemInput{
		TableName: aws.String("Users"),
		Item:      av,
	}

	_, err = database.DynamoDBClient.PutItem(input)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, user)
}

func UpdateUser(c echo.Context) error {
	userID := c.Param("id")
	var updatedUser models.User
	if err := c.Bind(&updatedUser); err != nil {
		return err
	}

	updatedUser.ID = userID

	av, err := dynamodbattribute.MarshalMap(updatedUser)
	if err != nil {
		return err
	}

	input := &dynamodb.PutItemInput{
		TableName: aws.String("Users"),
		Item:      av,
	}

	_, err = database.DynamoDBClient.PutItem(input)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, updatedUser)
}

func DeleteUser(c echo.Context) error {
	userID := c.Param("id")

	input := &dynamodb.DeleteItemInput{
		TableName: aws.String("Users"),
		Key: map[string]*dynamodb.AttributeValue{
			"ID": {
				S: aws.String(userID),
			},
		},
	}

	_, err := database.DynamoDBClient.DeleteItem(input)
	if err != nil {
		return err
	}

	return c.NoContent(http.StatusNoContent)
}
