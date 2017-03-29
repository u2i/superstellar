package persistence

import (
	"log"

	"superstellar/backend/pb"

	"strconv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

type ScoreBoardReader struct {
	adapter *DynamoDbAdapter
}

func NewScoreBoardReader(adapter *DynamoDbAdapter) *ScoreBoardReader {
	return &ScoreBoardReader{
		adapter: adapter,
	}
}

func (reader *ScoreBoardReader) ReadScoreBoard() *pb.ScoreBoard {
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

	resp, error := reader.adapter.dynamodb.Query(queryInput)
	if error != nil {
		log.Println("Cannot get items from DynamoDB", error)
		return nil
	}

	protoScoreBoardItems := make([]*pb.ScoreBoardItem, 0, *resp.Count)
	protoScoreBoard := &pb.ScoreBoard{Items: protoScoreBoardItems}

	for i := range resp.Items {
		item := resp.Items[i]

		score, err := strconv.Atoi(*item["score"].N)
		if err == nil {
			scoreBoardItem := &pb.ScoreBoardItem{Name: *item["name"].S, Score: uint32(score)}
			protoScoreBoard.Items = append(protoScoreBoard.Items, scoreBoardItem)
		}
	}

	return protoScoreBoard
}
