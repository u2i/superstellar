package persistence

import (
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

type ScoreBoardReader struct {
	adapter *DynamoDbAdapter
}

func NewScoreBoardReader(adapter *DynamoDbAdapter) *ScoreBoardSerializer {
	return &ScoreBoardSerializer{
		adapter: adapter,
	}
}

func (serializer *ScoreBoardSerializer) ReadScoreBoard() {
	queryInput := &dynamodb.QueryInput{
		TableName:              aws.String("SuperstellarScoreBoard"),
		IndexName:              aws.String("game-score-index"),
		KeyConditionExpression: aws.String("game = :game AND score > :min_score"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":game":      {S: aws.String("superstellar")},
			":min_score": {N: aws.String("0")},
		},
		ScanIndexForward: aws.Bool(false),
		Limit:            aws.Int64(10),
	}

	resp, error := serializer.adapter.dynamodb.Query(queryInput)
	if error != nil {
		log.Println("Cannot get items from DynamoDB", error)
	} else {
		log.Println(resp)
	}
}
