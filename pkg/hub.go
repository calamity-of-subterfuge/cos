package pkg

import (
	"log"

	"github.com/gorilla/websocket"
)

// HubError is an enum of errors that can cause the hub to stop managing.
// Unexpected errors will come with no HubError attached.
type HubError int

const (
	ErrConnectionGoingAway HubError = 1
	ErrCanceled            HubError = 2
)

// Error implements the error interface for HubError
func (e HubError) Error() string {
	switch e {
	case ErrConnectionGoingAway:
		return "connection going away"
	case ErrCanceled:
		return "canceled"
	default:
		return "unknown error"
	}
}

// GameConstructor describes something which can initialize games from
// a sendQueue, where the sendQueue is a channel that client packets can
// be sent to be forwarded to the server.
type GameConstructor func(sendQueue chan interface{}) Game

// Hub manages a lobby socket connection in order to detect and handle
// game notifications by connecting to the server and then initializing
// a game with a given GameConstructor, then managing the connection
// to use that Game.
type Hub struct {
	lobbySocketConn       *Conn
	lobbySocketRecvQueue  chan ReceivedMessage
	lobbySocketClosedChan chan string

	gameConstructor   GameConstructor
	gameHubsByUID     map[string]*GameHub
	gameFinishedQueue chan string
	cancelChan        chan struct{}
	welcomeMsg        map[string]interface{}
}

// NewHub initializes a hub that will take over the given lobby socket
// connection. Upon a receiving a notification about a new game this
// will handle connecting to the server and managing the websocket in
// order to run the Game produced by the given GameConstructor. The
// welcomeMsg should be the first packet received on the lobby connection
// and is used exclusively for debugging.
func NewHub(lobbyConn *websocket.Conn, welcomeMsg map[string]interface{}, gameConstructor GameConstructor) *Hub {
	recvQueue := make(chan ReceivedMessage, 64)
	closedChan := make(chan string, 1)
	return &Hub{
		lobbySocketConn:       NewConn(lobbyConn, "ls", recvQueue, closedChan),
		lobbySocketRecvQueue:  recvQueue,
		lobbySocketClosedChan: closedChan,

		gameConstructor:   gameConstructor,
		gameHubsByUID:     make(map[string]*GameHub),
		gameFinishedQueue: make(chan string, 16),
		cancelChan:        make(chan struct{}, 1),
		welcomeMsg:        welcomeMsg,
	}
}

// Manage the hub forever or until we are disconnected from the lobby or
// Cancel'd. This cannot be run in multiple routines simultaneously.
func (h *Hub) Manage() error {
	var manageEndReason error

	log.Printf("Hub manage loop started..")

manageLoop:
	for {
		select {
		case msg := <-h.lobbySocketRecvQueue:
			log.Printf("Notification from lobby-socket server: %v", msg.Message)
			typeRaw, found := msg.Message["type"]
			if !found {
				log.Printf("Ignoring notification (missing type)")
				break
			}

			typeStr, ok := typeRaw.(string)
			if !ok {
				log.Printf("Ignoring notification (type not a string)")
				break
			}

			switch typeStr {
			case "match-available":
				manageEndReason = h.handleMatchAvailable(msg)
				if manageEndReason != nil {
					break manageLoop
				}
			default:
				log.Printf("Ignoring notification (unknown type: %s)", typeStr)
			}
		case gameUID := <-h.gameFinishedQueue:
			log.Printf("Game finished: %s", gameUID)
			delete(h.gameHubsByUID, gameUID)
		case <-h.lobbySocketClosedChan:
			manageEndReason = ErrConnectionGoingAway
			break manageLoop
		case <-h.cancelChan:
			manageEndReason = ErrCanceled
			break manageLoop
		}
	}

	for _, gh := range h.gameHubsByUID {
		gh.Close()
	}

	// this avoids the game finished channel filling up and gives a chance
	// for the game hubs to actually finish
	for len(h.gameHubsByUID) > 0 {
		gameUID := <-h.gameFinishedQueue
		log.Printf("Game finished: %s", gameUID)
		delete(h.gameHubsByUID, gameUID)
	}

	return manageEndReason
}

func (h *Hub) handleMatchAvailable(msg ReceivedMessage) error {
	urlRaw, found := msg.Message["url"]
	if !found {
		log.Printf("Ignoring notification (missing URL)")
		return nil
	}

	url, ok := urlRaw.(string)
	if !ok {
		log.Printf("Ignoring notification (url not a string)")
		return nil
	}

	var jwtRaw interface{}
	jwtRaw, found = msg.Message["jwt"]
	if !found {
		log.Printf("Ignoring notification (missing jwt)")
		return nil
	}
	var jwt string
	jwt, ok = jwtRaw.(string)
	if !ok {
		log.Printf("Ignoring notification (jwt not a string)")
		return nil
	}

	gconn, err := ConnectGame(url, jwt)
	if err != nil {
		log.Printf("Failed to connect to %s: %v", url, err)
		return nil
	}

	uid := generateSecureToken(23)
	gh := NewGameHub(gconn, uid, h.gameFinishedQueue, h.gameConstructor)
	h.gameHubsByUID[uid] = gh

	go gh.Manage()

	log.Printf("Assigned match the uid %s", uid)
	return nil
}

// Cancels the hub. This may be called from any goroutine and, if Manage
// is being run, it will stop and return ErrCanceled
func (h *Hub) Cancel() {
	select {
	case h.cancelChan <- struct{}{}:
	default:
	}
}
