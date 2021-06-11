package client

import "github.com/calamity-of-subterfuge/cos/pkg/utils"

// TeamRoleUIDLookup provides a lookup of the set of UIDs for players
// which are within a particular team and role.
type TeamRoleUIDLookup map[int]map[utils.Role]map[string]struct{}

// Add the given uid with the given team and role
func (l TeamRoleUIDLookup) Add(team int, role utils.Role, uid string) bool {
	roleToUIDS, ok := l[team]
	if !ok {
		roleToUIDS = make(map[utils.Role]map[string]struct{})
		l[team] = roleToUIDS
	}

	var uids map[string]struct{}
	uids, ok = roleToUIDS[role]
	if !ok {
		uids = make(map[string]struct{})
		roleToUIDS[role] = uids
	}

	_, ok = uids[uid]
	if ok {
		return false
	}
	uids[uid] = struct{}{}
	return true
}

// Remove the given uid with the given team and role
func (l TeamRoleUIDLookup) Remove(team int, role utils.Role, uid string) bool {
	roleToUIDS, ok := l[team]
	if !ok {
		return false
	}

	var uids map[string]struct{}
	uids, ok = roleToUIDS[role]
	if !ok {
		return false
	}

	_, ok = uids[uid]
	if !ok {
		return false
	}

	delete(uids, uid)
	if len(uids) == 0 {
		delete(roleToUIDS, role)

		if len(roleToUIDS) == 0 {
			delete(l, team)
		}
	}
	return true
}

// EachOnTeam calls the given function on each player in the given team.
func (l TeamRoleUIDLookup) EachOnTeam(team int, fnc func(string)) {
	rolesToUIDs, ok := l[team]

	if !ok {
		return
	}

	for _, uids := range rolesToUIDs {
		for uid := range uids {
			fnc(uid)
		}
	}
}

// EachOnTeamWithRole calls the given function with the uid of every player
// in the given team with the given role
func (l TeamRoleUIDLookup) EachOnTeamWithRole(team int, role utils.Role, fnc func(string)) {
	rolesToUIDs, ok := l[team]
	if !ok {
		return
	}

	uids, hasWithRole := rolesToUIDs[role]
	if !hasWithRole {
		return
	}

	for uid := range uids {
		fnc(uid)
	}
}
