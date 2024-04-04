package handlers

import (
	"context"
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

	conn, err := database.GetClickHouseConnection()
	if err != nil {
		return err
	}

	defer conn.Close()

	batch, err := conn.PrepareBatch(context.Background(), "INSERT INTO users (id, name, surname) VALUES (?, ?, ?)")
	if err != nil {
		return err
	}

	err = batch.Append(
		user.ID,
		user.Name,
		user.Surname,
	)
	if err != nil {
		return err
	}

	if err := batch.Send(); err != nil {
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

func GetUsersByAttribute(c echo.Context) error {
	attribute := c.QueryParam("attribute")

	exprAttrValues := map[string]*dynamodb.AttributeValue{
		":val": {
			S: aws.String(attribute),
		},
	}

	expression := "name = :val"

	input := &dynamodb.QueryInput{
		TableName:                 aws.String("Users"),
		KeyConditionExpression:    aws.String(expression),
		ExpressionAttributeValues: exprAttrValues,
	}

	result, err := database.DynamoDBClient.Query(input)
	if err != nil {
		return err
	}

	var users []models.User

	for _, item := range result.Items {
		var user models.User
		if err := dynamodbattribute.UnmarshalMap(item, &user); err != nil {
			return err
		}
		users = append(users, user)
	}

	return c.JSON(http.StatusOK, users)
}

func GetUsersByEvent(c echo.Context) error {
	eventValue := c.QueryParam("event")

	exprAttrValues := map[string]*dynamodb.AttributeValue{
		":val": {
			S: aws.String(eventValue),
		},
	}

	expression := "event = :val"

	input := &dynamodb.QueryInput{
		TableName:                 aws.String("Users"),
		IndexName:                 aws.String("EventIndex"),
		KeyConditionExpression:    aws.String(expression),
		ExpressionAttributeValues: exprAttrValues,
		ProjectionExpression:      aws.String("ID"),
	}

	result, err := database.DynamoDBClient.Query(input)
	if err != nil {
		return err
	}

	var userIDs []string

	for _, item := range result.Items {
		userID := aws.StringValue(item["ID"].S)
		userIDs = append(userIDs, userID)
	}

	return c.JSON(http.StatusOK, userIDs)
}
