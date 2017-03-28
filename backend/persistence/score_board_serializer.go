package persistence

import (
	"superstellar/backend/communication"
	"superstellar/backend/state"

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

func (serializer *ScoreBoardSerializer) serializeObjectDestroyed(objectDestroyed *events.ObjectDestroyed,
	client *communication.Client) *dynamodb.PutItemInput {

	spaceship := objectDestroyed.DestroyedObject.(*state.Spaceship)

	return &dynamodb.PutItemInput{
		TableName: aws.String("SuperstellarScoreBoard"),
		Item: map[string]*dynamodb.AttributeValue{
			"id":                {S: aws.String(uuid.NewV4().String())},
			"name":              {S: aws.String(client.UserName())},
			"spawn_time":        {S: aws.String(spaceship.SpawnTimestamp().String())},
			"death_time":        {S: aws.String(objectDestroyed.Timestamp.String())},
			"score":             {N: aws.String(fmt.Sprint(spaceship.MaxHP))},
			"hits":              {N: aws.String(fmt.Sprint(spaceship.Hits))},
			"hits_received":     {N: aws.String(fmt.Sprint(spaceship.HitsReceived))},
			"kills":             {N: aws.String(fmt.Sprint(spaceship.Kills))},
			"projectiles_fired": {N: aws.String(fmt.Sprint(spaceship.ProjectilesFired))},
		},
	}
}

func (serializer *ScoreBoardSerializer) writeScoreBoard(objectDestroyed *events.ObjectDestroyed, client *communication.Client) {
	if _, ok := objectDestroyed.DestroyedObject.(*state.Spaceship); ok {
		putItemInput := serializer.serializeObjectDestroyed(objectDestroyed, client)

		_, error := serializer.adapter.DynamoDb().PutItem(putItemInput)
		if error != nil {
			log.Println("Cannot put item to DynamoDB", error)
		}
	}
}
func (serializer *ScoreBoardSerializer) HandleObjectDestroyed(objectDestroyed *events.ObjectDestroyed) {
	if client, ok := serializer.server.GetClient(objectDestroyed.DestroyedObject.Id()); ok {
		go serializer.writeScoreBoard(objectDestroyed, client)
	}
}
