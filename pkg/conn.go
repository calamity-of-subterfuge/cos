package pkg

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/gorilla/websocket"
)

// PreparablePacket can be implemented by packets to get a chance to initialize
// things prior to being marshaled and sent across the connection. This is most
// typically used for setting the Type field on the packet.
type PreparablePacket interface {
	// PrepareForMarshal is called prior to marshalling the packet which is to
	// be sent across the connection.
	PrepareForMarshal()
}

// ReceivedMessage describes a message along with the connection it was
// received on.
type ReceivedMessage struct {
	// ConnectionUID is the UID of the WebsocketChannelConn that this message
	// came from.
	ConnectionUID string

	// Message is the parsed message that was received. Note that this may have
	// been only part of the actual logical message frame, since the calamity of
	// subterfuge protocol allows multiple messages per message frame to reduce
	// overhead on small messages.
	Message map[string]interface{}
}

// Conn is a convenience wrapper around a basic websocket connection which uses
// channels for send/receive of packets in the format expected by the calamity
// of subterfuge lobby socket and game socket protocols. The connection itself
// manages the required goroutines that read from the send queue and write to
// the receive queue, which can be canceled using Close.
type Conn struct {
	// UID is the identifier of this connection which is forwarded alongside all
	// messages to the receiving channel. This allows multiple connections to
	// use the same receive channel if it's desirable to do so. It may be left
	// blank if the receive channel only has a single connection and hence the
	// UID of the connection is superfluous.
	UID string

	// SendQueue is the channel which the Conn reads from in order to write to the
	// actual websocket.
	SendQueue chan interface{}

	recvQueue    chan ReceivedMessage
	closedQueue  chan string
	cancelSignal chan struct{}
	conn         *websocket.Conn
}

// NewConn takes over management of the given websocket connection and returns
// the managed connection. It is not safe to use the websocket directly once
// this function is called.
//
// The uid is used to distinguish messages from this connection if there are
// multiple connections using the same receive queue and closed queue. It may
// be left as a blank string if the receive and closed queues only have one
// connection and hence the uid is not required to distinguish the source.
//
// A message is written to the receive queue whenever the server sends us a
// message. Our uid is written to the closedQueue exactly once when the
// underlying websocket connection is closed.
func NewConn(conn *websocket.Conn, uid string, recvQueue chan ReceivedMessage, closedQueue chan string) *Conn {
	res := &Conn{
		UID:          uid,
		SendQueue:    make(chan interface{}, 128),
		recvQueue:    recvQueue,
		closedQueue:  closedQueue,
		cancelSignal: make(chan struct{}, 1),
		conn:         conn,
	}

	go res.manageSend()
	go res.manageRecv()

	return res
}

// Close the connection if it's not already closed. If the socket is
// currently open, this will result in our uid being written to the
// closedQueue after a short delay.
func (c *Conn) Close() {
	select {
	case c.cancelSignal <- struct{}{}:
	default:
	}
}

func (c *Conn) manageSend() {
	lowerTimeout := CONN_READ_TIMEOUT
	if CONN_WRITE_TIMEOUT < lowerTimeout {
		lowerTimeout = CONN_WRITE_TIMEOUT
	}

	pingInterval := (lowerTimeout * 9) / 10
	pingTicker := time.NewTicker(pingInterval)

	packets := make([]interface{}, 0, 1)

outerLoop:
	for {
		select {
		case packet := <-c.SendQueue:
			// Batching sends can significantly improve performance
			packets = append(packets, packet)
		readPacketsLoop:
			for len(packets) < 16 {
				select {
				case nextPacket := <-c.SendQueue:
					packets = append(packets, nextPacket)
				default:
					break readPacketsLoop
				}
			}

			err := c.sendPackets(packets)
			packets = packets[:0]
			if err != nil {
				if !errors.Is(err, net.ErrClosed) {
					log.Printf("error sending packets to conn %s: %v", c.UID, err)
				}
				c.Close()
				break outerLoop
			}
		case <-pingTicker.C:
			err := c.conn.SetWriteDeadline(time.Now().Add(CONN_WRITE_TIMEOUT))
			if err != nil {
				log.Printf("Error setting write dealding for ping to %s: %v", c.UID, err)
				c.Close()
				break outerLoop
			}

			err = c.conn.WriteMessage(websocket.PingMessage, nil)
			if err != nil {
				if !errors.Is(err, net.ErrClosed) {
					log.Printf("Failed to write ping to connection %s: %v", c.UID, err)
				}
				c.Close()
				break outerLoop
			}
		case <-c.cancelSignal:
			// It's nice to try and send the last couple packets here so that
			// when using this you don't have to do this awkward thing where
			// when you want to close the connection cleanly you need to write
			// the packets and wait "a bit". This isn't perfect since if there's
			// too many packets in the queue it still won't get them all, but
			// that should basically never happen
		readPacketsLoop2:
			for len(packets) < 16 {
				select {
				case nextPacket := <-c.SendQueue:
					packets = append(packets, nextPacket)
				default:
					break readPacketsLoop2
				}
			}

			if len(packets) > 0 {
				err := c.sendPackets(packets)
				if err != nil && !errors.Is(err, net.ErrClosed) {
					log.Printf("error sending final packets to conn %s: %v", c.UID, err)
				}
			}
			c.Close()
			break outerLoop
		}
	}

	cerr := c.conn.SetWriteDeadline(time.Now().Add(CONN_WRITE_TIMEOUT))
	if cerr != nil {
		log.Printf("failed to set write deadline on %s for close code: %v", c.UID, cerr)
	} else {
		cerr = c.conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, "forcibly disconnecting"))
		if cerr != nil && !errors.Is(cerr, net.ErrClosed) {
			log.Printf("failed to send close code to conn %s: %v", c.UID, cerr)
		}
	}

	cerr = c.conn.Close()
	if cerr != nil && !errors.Is(cerr, net.ErrClosed) {
		log.Printf("failed to close connection %s on send close: %v", c.UID, cerr)
	}
	c.closedQueue <- c.UID
}

func (c *Conn) sendPackets(packets []interface{}) error {
	for _, pkt := range packets {
		if preparablePacket, ok := pkt.(PreparablePacket); ok {
			preparablePacket.PrepareForMarshal()
		}
	}

	marshalledPackets, err := json.Marshal(packets)
	if err != nil {
		for _, pkt := range packets {
			_, subErr := json.Marshal(pkt)
			if subErr != nil {
				return fmt.Errorf("failed to marshal packet %v: %w", pkt, err)
			}
		}
		return fmt.Errorf("failed to marshal packets despite each individual packet marshalling fine! packets: %v, err: %w", packets, err)
	}

	err = c.conn.SetWriteDeadline(time.Now().Add(CONN_WRITE_TIMEOUT))
	if err != nil {
		return fmt.Errorf("failed to set write deadline: %w", err)
	}

	err = c.conn.WriteMessage(websocket.TextMessage, marshalledPackets)
	if err != nil {
		return fmt.Errorf("failed to write message: %w", err)
	}
	return nil
}

func (c *Conn) manageRecv() {
	// Receive is naturally cancelled promptly by manageSend
	// closing the websocket

	c.conn.SetPongHandler(func(string) error {
		return c.conn.SetReadDeadline(time.Now().Add(CONN_READ_TIMEOUT))
	})

	for {
		err := c.conn.SetReadDeadline(time.Now().Add(CONN_READ_TIMEOUT))
		if err != nil {
			log.Printf("Failed to set read deadline for %s: %v", c.UID, err)
			c.Close()
			break
		}

		var messageType int
		var message []byte
		messageType, message, err = c.conn.ReadMessage()
		if err != nil {
			log.Printf("Failed to read message from %s: %v", c.UID, err)
			c.Close()
			break
		}

		if messageType != websocket.TextMessage {
			log.Printf("Invalid incoming message type from %s: %v", c.UID, messageType)
			c.Close()
			break
		}

		decoder := json.NewDecoder(bytes.NewBuffer(message))
		decoder.UseNumber()

		var decodedMessage interface{}
		err = decoder.Decode(&decodedMessage)
		if err != nil {
			log.Printf("Failed to decode incoming message from %s: %v", c.UID, err)
			c.Close()
			break
		}

		if arr, ok := decodedMessage.([]interface{}); ok {
			for _, packet := range arr {
				c.recvQueue <- ReceivedMessage{
					ConnectionUID: c.UID,
					Message:       packet.(map[string]interface{}),
				}
			}
		} else if packet, ok := decodedMessage.(map[string]interface{}); ok {
			c.recvQueue <- ReceivedMessage{
				ConnectionUID: c.UID,
				Message:       packet,
			}
		} else {
			log.Printf("Unknown format for incoming message from %s: %v", c.UID, decodedMessage)
			c.Close()
			break
		}
	}
}
