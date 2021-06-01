package srvpkts

// PlayerSync contains all of the information required for a client who has
// never learned about a player before to render them.
type PlayerSync struct {
	GameObjectSync

	// Role is the role of the player, see core.Role and core.RoleToName
	Role string `mapstructure:"role" json:"role"`

	// Team is the team of the player, generally 1-6 inclusive
	Team int `mapstructure:"team" json:"team"`
}
