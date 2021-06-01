package client

import "time"

// ChatMessage describes a chat message
type ChatMessage struct {
	// GameTime when this message was received by the server
	GameTime float64

	// Time is the real time when the message was received by the server
	Time time.Time

	// ChatAuthor is a pointer to the chat author for this chat message;
	// this pointer is only valid while the message is in the Chat buffer,
	// otherwise it may be stale.
	ChatAuthor *ChatAuthor

	// Text of the message
	Text string
}
