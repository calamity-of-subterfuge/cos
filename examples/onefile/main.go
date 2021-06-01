package onefile

import (
	"log"
	"math/rand"
	"time"

	cos "github.com/calamity-of-subterfuge/cos/pkg"
	"github.com/calamity-of-subterfuge/cos/pkg/client"
	"github.com/calamity-of-subterfuge/cos/pkg/clipkts"
	"github.com/calamity-of-subterfuge/cos/pkg/srvpkts"
	"github.com/calamity-of-subterfuge/cos/pkg/utils"
)

// Game typically goes in a different file and is your implementation of
// cos.Game
type Game struct {
	// sendQueue is how you send messages to the server
	sendQueue chan interface{}

	// state is the state of the world for you as a client. you will
	// need either this or your own version of it.
	state *client.State

	// chat maintains the chat history, local chat authors, and recent chat
	// authors. you'll probably want this, but you might be able get away with
	// just watching for ChatMessagePacket's
	chat *client.Chat

	// timeUntilNextHello is just a dummy value to show how to send
	// packets; see Tick
	timeUntilNextHello time.Duration
}

// NewGame initializes a new Game which can send packets using the
// given sendQueue
func NewGame(sendQueue chan interface{}) cos.Game {
	return &Game{
		sendQueue:          sendQueue,
		state:              client.NewState(),
		chat:               client.NewChat(100),
		timeUntilNextHello: time.Second * 5,
	}
}

func (g *Game) OnReceiveMessage(packet srvpkts.Packet) {
	g.state.HandleMessage(packet)
	g.chat.HandleMessage(packet)

	switch v := packet.(type) {
	case *srvpkts.ChatMessagePacket:
		log.Printf("message from %v (%v): %s", g.chat.LocalChatAuthorsByUID[v.AuthorUID].Name, v.AuthorUID, v.Text)
	}
}

func (g *Game) OnDisconnected() {}
func (g *Game) Tick(delta time.Duration) {
	if g.state.GameTime == 0 {
		// no game sync yet
		return
	}

	g.timeUntilNextHello -= delta

	if g.timeUntilNextHello <= 0 {
		g.timeUntilNextHello = time.Second * 5
		// notice how Type does not need to be filled in
		g.sendQueue <- clipkts.SendLocalMessagePacket{
			Text: "hello world!",
		}
	}
}

// RunExample actually executes the example; typically this would be
// just called main() and the package would be main
func RunExample() {
	// AI personalities should use some amount of randomness. Do not forget
	// to seed rand!
	rand.Seed(time.Now().UnixNano())

	// you should load these from somewhere! typically a configuration file,
	// via environment variables, or via flags (https://golang.org/pkg/flag/)
	var email string = "example@example.com"
	var grantIden string = "pa_xyz"
	var secret = "mysecret"

	// these are fine to keep in the code but depend on your infrastructure
	// so it can be convenient to load from somewhere.
	var maxConcurrentInstances int = 2
	var clientAllowList []string = make([]string, 0)

	// these should be checked in with your repo
	var aiName string = "ExampleAI"
	var aiUID string = "example-ai"
	var version string = "0.0.1"
	var role utils.Role = utils.RoleEconomyAI

	aiCfg := &cos.AIConfig{
		AIName:                 aiName,
		AIUID:                  aiUID,
		Version:                version,
		Role:                   role,
		ClientAllowList:        clientAllowList,
		MaxConcurrentInstances: maxConcurrentInstances,
	}

	cos.Play(
		&cos.Config{
			Email:     email,
			GrantIden: grantIden,
			Secret:    secret,
			AIConfig:  aiCfg,
		},
		NewGame,
	)
}
