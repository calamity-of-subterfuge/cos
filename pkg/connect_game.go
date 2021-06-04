package pkg

import (
	"fmt"
	"net/http"
	"time"

	"github.com/calamity-of-subterfuge/cos/pkg/utils"
	"github.com/gorilla/websocket"
)

// ConnectGame will connect to the game server at the given url, authenticating
// with the given JWT.
func ConnectGame(url string, jwt string) (*websocket.Conn, error) {
	headers := make(http.Header)
	headers.Add("Origin", utils.WEBSOCKET_ORIGIN)

	conn, _, err := websocket.DefaultDialer.Dial(url, headers)
	if err != nil {
		return nil, fmt.Errorf("dialing %s: %w", url, err)
	}

	err = conn.SetWriteDeadline(time.Now().Add(utils.CONN_WRITE_TIMEOUT))
	if err != nil {
		closeConn(websocket.CloseInternalServerErr, conn)
		return nil, fmt.Errorf("set write deadline: %w", err)
	}
	err = conn.WriteMessage(websocket.TextMessage, []byte(jwt))
	if err != nil {
		closeConn(websocket.ClosePolicyViolation, conn)
		return nil, fmt.Errorf("writing JWT: %w", err)
	}

	err = conn.SetWriteDeadline(time.Time{})
	if err != nil {
		closeConn(websocket.CloseInternalServerErr, conn)
		return nil, fmt.Errorf("clear write deadline: %w", err)
	}

	return conn, nil
}
