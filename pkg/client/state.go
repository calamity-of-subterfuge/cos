package client

import (
	"log"

	"github.com/calamity-of-subterfuge/cos/pkg/srvpkts"
	"github.com/calamity-of-subterfuge/cos/pkg/utils"
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

	// MyRole is the role of the player for this client
	MyRole utils.Role

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

	onSelfLoaded                    []func(*Player)
	onControllableSmartObjectLoaded []func(*SmartObject)
	onSelfLost                      []func(*Player)
	onControllableSmartObjectLost   []func(*SmartObject)
}

// NewState initializes a blank state that will need the game sync packet
// in order to fill into normal representation. Typically the state can
// be considered invalid if the GameTime is 0.
func NewState() *State {
	return &State{}
}

// OnSelfLoaded will register the given listener to be called whenever
// the Player with uid s.MyUID is loaded from a packet. Typically this
// is on game sync.
func (s *State) OnSelfLoaded(listener func(*Player)) {
	if s.onSelfLoaded == nil {
		s.onSelfLoaded = []func(*Player){listener}
	} else {
		s.onSelfLoaded = append(s.onSelfLoaded, listener)
	}
}

// OnSelfLost will register the given listener to be called whenever
// the Player with uid s.MyUID is removed from the game. Typically
// is the first stage of a game sync packet.
func (s *State) OnSelfLost(listener func(*Player)) {
	if s.onSelfLost == nil {
		s.onSelfLost = []func(*Player){listener}
	} else {
		s.onSelfLost = append(s.onSelfLost, listener)
	}
}

// OnControllableSmartObjectLoaded is called when a new smart object
// that the player can control is loaded from a packet, such as via
// a game sync or because it just came into vision or it was just
// created.
func (s *State) OnControllableSmartObjectLoaded(listener func(*SmartObject)) {
	if s.onControllableSmartObjectLoaded == nil {
		s.onControllableSmartObjectLoaded = []func(*SmartObject){listener}
	} else {
		s.onControllableSmartObjectLoaded = append(s.onControllableSmartObjectLoaded, listener)
	}
}

// OnControllableSmartObjectLost is called whenever a smart object which
// is controlled by the player is lost, such as at the beginning of a game
// sync, because it died, or because we lost vision of it.
func (s *State) OnControllableSmartObjectLost(listener func(*SmartObject)) {
	if s.onControllableSmartObjectLost == nil {
		s.onControllableSmartObjectLost = []func(*SmartObject){listener}
	} else {
		s.onControllableSmartObjectLost = append(s.onControllableSmartObjectLost, listener)
	}
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
		if ov, found := s.PlayersByUID[v.UID]; found {
			if v.UID == s.MyUID && s.onSelfLost != nil {
				for _, listener := range s.onSelfLost {
					listener(ov)
				}
			}

			delete(s.PlayersByUID, v.UID)
		} else if ov, found := s.SmartObjectsByUID[v.UID]; found {
			if ov.ControllingTeam == s.MyTeam && ov.ControllingRole == s.MyRole && s.onControllableSmartObjectLost != nil {
				for _, listener := range s.onControllableSmartObjectLost {
					listener(ov)
				}
			}

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
		newPlayer := (&Player{}).Sync(&v.Object)
		s.PlayersByUID[v.Object.UID] = newPlayer

		if newPlayer.GameObject.UID == s.MyUID && s.onSelfLoaded != nil {
			for _, listener := range s.onSelfLoaded {
				listener(newPlayer)
			}
		}
	case *srvpkts.SmartObjectAddedPacket:
		s.updateGameTime(v.GameTime)
		newSO := (&SmartObject{}).Sync(&v.Object)
		s.SmartObjectsByUID[v.Object.UID] = newSO

		if newSO.ControllingTeam == s.MyTeam && newSO.ControllingRole == s.MyRole && s.onControllableSmartObjectLoaded != nil {
			for _, listener := range s.onControllableSmartObjectLoaded {
				listener(newSO)
			}
		}
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
	if s.MyUID != "" && s.onSelfLost != nil {
		me, found := s.PlayersByUID[s.MyUID]
		if found {
			for _, listener := range s.onSelfLost {
				listener(me)
			}
		}
	}

	if s.SmartObjectsByUID != nil && s.onControllableSmartObjectLost != nil {
		for _, so := range s.SmartObjectsByUID {
			if so.ControllingTeam == s.MyTeam && so.ControllingRole == s.MyRole {
				for _, listener := range s.onControllableSmartObjectLost {
					listener(so)
				}
			}
		}
	}

	s.MyUID = packet.Player.UID
	s.MyTeam = packet.Player.Team
	s.MyRole = utils.RoleFromName(packet.Player.Role)
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

	if s.onSelfLoaded != nil {
		me, found := s.PlayersByUID[s.MyUID]
		if found {
			for _, listener := range s.onSelfLoaded {
				listener(me)
			}
		}
	}

	if s.onControllableSmartObjectLoaded != nil {
		for _, so := range s.SmartObjectsByUID {
			if so.ControllingRole == s.MyRole && so.ControllingTeam == s.MyTeam {
				for _, listener := range s.onControllableSmartObjectLoaded {
					listener(so)
				}
			}
		}
	}
}

func (s *State) updateGameTime(gameTime float64) {
	if s.GameTime < gameTime {
		s.GameTime = gameTime
	}
}
