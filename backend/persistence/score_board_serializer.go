package persistence

import (
	"superstellar/backend/communication"

	"superstellar/backend/events"

	"log"

	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/satori/go.uuid"
)

type ScoreBoardSerializer struct {
	server  *communication.Server
	adapter *DynamoDbAdapter
}

func NewScoreBoardSerializer(server *communication.Server, adapter *DynamoDbAdapter) *ScoreBoardSerializer {
	return &ScoreBoardSerializer{
		server:  server,
		adapter: adapter,
	}
}

func (serializer *ScoreBoardSerializer) serializeUserDied(userDied *events.UserDied,
	client *communication.Client) *dynamodb.PutItemInput {
	spaceship := userDied.ShotSpaceship

	return &dynamodb.PutItemInput{
		TableName: aws.String("ScoreBoard"),
		Item: map[string]*dynamodb.AttributeValue{
			"id":                {S: aws.String(uuid.NewV4().String())},
			"name":              {S: aws.String(client.UserName())},
			"spawn_time":        {S: aws.String(spaceship.SpawnTimestamp().String())},
			"death_time":        {S: aws.String(userDied.Timestamp.String())},
			"score":             {N: aws.String(fmt.Sprint(spaceship.MaxHP))},
			"hits":              {N: aws.String(fmt.Sprint(spaceship.Hits))},
			"hits_received":     {N: aws.String(fmt.Sprint(spaceship.HitsReceived))},
			"kills":             {N: aws.String(fmt.Sprint(spaceship.Kills))},
			"projectiles_fired": {N: aws.String(fmt.Sprint(spaceship.ProjectilesFired))},
		},
	}
}

func (serializer *ScoreBoardSerializer) writeScoreBoard(userDied *events.UserDied, client *communication.Client) {
	putItemInput := serializer.serializeUserDied(userDied, client)

	_, error := serializer.adapter.DynamoDb().PutItem(putItemInput)
	if error != nil {
		log.Println("Cannot put item to DynamoDB", error)
	}
}

func (serializer *ScoreBoardSerializer) HandleUserDied(userDied *events.UserDied) {
	client, ok := serializer.server.GetClient(userDied.ClientID)
	if ok {
		go serializer.writeScoreBoard(userDied, client)
	}
}
