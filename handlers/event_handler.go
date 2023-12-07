package handlers

import (
	"net/http"
	"profile-api/database"
	"profile-api/models"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func CreateEvent(c echo.Context) error {
	userID := c.Param("id")
	var event models.Event
	if err := c.Bind(&event); err != nil {
		return err
	}

	event.ID = uuid.New().String()
	event.Timestamp = time.Now().Format(time.RFC3339)

	av, err := dynamodbattribute.MarshalMap(event)
	if err != nil {
		return err
	}

	input := &dynamodb.UpdateItemInput{
		TableName: aws.String("Users"),
		Key: map[string]*dynamodb.AttributeValue{
			"ID": {
				S: aws.String(userID),
			},
		},
		UpdateExpression: aws.String("SET #events = list_append(if_not_exists(#events, :empty_list), :event)"),
		ExpressionAttributeNames: map[string]*string{
			"#events": aws.String("Events"),
		},
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":event": {
				L: []*dynamodb.AttributeValue{
					{
						M: av,
					},
				},
			},
			":empty_list": {
				L: []*dynamodb.AttributeValue{},
			},
		},
		ReturnValues: aws.String("ALL_NEW"),
	}

	result, err := database.DynamoDBClient.UpdateItem(input)
	if err != nil {
		return err
	}

	var updatedUser models.User
	err = dynamodbattribute.UnmarshalMap(result.Attributes, &updatedUser)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, updatedUser.Events[len(updatedUser.Events)-1])
}
