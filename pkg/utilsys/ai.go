package utilsys

import (
	"log"
	"time"
)

type actionActor struct {
	actor  interface{}
	action Action
}

// AI runs Actions on all the actors within the world.
type AI struct {
	world      interface{}
	coreAction ActionBuilder

	actors []actionActor
}

// NewAI constructs a new AI within the given world, which uses the given
// coreAction for all actors. Typically coreAction is a Thinker, though
// this is not enforced.
func NewAI(world interface{}, coreAction ActionBuilder) *AI {
	return &AI{
		world:      world,
		coreAction: coreAction,
		actors:     make([]actionActor, 0),
	}
}

// AddActor adds the given actor to be handled by this AI.
//
// performance: O(1) amortized
func (ai *AI) AddActor(actor interface{}) {
	action := ai.coreAction.Build()
	action.Attached(ai.world, actor)

	ai.actors = append(ai.actors, actionActor{
		actor:  actor,
		action: action,
	})
}

// RemoveActor removes the given actor from being handled by this AI.
//
// performance: O(n) where n is the number of actors
func (ai *AI) RemoveActor(actor interface{}) {
	for idx := 0; idx < len(ai.actors); idx++ {
		if ai.actors[idx].actor == actor {
			ai.actors = append(ai.actors[:idx], ai.actors[idx+1:]...)
			return
		}
	}
}

// Tick all of the actions for actors handled by this AI, informing
// them the given amount of time has passed.
func (ai *AI) Tick(delta time.Duration) {
	for _, actorAction := range ai.actors {

		// We cut the loop off after 2 times to avoid an infinite loop,
		// but the idea is to allow for a full cycle of
		// 1. Success state cause Reset -> Requested
		// 3. Requested state causes Success
	actionUpdateLoop:
		for i := 0; i < 2; i++ {
			switch actorAction.action.State() {
			case ActionStateRequested:
				fallthrough
			case ActionStateExecuting:
				actorAction.action.Execute(delta)
				break actionUpdateLoop
			case ActionStateCanceled:
				actorAction.action.FinishCanceling(delta)
				break actionUpdateLoop
			case ActionStateSuccess:
				fallthrough
			case ActionStateFailure:
				actorAction.action.Reset()
			default:
				log.Panicf("action has bad State(): %v", actorAction.action.State())
			}
		}
	}
}
