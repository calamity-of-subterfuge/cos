package pkg

import (
	"crypto/rand"
	"encoding/base64"
	"log"
	"time"

	"github.com/calamity-of-subterfuge/cos/pkg/utils"
	"github.com/gorilla/websocket"
)

func closeConn(code int, conn *websocket.Conn) {
	err := conn.SetWriteDeadline(time.Now().Add(utils.CONN_WRITE_TIMEOUT))
	if err != nil {
		conn.Close()
		return
	}

	conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(code, ""))
	conn.Close()
}

func generateSecureToken(length int) string {
	b := make([]byte, length)
	if _, err := rand.Read(b); err != nil {
		log.Fatalf("Failed to generate secure token")
	}
	return base64.URLEncoding.EncodeToString(b)
}
