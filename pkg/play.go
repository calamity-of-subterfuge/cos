package pkg

import (
	"errors"
	"log"
	"math"
	"time"

	"github.com/gorilla/websocket"
)

// Config is the configuration used for the standard Play loop.
type Config struct {
	// Email of the calamity of subterfuge account to login with
	Email string

	// GrantIden identifies the non-human authentication method used for
	// the acount
	GrantIden string

	// Secret is the secret that allows the use of the grant
	Secret string

	// AIConfig is the configuration for the AI
	AIConfig *AIConfig
}

func Play(cfg *Config, gameConstructor GameConstructor) {
	for {
		log.Println("Logging in...")
		var auth *AuthToken
		var err error
		retryCounter := 0
		for {
			auth, err = Login(cfg.Email, cfg.GrantIden, cfg.Secret)
			if err == nil {
				break
			}

			sleepSeconds := int64(math.Pow(2, float64(retryCounter))) * 60
			log.Printf("Error logging in, retrying in %d seconds: %v", sleepSeconds, err)
			time.Sleep(time.Duration(sleepSeconds) * time.Second)
			if retryCounter < 4 {
				retryCounter++
			}
		}
		log.Println("Successfully logged in; connecting to lobby socket server...")

		var socketConn *websocket.Conn
		var welcomeMessage map[string]interface{}
		retryCounter = 0
		socketConn, welcomeMessage, err = QueueAI(cfg.AIConfig, auth)
		for err != nil && retryCounter < 5 {
			sleepSeconds := int64(math.Pow(2, float64(retryCounter))) * 60
			log.Printf("Failed to queue AI, retrying in %d seconds: %v", sleepSeconds, err)
			time.Sleep(time.Duration(sleepSeconds) * time.Second)
			retryCounter = retryCounter + 1
			socketConn, welcomeMessage, err = QueueAI(cfg.AIConfig, auth)
		}
		if retryCounter == 5 {
			log.Println("Too many failures to queue ai in a row; relogging in")
			continue
		}

		hub := NewHub(socketConn, welcomeMessage, gameConstructor)
		err = hub.Manage()
		if err != nil {
			if errors.Is(err, ErrCanceled) {
				break
			} else {
				log.Printf("Error while managing the hub, relogging in in 5 seconds: %v", err)
			}
		}
		time.Sleep(5 * time.Second)
	}
}
