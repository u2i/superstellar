package simulation

import (
	"container/list"
	"reflect"
	"superstellar/backend/state"
	"superstellar/backend/types"
)

type CollisionManager struct {
	space *state.Space
}

func NewCollisionManager(space *state.Space) *CollisionManager {
	return &CollisionManager{
		space: space,
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
}
