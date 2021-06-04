# utilsys

This package provides a completely self-contained utility system structure,
which is a very well-known method of generating complex, emergent AI behavior
without an excessive amount of computational resources. There is no requirement
that clients use this technique for generating behavior.

An introduction to utility systems can be found
[on wikipedia](https://en.wikipedia.org/wiki/Utility_system).

The naming conventions of this package are inspired by the rust crate
https://docs.rs/big-brain however it is not a port.

## Overview

- An `Action` is capable of manipulating the state.
- A `Scorer` assesses the environment and produces a measure of the utility of
  an `Action`.
- A `Thinker` is a type of `Action` which delegates to other `Actions` based on
  their utility.
- A `Qualifier` is a type of `Scorer` which delegates to other `Scorer`s to
  combine or rescale them.
- A `BayesScorer` is a special type of `Scorer` where the action has a utility
  of 1, but success is probabilistic and can be estimated using
  [Bayes' rule](https://www.youtube.com/watch?v=lG4VkPoG3ko).
  `BayesScorer` rely on a the standard odds (e.g. in general 1 succeed for every
  5 failures) and a set of `BayesFactor`s (e.g, 3x as likely to succeed when Y
  is researched) to produce the final odds (e.g. 3 successes for every 5
  failures) which are then converted to the final probability (e.g. 3 / (3+5) =
  0.375).

### Scoring Convention

Although not required, it's recommended that `Scorer`s follow a scoring
convention where all scores are between `0` and `1` (inclusive), where
the values are interpreted as:

- `<0.1`: Useless, valueless, or impossible. For example, the score for healing
  a full health entity.
- `0.1`: Possible, but doesn't do anything to achieve your goals.
- `0.1 - 1`: There is some value in doing this, and that value increases
  linearly with the score.
- `1`: It is not possible to get more worth or value from the action. For
  example, a failure to perform this action right now will immediately result in
  the game being lost.

## Example

### Creating an Action

Actions generally consist of two components - the exposed action, which acts
just as a stateless builder for actions, and the private action, which actually
does the thing and is stateful.

```go
type meanderAction struct{
    // the default value is utilsys.ActionStateInit, which means that
    // the next thing to happen to this action is Attached()
    state utilsys.ActionState
}

func (a *meanderAction) State() utilsys.ActionState {
    return a.state
}

func (a *meanderAction) Attached(world interface{}, actor interface{}) {
    // Typically you would do some kind of type assertion here
    // on world and actor then save them to the meanderAction
    // so they can be used on Execute()

    // When this is called state is ActionStateInit and usually this should
    // change the state to ActionStateRequested

    // This SHOULD NOT manipulate the world or actor. That should be done
    // in Execute.

    a.state = utilsys.ActionStateRequested
}

func (a *meanderAction) Execute(delta time.Duration) {
    // This is called if state is ActionStateRequested or ActionStateExecuting,
    // one per tick, and should eventually set the state to ActionStateSuccess
    // or ActionStateFailure. For the utilsys ActionStateRequested and
    // ActionStateExecuting are basically the same, but you can use them to
    // distinguish the first call to Execute with future calls to Execute

    if a.state == utilsys.ActionStateRequested {
        // For the purposes of demo lets pretend we need another tick
        a.state = utilsys.ActionStateExecuting
        return
    }

    a.state = utilsys.ActionStateSuccess
}

func (a *meanderAction) Cancel() {
    // This is called when the state is ActionStateRequested or ActionStateExecuting,
    // but for some reason we want the action to end as soon as possible. Invoking
    // this function should eventually result in ActionStateSuccess or ActionStateFailure.
    a.state = utilsys.ActionStateCanceled
}

func (a *meanderAction) FinishCanceling(delta time.Duration) {
    // FinishCanceling is called if state is ActionStateCanceled once per tick
    // and should eventually set the state to ActionStateSuccess or ActionStateFailure.
    a.state = utilsys.ActionStateFailed
}

func (a *menaderAction) Reset() {
  // This is called when the state is ActionStateSuccess or ActionStateFailure and
  // should result in ActionStateRequested, ActionStateExecuting, ActionStateSuccess,
  // or ActionStateFailure. This should NOT manipulate the world or actor.
  a.state = utilsys.ActionStateRequested
}

type MeanderAction struct{}
func (a MeanderAction) Build() utilsys.Action {
    return &meanderAction{}
}
```

### Creating a Scorer

Scorers follow essentially the same pattern as `Action`, but simpler.

```go
type foodScorer struct{
    // normally you'd use a stricter type for these here, since you
    // would have casted them in Attached()

    world interface{}
    actor interface{}
}

func (s *foodScorer) Attach(world interface{}, actor interface{}) {
    s.world = world
    s.actor = actor
}

func (s *foodScorer) Score() float64 {
    // some calculation here to get a number between 0 and 1. A good
    // default choice is that your Score() functions always use the
    // full range and then you use Qualifier's as necessary to rescale
    // or clip scores.
    return 0.5
}

type FoodScorer struct{}
func (s FoodScorer) Build() utilsys.Scorer {
    return &foodScorer{}
}
```

### Building the AI

Notice how there are few pointers as everything constructed at this step, with
the exception of `world`, is intended to be essentially stateless. This is not a
hard requirement but if it does not come naturally you may be interpreting the
interfaces incorrectly. Remember, `MeanderAction` is really an `ActionBuilder`
at this step.

This example uses `HighestScore` technique which is the simplest type of
thinker, which just runs whatever has the highest score at a given point,
breaking ties uniformly at random. This is what is classically meant by
a utility system.

Other types you should definitely consider, especially when nesting:

- A `FirstToScore` system is less theoretically pure but is more stable. It
  runs the first action in the list whose score meets or exceeds a threshold,
  falling back to `HighestScore`.
- A `LinearProbabilistic` system adds a lot of randomness to the AI. It selects
  an action from the given list of choices at random, where the odds of
  selecting an action is proportional to its score. Note this will often result
  in very low-score selections. One can specify a threshold for the minimum
  score to be included (so long as at least something reaches the minimum score)
  to avoid particularly ridiculous behavior.
- A `SoftMaxProbabilistic` system adds a bit of randomness to the AI. It selects
  an action from the given list of choices at random, where the odds of
  selecting an action is proportional to e^(factor*score). This is much more
  predictable than the `LinearProbabilistic` system for a reasonable factor value
  (usually between 5 and 30). Higher factors mean less random.

```go
var world interface{} // typically your pkg.Game
utilsys.NewAI(
    world,
    utilsys.NewHighestScoreThinker([]utilsys.ScoredActionBuilder{
        {
            Action: MeanderAction{},
            Scorer: utilsys.FixedScorer{Score: 0.1}
        },
        {
            Action: AcquireResourceAction{Resource: "gold"},
            Scorer: AcquireResourceScorer{}
        },
    })
)
```

### Using the AI

Using the AI just requires that you add all the actors to the AI via

```go
// actor is typically a Player or SmartObject, some actions could be either.
// Typically you would get this from (*client.State).OnSelfLoaded or
// (*client.State).OnControllableSmartObjectLoaded
var actor interface{}

ai.AddActor(actor)
```

Make sure you remember to detach any actors you no longer want to
control:

```go
// Typically this would happen from (*client.State).OnSelfLost or
// (*client.State).OnControllableSmartObjectLost
var actor interface{}

ai.RemoveActor(actor)
```

And then you regularly tick the AI:


```go
// Typically this would happen from your (pkg.Game).Tick
var delta time.Duration

ai.Tick(delta)
```

And that's all there is to it!
