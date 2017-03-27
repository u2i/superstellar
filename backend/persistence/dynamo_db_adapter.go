package persistence

import (
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

type DynamoDbAdapter struct {
	config   *aws.Config
	dynamodb *dynamodb.DynamoDB
}

func NewDynamoDbWriter() *DynamoDbAdapter {
	region := os.Getenv("DYNAMODB_REGION")
	endpoint := os.Getenv("DYNAMODB_ENDPOINT")

	config := &aws.Config{
		Region:   &region,
		Endpoint: &endpoint,
	}

	return &DynamoDbAdapter{
		config:   config,
		dynamodb: dynamodb.New(session.Must(session.NewSession()), config),
	}
}

func (adapter *DynamoDbAdapter) DynamoDb() *dynamodb.DynamoDB {
	return adapter.dynamodb
}
