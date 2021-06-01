package client

import (
	"log"

	"github.com/calamity-of-subterfuge/cos/pkg/srvpkts"
)

// State describes the general state of the world from the perspective of the
// client right now. This is generally used to take over the generic packet
// management for AIs so they can focus on unit / player controls rather than
// just synchronizing the client state.
type State struct {
	// MyUID is the uid of the game object of the player for this client
	MyUID string

	// MyTeam is the team that the player for this client is on
	MyTeam int

	// GameTime is the current game time
	GameTime float64

	// PlayersByUID contains all the currently visible players mapped from
	// their uid.
	PlayersByUID map[string]*Player

	// StaticObjects contains all of the static objects in the game. These
	// do not move or change and you can see them regardless of vision.
	StaticObjects []GameObject

	// SmartObjectsByUID contains all of the currently visible smart objects
	// mapped from their uid.
	SmartObjectsByUID map[string]*SmartObject

	// GameObjectsByUID contains any generic game objects within vision.
	// Besides handling them for the purpose of collision these are largely
	// irrelevant to gameplay.
	GenericObjectsByUID map[string]*GameObject

	// ResourcesByUID contains all the resources on the clients team, mapped
	// from their uid.
	ResourcesByUID map[string]*Resource
}

// NewState initializes a blank state that will need the game sync packet
// in order to fill into normal representation. Typically the state can
// be considered invalid if the GameTime is 0.
func NewState() *State {
	return &State{}
}

// HandleMessage should be called whenever a new server packet is received. If
// the packet is relevant to the client state, this updates the client state
// appropriately.
func (s *State) HandleMessage(packet srvpkts.Packet) {
	switch v := packet.(type) {
	case *srvpkts.GameObjectAddedPacket:
		s.updateGameTime(v.GameTime)
		s.GenericObjectsByUID[v.Object.UID] = (&GameObject{}).Sync(&v.Object)
	case *srvpkts.GameObjectRemovedPacket:
		s.updateGameTime(v.GameTime)
		if _, found := s.PlayersByUID[v.UID]; found {
			delete(s.PlayersByUID, v.UID)
		} else if _, found := s.SmartObjectsByUID[v.UID]; found {
			delete(s.PlayersByUID, v.UID)
		} else {
			delete(s.GenericObjectsByUID, v.UID)
		}
	case *srvpkts.GameObjectUpdatePacket:
		s.updateGameTime(v.GameTime)
		if plyr, found := s.PlayersByUID[v.UID]; found {
			plyr.Update(v)
		} else if genObj, found := s.GenericObjectsByUID[v.UID]; found {
			genObj.Update(v)
		} else {
			log.Printf("ignoring update to unknown object %s", v.UID)
		}
	case *srvpkts.GameSyncPacket:
		s.handleGameSync(v)
	case *srvpkts.PlayerAddedPacket:
		s.updateGameTime(v.GameTime)
		s.PlayersByUID[v.Object.UID] = (&Player{}).Sync(&v.Object)
	case *srvpkts.SmartObjectAddedPacket:
		s.updateGameTime(v.GameTime)
		s.SmartObjectsByUID[v.Object.UID] = (&SmartObject{}).Sync(&v.Object)
	case *srvpkts.SmartObjectUpdatePacket:
		s.updateGameTime(v.GameTime)
		s.SmartObjectsByUID[v.UID].Update(v)
	case *srvpkts.TeamResourceChangedPacket:
		s.updateGameTime(v.GameTime)
		for uid, amt := range v.Resources {
			s.ResourcesByUID[uid].Amount = amt
		}
	}
}

func (s *State) handleGameSync(packet *srvpkts.GameSyncPacket) {
	s.MyUID = packet.Player.UID
	s.MyTeam = packet.Player.Team
	s.GameTime = packet.GameTime

	s.PlayersByUID = make(map[string]*Player, len(packet.Players))
	for _, plyr := range packet.Players {
		s.PlayersByUID[plyr.UID] = (&Player{}).Sync(&plyr)
	}

	s.StaticObjects = make([]GameObject, 0, len(packet.DumbObjects))
	for _, obj := range packet.DumbObjects {
		s.StaticObjects = append(s.StaticObjects, *(&GameObject{}).Sync(&obj))
	}

	s.SmartObjectsByUID = make(map[string]*SmartObject, len(packet.SmartObjects))
	for _, obj := range packet.SmartObjects {
		s.SmartObjectsByUID[obj.UID] = (&SmartObject{}).Sync(&obj)
	}

	s.GenericObjectsByUID = make(map[string]*GameObject)
	s.ResourcesByUID = make(map[string]*Resource, len(packet.Resources))
	for _, resPkt := range packet.Resources {
		s.ResourcesByUID[resPkt.UID] = (&Resource{Amount: packet.Team.Resources[resPkt.UID]}).Sync(&resPkt)
	}
}

func (s *State) updateGameTime(gameTime float64) {
	if s.GameTime < gameTime {
		s.GameTime = gameTime
	}
}
