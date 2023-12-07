package database

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

var DynamoDBClient *dynamodb.DynamoDB

func InitDynamoDB() {
	os.Setenv("AWS_PROFILE", "go_project")

	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String("eu-west-1"),
		Credentials: credentials.NewSharedCredentials("", "go_project"),
	})

	if err != nil {
		fmt.Println(err)
		return
	}

	DynamoDBClient = dynamodb.New(sess)
}
