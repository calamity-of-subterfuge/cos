package utilsys

import "time"

type ActionState int

const (
	// The initial state of an Action; implies it needs to be Attach'd.
	ActionStateInit ActionState = 0

	// The action has never been Execute'd before and should have Execute
	// called
	ActionStateRequested ActionState = 1

	// The action has had Execute called before but still needs Execute
	// to be called
	ActionStateExecuting ActionState = 2

	// Something has requested the action be canceled, but it is not done
	// canceling so FinishCanceling should be called
	ActionStateCanceled ActionState = 3

	// The action succeeded
	ActionStateSuccess ActionState = 4

	// The action failed
	ActionStateFailure ActionState = 5
)

// Action is a stateful object that acts upon a given actor in a given world.
// These are produced by action builders, which are what go into the AI.
type Action interface {
	// State returns the state of this action. Actions are responsible for
	// ensuring their state goes through the correct phases.
	State() ActionState

	// Attached is called when the action is in the state Init to let them know
	// which world and actor they are acting upon. The Action should typically
	// verify these are of the appropriate type and store them in a stricter
	// type, then move to Requested, Success, or Failure as appropriate. A
	// single Action is not reused across actors.
	Attached(world, actor interface{})

	// Execute this action within the world on the actor, after the given
	// amount of elapsed time since the last call (or an arbitrary value
	// if never called before). Should eventually result in the action
	// transitioning to Success or Failure.
	Execute(delta time.Duration)

	// Cancel this action, which tells the Action to do whatever is necessary
	// to get to the Success or Failure state as quickly as possible. Typically
	// this will either imemdiately update the state of the Action to Success
	// or Failure, or move the action into the Canceled state. This is NOT called
	// if the actor is removed from the AI - the Action simply will no longer
	// receive callbacks in that event, to avoid tedious nil handling within
	// each Action.
	Cancel()

	// FinishCanceling is called when the Action is in the Cancel state, and
	// should eventually move the action to the Success or Failure state. It
	// is passed the elapsed time since the last call or an arbitrary value
	// if never called before.
	FinishCanceling(delta time.Duration)

	// Reset is called only in the Success or Failure state, and acts as the
	// equivalent of Attached except for an instance that's already been used
	// before.
	Reset()
}

// ActionBuilder is something capable of building unattached actions and is
// generally stateless
type ActionBuilder interface {
	// Build the action
	Build() Action
}
