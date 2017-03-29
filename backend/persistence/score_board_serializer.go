package persistence

import (
	"fmt"
	"log"
	"superstellar/backend/events"
	"superstellar/backend/state"
	"superstellar/backend/utils"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

type ScoreBoardSerializer struct {
	userNameRegistry *utils.UserNamesRegistry
	adapter          *DynamoDbAdapter
	idManager        *utils.IdManager
	eventDispatcher  *events.EventDispatcher
}

func NewScoreBoardSerializer(userNameRegistry *utils.UserNamesRegistry, adapter *DynamoDbAdapter,
	idManager *utils.IdManager, eventDispatcher *events.EventDispatcher) *ScoreBoardSerializer {
	return &ScoreBoardSerializer{
		userNameRegistry: userNameRegistry,
		adapter:          adapter,
		idManager:        idManager,
		eventDispatcher:  eventDispatcher,
	}
}

func (serializer *ScoreBoardSerializer) serializeObjectDestroyed(objectDestroyed *events.ObjectDestroyed) *dynamodb.PutItemInput {
	spaceship := objectDestroyed.DestroyedObject.(*state.Spaceship)

	return &dynamodb.PutItemInput{
		TableName: aws.String("SuperstellarScoreBoard"),
		Item: map[string]*dynamodb.AttributeValue{
			"game":                {S: aws.String("superstellar")},
			"name":                {S: aws.String(serializer.userNameRegistry.GetUserName(spaceship.Id()))},
			"spawn_time":          {S: aws.String(spaceship.SpawnTimestamp().String())},
			"death_time":          {S: aws.String(objectDestroyed.Timestamp.String())},
			"score":               {N: aws.String(fmt.Sprint(spaceship.MaxHP))},
			"hits":                {N: aws.String(fmt.Sprint(spaceship.Hits))},
			"hits_received":       {N: aws.String(fmt.Sprint(spaceship.HitsReceived))},
			"kills":               {N: aws.String(fmt.Sprint(spaceship.Kills))},
			"destroyed_asteroids": {N: aws.String(fmt.Sprint(spaceship.DestroyedAsteroids))},
			"projectiles_fired":   {N: aws.String(fmt.Sprint(spaceship.ProjectilesFired))},
		},
	}
}

func (serializer *ScoreBoardSerializer) writeScoreBoard(objectDestroyed *events.ObjectDestroyed) {
	if spaceship, ok := objectDestroyed.DestroyedObject.(*state.Spaceship); ok {
		putItemInput := serializer.serializeObjectDestroyed(objectDestroyed)

		_, error := serializer.adapter.DynamoDb().PutItem(putItemInput)
		if error != nil {
			log.Println("Cannot put item to DynamoDB", error)
		} else {
			serializer.fireScoreSentEvent(spaceship.MaxHP)
			log.Println("ScoreBoard item sent to DynamoDB")
		}
	}
}

func (serializer *ScoreBoardSerializer) HandleObjectDestroyed(objectDestroyed *events.ObjectDestroyed) {
	if serializer.idManager.IsPlayerId(objectDestroyed.DestroyedObject.Id()) {
		go serializer.writeScoreBoard(objectDestroyed)
	}
}

func (serializer *ScoreBoardSerializer) fireScoreSentEvent(score uint32) {
	scoreSent := &events.ScoreSent{
		Score: score,
	}

	serializer.eventDispatcher.FireScoreSent(scoreSent)
}
