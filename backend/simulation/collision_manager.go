package simulation

import (
	"container/list"
	"log"
	"math"
	"math/rand"
	"reflect"
	"superstellar/backend/events"
	"superstellar/backend/state"
	"superstellar/backend/types"
	"time"
)

type CollisionManager struct {
	space           *state.Space
	eventDispatcher *events.EventDispatcher
}

func NewCollisionManager(space *state.Space, eventDispatcher *events.EventDispatcher) *CollisionManager {
	return &CollisionManager{
		space:           space,
		eventDispatcher: eventDispatcher,
	}
}

func (manager *CollisionManager) resolveCollisions() {
	collided := make(map[state.Object]bool)
	oldVelocity := make(map[state.Object]*types.Vector)

	for _, object := range manager.space.Objects {

		collided[object] = true

		for _, otherObject := range manager.space.Objects {
			if !collided[otherObject] && object.DetectCollision(otherObject) {
				if _, exists := oldVelocity[object]; !exists {
					oldVelocity[object] = object.Velocity().Multiply(-1.0)
				}

				if _, exists := oldVelocity[otherObject]; !exists {
					oldVelocity[otherObject] = otherObject.Velocity().Multiply(-1.0)
				}

				manager.collide(object, otherObject)
			}
		}
	}

	queue := list.New()
	collidedThisTurn := make(map[state.Object]bool)
	visited := make(map[state.Object]bool)

	for object := range oldVelocity {
		queue.PushBack(object)
		collidedThisTurn[object] = true
		visited[object] = true
	}

	for e := queue.Front(); e != nil; e = e.Next() {
		object := e.Value.(state.Object)
		collidedThisTurn[object] = true
		object.SetPosition(object.Position().Add(oldVelocity[object]))

		for _, otherObject := range manager.space.Objects {
			if !collidedThisTurn[otherObject] && object.DetectCollision(otherObject) {
				oldVelocity[otherObject] = otherObject.Velocity().Multiply(-1.0)
				if !visited[otherObject] {
					visited[otherObject] = true
					queue.PushBack(otherObject)
				}

				manager.collide(object, otherObject)
			}
		}
	}

	// TODO kod przeciwzakrzepowy - wywalic jak zrobimy losowe spawnowanie
	collided2 := make(map[state.Object]bool)

	for _, object := range manager.space.Objects {
		collided2[object] = true
		for _, otherObject := range manager.space.Objects {
			if !collided2[otherObject] && object.DetectCollision(otherObject) {
				log.Printf("COLLISON")
				if val, exists := oldVelocity[object]; exists {
					log.Printf("ov1: %f %f", val.X, val.Y)
				}
				if val, exists := oldVelocity[otherObject]; exists {
					log.Printf("ov2: %f %f", val.X, val.Y)
				}
				log.Printf("v1: %f %f", object.Velocity().X, object.Velocity().Y)
				log.Printf("v2: %f %f", otherObject.Velocity().X, otherObject.Velocity().Y)
				log.Printf("p1: %d %d", object.Position().X, object.Position().Y)
				log.Printf("p2: %d %d", otherObject.Position().X, otherObject.Position().Y)

				randAngle := rand.Float64() * 2 * math.Pi
				randMove := types.NewVector(5000, 0).Rotate(randAngle)
				object.SetPosition(object.Position().Add(randMove))
			}
		}
	}
	// koniec kodu przeciwzakrzepowego
}

func (manager *CollisionManager) collide(objectA state.Object, objectB state.Object) {
	typeA := reflect.TypeOf(objectA)
	typeB := reflect.TypeOf(objectA)

	spaceshipType := reflect.TypeOf(&state.Spaceship{})
	simpleCollision := &SimpleCollision{}

	var collision Collision

	collision = simpleCollision

	if typeA == spaceshipType && typeB == spaceshipType {
		collision = simpleCollision
	} else {
		collision = simpleCollision
	}

	if collision != nil {
		collision.collide(objectA, objectB)
	}

	manager.checkHp(objectA, objectB)
	manager.checkHp(objectB, objectA)

}

func (manager *CollisionManager) checkHp(victim state.Object, predator state.Object) {
	if victim.Hp() <= 0 {
		manager.space.RemoveObject(victim.Id())

		if reflect.TypeOf(victim) == reflect.TypeOf(&state.Spaceship{}) {
			objectDestroyedMessage := &events.ObjectDestroyed{
				DestroyedObject: victim,
				DestroyedBy:     predator,
				//ShotSpaceship: victim.(*state.Spaceship),
				Timestamp: time.Now(),
			}
			manager.eventDispatcher.FireObjectDestroyed(objectDestroyedMessage)
		}
	}
}
