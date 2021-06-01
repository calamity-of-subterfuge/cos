package srvpkts

// ChatAuthorSync contains all the information the client needs to
// learn about a chat author
type ChatAuthorSync struct {
	// UID is the identifier for this chat author
	UID string `mapstructure:"uid" json:"uid"`

	// Name is the display name for this chat author
	Name string `mapstructure:"name" json:"name"`

	// Color is the font color for this player
	Color string `mapstructure:"color" json:"color"`

	// BonusClasses contains additional styles for this players text
	BonusClasses []string `mapstructure:"bonus_classes" json:"bonus_classes"`
}
