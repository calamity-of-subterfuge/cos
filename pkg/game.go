package pkg

import (
	"time"

	"github.com/calamity-of-subterfuge/cos/pkg/srvpkts"
)

// Game describes something which actually plays the game, i.e., controls the
// AI. When implementing a Game it can be helpful to use a world.State to handle
// the generic client state, but it is not required.
type Game interface {
	// OnReceiveMessage should be called whenever the server sends
	// us a complete message which was successfully parsed.
	OnReceiveMessage(srvpkts.Packet)

	// OnDisconnected should be called if the websocket closed
	OnDisconnected()

	// Tick this game to account for the given amount of elapsed time,
	// called regularly
	Tick(time.Duration)
}
