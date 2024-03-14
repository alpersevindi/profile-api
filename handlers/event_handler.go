package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"profile-api/database"
	"profile-api/models"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/confluentinc/confluent-kafka-go/kafka"
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

	p, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": "localhost:9092",
		"client.id":         "foo",
		"acks":              "all",
	})
	if err != nil {
		fmt.Printf("Failed to create producer: %s\n", err)
	}
	topic := "ES"

	go func() {
		consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
			"bootstrap.servers": "localhost:9092",
			"group.id":          "foo",
			"auto.offset.reset": "smallest",
		})
		if err != nil {
			log.Fatal(err)
		}
		err = consumer.Subscribe(topic, nil)
		if err != nil {
			log.Fatal(err)
		}

		for {
			ev := consumer.Poll(100)
			switch e := ev.(type) {
			case *kafka.Message:
				fmt.Printf("consumed message: %s\n", string(e.Value))
			case *kafka.Error:
				fmt.Printf("%s\n", e)
			}
		}
	}()

	msg, err := json.Marshal(map[string]interface{}{
		"Id":       event.ID,
		"Type":     event.Type,
		"Time":     event.Timestamp,
		"Name":     event.Product.Name,
		"Price":    event.Product.Price,
		"Quantity": event.Product.Quantity,
	})

	dch := make(chan kafka.Event, 10000)
	err = p.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          msg,
	},
		dch,
	)
	if err != nil {
		log.Fatal(err)
	}
	<-dch

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("eu-west-1"),
	})
	if err != nil {
		log.Fatal("Failed to create AWS session:", err)
	}

	s3Client := s3.New(sess)

	s3result, err := s3Client.PutObject(&s3.PutObjectInput{
		Bucket: aws.String("profile-api-event-bucket"),
		Key:    aws.String(""),
		Body:   bytes.NewReader(msg),
	})
	if err != nil {
		log.Println("Error uploading event to S3:", err)
	} else {
		log.Printf("file uploaded to, %s\n", s3result)

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

	conn, err := database.GetClickHouseConnection()
	if err != nil {
		return err
	}

	defer conn.Close()

	batch, err := conn.PrepareBatch(context.Background(), "INSERT INTO user_events (id, type, user_id, price, timestamp) VALUES (?, ?, ?, ?, ?)")
	if err != nil {
		return err
	}

	err = batch.Append(
		event.ID,
		event.Type,
		userID,
		event.Product.Price,
		event.Timestamp,
	)
	if err != nil {
		return err
	}

	if err := batch.Send(); err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, updatedUser.Events[len(updatedUser.Events)-1])
}
