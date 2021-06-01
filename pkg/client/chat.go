package client

import (
	"github.com/calamity-of-subterfuge/cos/lib/rbuf"
	"github.com/calamity-of-subterfuge/cos/pkg/srvpkts"
	"github.com/calamity-of-subterfuge/cos/pkg/utils"
)

// Chat maintains the chat history and current set of chat authors
type Chat struct {
	// LocalChatAuthorsByUID contains all of the chat authors that would hear
	// any messages we send locally right now.
	LocalChatAuthorsByUID map[string]*ChatAuthor

	// RecentChatAuthorsByUID contains all of the chat authors that are in the
	// message history
	RecentChatAuthorsByUID map[string]*ChatAuthor

	// RecentChatAuthorUIDsToCount maps from the keys in RecentChatAuthorsByUID
	// to the number of messages by that chat author in the message history. Each
	// value is positive - when it reaches 0 the key is deleted from the map
	// and the chat author is removed from RecentChatAuthorsByUID.
	RecentChatAuthorUIDsToCount map[string]int

	// MessageHistory is the ring of recent messages
	MessageHistory *rbuf.PointerRingBuf

	messageHistoryLength int
}

// NewChat initializes a blank chat event that will need a game sync
// event to fully initialize
func NewChat(messageHistoryLength int) *Chat {
	return &Chat{
		messageHistoryLength: messageHistoryLength,
	}
}

// HandleMessage should be called on any packet received from the server to
// update the state of the chat.
func (c *Chat) HandleMessage(packet srvpkts.Packet) {
	switch v := packet.(type) {
	case *srvpkts.GameSyncPacket:
		c.gameSync(v)
	case *srvpkts.ChatAuthorAddedPacket:
		c.chatAuthorAdded(v)
	case *srvpkts.ChatAuthorUpdatePacket:
		c.chatAuthorUpdate(v)
	case *srvpkts.ChatAuthorRemovedPacket:
		c.chatAuthorRemoved(v)
	case *srvpkts.ChatMessagePacket:
		c.chatMessage(v)
	}
}

func (c *Chat) gameSync(packet *srvpkts.GameSyncPacket) {
	c.LocalChatAuthorsByUID = make(map[string]*ChatAuthor)
	c.RecentChatAuthorsByUID = make(map[string]*ChatAuthor)
	c.RecentChatAuthorUIDsToCount = make(map[string]int)
	c.MessageHistory = rbuf.NewPointerRingBuf(c.messageHistoryLength)

	for uid, author := range packet.ChatAuthors {
		c.LocalChatAuthorsByUID[uid] = (&ChatAuthor{}).Update(&author)
	}
}

func (c *Chat) chatAuthorAdded(packet *srvpkts.ChatAuthorAddedPacket) {
	if recentAuthor, found := c.RecentChatAuthorsByUID[packet.UID]; found {
		recentAuthor.Update(&packet.ChatAuthorSync)
		c.LocalChatAuthorsByUID[packet.UID] = recentAuthor
	} else {
		c.LocalChatAuthorsByUID[packet.UID] = (&ChatAuthor{}).Update(&packet.ChatAuthorSync)
	}
}

func (c *Chat) chatAuthorUpdate(packet *srvpkts.ChatAuthorUpdatePacket) {
	c.LocalChatAuthorsByUID[packet.UID].Update(&packet.ChatAuthorSync)
}

func (c *Chat) chatAuthorRemoved(packet *srvpkts.ChatAuthorRemovedPacket) {
	delete(c.LocalChatAuthorsByUID, packet.UID)
}

func (c *Chat) chatMessage(packet *srvpkts.ChatMessagePacket) {
	if c.MessageHistory.Readable == c.MessageHistory.N {
		popped := c.MessageHistory.A[c.MessageHistory.Beg].(*ChatMessage)

		oldCnt := c.RecentChatAuthorUIDsToCount[popped.ChatAuthor.UID]
		if oldCnt == 1 {
			delete(c.RecentChatAuthorUIDsToCount, popped.ChatAuthor.UID)
			delete(c.RecentChatAuthorsByUID, popped.ChatAuthor.UID)
		} else {
			c.RecentChatAuthorUIDsToCount[popped.ChatAuthor.UID] = oldCnt - 1
		}
	}

	var author *ChatAuthor
	var found bool
	if author, found = c.RecentChatAuthorsByUID[packet.AuthorUID]; found {
		c.RecentChatAuthorUIDsToCount[packet.AuthorUID]++
	} else {
		c.RecentChatAuthorUIDsToCount[packet.AuthorUID] = 1
		author = c.LocalChatAuthorsByUID[packet.AuthorUID]
		c.RecentChatAuthorsByUID[packet.AuthorUID] = author
	}

	toPush := make([]interface{}, 1)
	toPush[0] = &ChatMessage{
		GameTime:   packet.GameTime,
		Time:       utils.TimeFromUnix(packet.Time),
		ChatAuthor: author,
		Text:       packet.Text,
	}
	c.MessageHistory.PushAndMaybeOverwriteOldestData(toPush)
}
