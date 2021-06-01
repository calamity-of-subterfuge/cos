package client

import "github.com/calamity-of-subterfuge/cos/pkg/srvpkts"

// ChatAuthor describes a chat author from the perspective of the client
type ChatAuthor struct {
	// UID is the identifier for the chat author
	UID string

	// Name is the display name of the chat author
	Name string

	// Color is the CSS class that represents the color that the chat author
	// should be rendered with. Any CSS class loaded on the play page of the
	// website is valid.
	Color string

	// BonusClasses contains all of the additional CSS classes for the text of
	// this author. Any CSS class loaded on the play page of the website is
	// valid.
	BonusClasses []string
}

// Update this chat author to match the given sync information, then return
// this chat author
func (a *ChatAuthor) Update(sync *srvpkts.ChatAuthorSync) *ChatAuthor {
	a.UID = sync.UID
	a.Name = sync.Name
	a.Color = sync.Color
	a.BonusClasses = sync.BonusClasses
	return a
}
