package utils

import "fmt"

// Role is an enum describing all the roles within the game that
// players can have. Each player has uniquely one role.
type Role int

const (
	// RoleInvalid is the default value meaning that there was a bug
	// that caused the player to not have a role set.
	RoleInvalid Role = 0

	// RolePlayer is the role of humans
	RolePlayer Role = 1

	// RoleEconomyAI is the role of AI's acting in a pair to control the
	// fiscal policy of the team
	RoleEconomyAI Role = 2

	// RoleMilitaryAI is the role of AI's acting in a pair to control the
	// defense and foreign policy of the team
	RoleMilitaryAI Role = 3

	// RoleScienceAI is the role of AI's acting in a pair to control the
	// scientific policy of the team
	RoleScienceAI Role = 4
)

// RoleToName converts a role constant to its typical name.
func RoleToName(role Role) string {
	switch role {
	case RoleInvalid:
		return "invalid"
	case RolePlayer:
		return "player"
	case RoleEconomyAI:
		return "economy"
	case RoleMilitaryAI:
		return "military"
	case RoleScienceAI:
		return "science"
	default:
		return fmt.Sprint(role)
	}
}

func RoleFromName(role string) Role {
	switch role {
	case "player":
		return RolePlayer
	case "economy":
		return RoleEconomyAI
	case "military":
		return RoleMilitaryAI
	case "science":
		return RoleScienceAI
	default:
		return RoleInvalid
	}
}

// Name is an alternative way to call RoleToName
func (r Role) Name() string {
	return RoleToName(r)
}
