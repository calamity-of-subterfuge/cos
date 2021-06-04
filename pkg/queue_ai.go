package pkg

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/calamity-of-subterfuge/cos/pkg/utils"
	"github.com/gorilla/websocket"
	"golang.org/x/net/http2"
)

// AIConfig describes an AI that you are queueing and any additional
// parameters for the server.
type AIConfig struct {
	// AIName is the name of our AI personality
	AIName string

	// AIUID is the unique identifier we assigned to our AI personality,
	// which is typically 23 random bytes encoded in some url safe format
	AIUID string

	// Version is a valid semantic identifier for our AI personality
	Version string

	// Role is the role this AI plays - economy, military, or science
	Role utils.Role

	// ClientAllowList are the UIDs of the users which are allowed to
	// select this AI Personality, unless this list is empty in which
	// case anyone can select this AI personality
	ClientAllowList []string

	// MaxConcurrentInstances is the maximum number of games which can
	// be played simultaneously by this machine.
	MaxConcurrentInstances int
}

// QueueAI will register the AI with the lobby server and use the response
// to connect to the lobby socket server, returning the already authenticated
// websocket. The second result is the welcome message which should be stored
// for debugging errors with the server.
func QueueAI(cfg *AIConfig, auth *AuthToken) (*websocket.Conn, map[string]interface{}, error) {
	resp, err := requestLobbySocketServer(cfg, auth)
	if err != nil {
		return nil, nil, fmt.Errorf("requesting lobby-socket server: %w", err)
	}

	var conn *websocket.Conn
	headers := make(http.Header)
	headers.Add("Origin", utils.WEBSOCKET_ORIGIN)

	conn, _, err = websocket.DefaultDialer.Dial(resp.URL, headers)
	if err != nil {
		return nil, nil, fmt.Errorf("dialing %s: %w", resp.URL, err)
	}

	err = conn.SetWriteDeadline(time.Now().Add(utils.CONN_WRITE_TIMEOUT))
	if err != nil {
		closeConn(websocket.CloseInternalServerErr, conn)
		return nil, nil, fmt.Errorf("set write deadline: %w", err)
	}
	err = conn.WriteMessage(websocket.TextMessage, []byte(resp.JWT))
	if err != nil {
		closeConn(websocket.ClosePolicyViolation, conn)
		return nil, nil, fmt.Errorf("writing JWT: %w", err)
	}

	err = conn.SetWriteDeadline(time.Time{})
	if err != nil {
		closeConn(websocket.CloseInternalServerErr, conn)
		return nil, nil, fmt.Errorf("clear write deadline: %w", err)
	}
	err = conn.SetReadDeadline(time.Now().Add(utils.CONN_READ_TIMEOUT))
	if err != nil {
		closeConn(websocket.CloseInternalServerErr, conn)
		return nil, nil, fmt.Errorf("set read deadline: %w", err)
	}

	var msgType int
	var welcomeResponse []byte
	msgType, welcomeResponse, err = conn.ReadMessage()
	if err != nil {
		closeConn(websocket.ClosePolicyViolation, conn)
		return nil, nil, fmt.Errorf("reading welcome message: %w", err)
	}

	if msgType != websocket.TextMessage {
		closeConn(websocket.CloseUnsupportedData, conn)
		return nil, nil, fmt.Errorf("wrong message type (%d) for welcome message", msgType)
	}

	var parsedWelcomeRaw interface{}
	err = json.Unmarshal(welcomeResponse, &parsedWelcomeRaw)
	if err != nil {
		closeConn(websocket.ClosePolicyViolation, conn)
		return nil, nil, fmt.Errorf("could not parse welcome message (%s): %w", string(welcomeResponse), err)
	}

	err = conn.SetReadDeadline(time.Time{})
	if err != nil {
		closeConn(websocket.CloseInternalServerErr, conn)
		return nil, nil, fmt.Errorf("clear read deadline: %w", err)
	}

	var parsedWelcome map[string]interface{}
	switch v := parsedWelcomeRaw.(type) {
	case []interface{}:
		parsedWelcome = v[0].(map[string]interface{})
	case map[string]interface{}:
		parsedWelcome = v
	default:
		return nil, nil, fmt.Errorf("could not parse welcome message (%s): unknown type", string(welcomeResponse))
	}

	return conn, parsedWelcome, nil
}

type queueAIResponse struct {
	JWT string `json:"jwt"`
	URL string `json:"url"`
}

func requestLobbySocketServer(cfg *AIConfig, auth *AuthToken) (*queueAIResponse, error) {
	body := map[string]interface{}{
		"name":                     cfg.AIName,
		"uid":                      cfg.AIUID,
		"version":                  cfg.Version,
		"role":                     cfg.Role.Name(),
		"client_allow_list":        cfg.ClientAllowList,
		"max_concurrent_instances": cfg.MaxConcurrentInstances,
	}

	bodyMarshalled, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("marshalling body: %w", err)
	}

	client := &http.Client{Transport: &http2.Transport{}}
	var resp *http.Response
	var req *http.Request
	req, err = http.NewRequest(
		"POST",
		fmt.Sprintf("%s/api/1/play/ai", utils.API_BASE),
		bytes.NewBuffer(bodyMarshalled),
	)
	if err != nil {
		return nil, fmt.Errorf("preparing request: %w", err)
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("bearer %s", auth.Token))

	resp, err = client.Do(req)

	if err != nil {
		return nil, fmt.Errorf("on POST: %w", err)
	}

	var respBody []byte
	respBody, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading body: %w", err)
	}

	err = resp.Body.Close()
	if err != nil {
		return nil, fmt.Errorf("closing body: %w", err)
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("unexpected status code (got %d); body=%s", resp.StatusCode, string(respBody))
	}

	var parsedBody queueAIResponse
	err = json.Unmarshal(respBody, &parsedBody)
	if err != nil {
		return nil, fmt.Errorf("interpresting body (body=%s): %w", string(respBody), err)
	}

	return &parsedBody, nil
}
