package database

import (
	"os"

	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
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
		return
	}

	DynamoDBClient = dynamodb.New(sess)
}

func GetClickHouseConnection() (driver.Conn, error) {
	conn, err := clickhouse.Open(&clickhouse.Options{
		Addr: []string{"127.0.0.1:9000"},
	})

	if err != nil {
		panic(err)
	}

	return conn, nil
}
