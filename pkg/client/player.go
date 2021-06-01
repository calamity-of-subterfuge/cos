package client

import (
	"github.com/calamity-of-subterfuge/cos/v2/pkg/srvpkts"
	"github.com/calamity-of-subterfuge/cos/v2/pkg/utils"
)

// Player describes a player in the world from the perspective of the client
type Player struct {
	// GameObject for this player
	GameObject *GameObject

	// Role this player is fulfilling
	Role utils.Role

	// Team this player is on
	Team int
}

// Sync this player to match the given sync information
func (p *Player) Sync(sync *srvpkts.PlayerSync) *Player {
	p.GameObject.Sync(&sync.GameObjectSync)
	p.Role = utils.RoleFromName(sync.Role)
	p.Team = sync.Team
	return p
}

// Update this player with the given information
func (p *Player) Update(upd *srvpkts.GameObjectUpdatePacket) *Player {
	p.GameObject.Update(upd)
	return p
}
