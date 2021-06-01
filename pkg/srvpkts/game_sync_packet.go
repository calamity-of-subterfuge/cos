package srvpkts

// GameSyncPacketPlayer is used only within GameSyncPackets to describe
// the player
type GameSyncPacketPlayer struct {
	// UID is the UID of the primary game object for this player
	UID string `mapstructure:"uid" json:"uid"`

	// Team is the team that the player is on
	Team int `mapstructure:"team" json:"team"`

	// Role is the role that the player fulfilling
	Role string `mapstructure:"role" json:"role"`
}

// GameSyncPacketTeam is used only within GameSyncPackets to describe the player
// team
type GameSyncPacketTeam struct {
	// Resources is how many resources the team has, where the keys are
	// resource UIDs and the values are the amount of the resource that
	// the player team has
	Resources map[string]int `mapstructure:"resources" json:"resources"`
}

// GameSyncPacket provides everything necessary for a player which has no
// idea about the state of the game to fill in the state.
type GameSyncPacket struct {
	// Type is game-sync
	Type string `mapstructure:"type" json:"type"`

	// GameTime in fractional seconds.
	GameTime float64 `mapstructure:"game_time" json:"game_time"`

	// Player describes high-level information about the actual player
	Player GameSyncPacketPlayer `mapstructure:"game_time" json:"player"`

	// Team describes high-level information about the team the player is on
	Team GameSyncPacketTeam `mapstructure:"team" json:"team"`

	// Resources describes the resources that are in the game, where the
	// keys are resource uids and the values have info about that resource
	Resources map[string]ResourceSync `mapstructure:"resources" json:"resources"`

	// Players contains all the players that you can see, which at least
	// contains you. The keys are game object uids. The client should
	// assume that anything we send here is fair game to display.
	Players map[string]PlayerSync `mapstructure:"players" json:"players"`

	// DumbObjects contains all the dumb game objects on the map that you
	// can see. This may contain objects way out of vision that are known
	// in advance. The client should assume anything we send here is fair
	// game to display. The keys are the game object uids.
	DumbObjects map[string]GameObjectSync `mapstructure:"dumb_objects" json:"dumb_objects"`

	// SmartObjects contains all the smart game objects on the map that
	// you can see. The keys are the game object uids.
	SmartObjects map[string]SmartObjectSync `mapstructure:"smart_objects" json:"smart_objects"`

	// ChatAuthors contains all the things that the player can communicate
	// with right now, including the player itself.
	ChatAuthors map[string]ChatAuthorSync `mapstructure:"chat_authors" json:"chat_authors"`
}

func (p *GameSyncPacket) GetType() string {
	return "game-sync"
}

func (p *GameSyncPacket) PrepareForMarshal() {
	p.Type = p.GetType()
}

func init() {
	registerPacketParser("game-sync", func(parsed map[string]interface{}) (Packet, error) {
		return parseSinglePacketOfType(parsed, &GameSyncPacket{})
	})
}
