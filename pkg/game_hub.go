package pkg

import (
	"log"
	"time"

	"github.com/calamity-of-subterfuge/cos/pkg/srvpkts"
	"github.com/gorilla/websocket"
)

// GameHub manages a single server websocket connection in order to run a
// single game, notifying a particular channel upon completion and allowing
// cancellation
type GameHub struct {
	UID string

	game              Game
	conn              *Conn
	recvQueue         chan ReceivedMessage
	connClosed        chan string
	finishNotifyQueue chan string
	cancelChan        chan struct{}
}

// NewGameHub takes over management of the given game server websocket to
// run a game initialized using the given GameConstructor. This does not
// start managing the game; that should be done in a dedicated goroutine
// by calling the long-running function Manage()
func NewGameHub(conn *websocket.Conn, uid string, finishNotifyQueue chan string, gameConstructor GameConstructor) *GameHub {
	recvQueue := make(chan ReceivedMessage, 1024)
	closedQueue := make(chan string, 1)
	wrappedConn := NewConn(conn, uid, recvQueue, closedQueue)
	game := gameConstructor(wrappedConn.SendQueue)

	return &GameHub{
		UID:               uid,
		conn:              wrappedConn,
		game:              game,
		recvQueue:         recvQueue,
		connClosed:        closedQueue,
		finishNotifyQueue: finishNotifyQueue,
		cancelChan:        make(chan struct{}, 1),
	}
}

// Closes this game hub if it is being managed right now. GameHubs cannot be
// reused.
func (h *GameHub) Close() {
	select {
	case h.cancelChan <- struct{}{}:
	default:
	}
}

// Manage this game hub. Typically run on a dedicated goroutine, this will
// monitor the channels for this game hub in order to execute the game.
func (h *GameHub) Manage() {
	ticker := time.NewTicker(time.Second / 60)
	lastTick := time.Now()
	lastWarnBehindAt := time.Now()

manageLoop:
	for {
		select {
		case msg := <-h.recvQueue:
			srvPacket, err := srvpkts.ParseSinglePacket(msg.Message)
			if err != nil {
				log.Printf("ignoring bad packet from server: %v (%v)", msg.Message, err)
				break
			}

			h.game.OnReceiveMessage(srvPacket)
		case curTick := <-ticker.C:
			if time.Since(curTick) > time.Second/5 {
				// Eat the tick to avoid falling so far behind; we'll tick again
				// later with a big time.Duration
				if time.Since(lastWarnBehindAt) > 5*time.Minute {
					log.Printf(
						"WARN: eating ticks because they are too old (%v) - "+
							"this warning happens only once per 5 minutes "+
							"and means that the AI is overloaded. This results "+
							"in ticks with large Duration's, which can lead to"+
							"instability",
						time.Since(curTick),
					)
					lastWarnBehindAt = time.Now()
				}
				break
			}

			h.game.Tick(curTick.Sub(lastTick))
			lastTick = curTick
		case <-h.connClosed:
			break manageLoop
		case <-h.cancelChan:
			break manageLoop
		}
	}

	log.Printf("GameHub %s shutting down", h.UID)
	h.game.OnDisconnected()
	h.Close()
	h.conn.Close()
	ticker.Stop()
	h.finishNotifyQueue <- h.UID
}
