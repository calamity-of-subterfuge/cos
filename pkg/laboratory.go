package pkg

// MAX_LABORATORY_PLACE_DISTANCE is the maximum distance between an science AI
// trying to place a laboratory and the laboratory
const MAX_LABORATORY_PLACE_DISTANCE = 1.2

// LABORATORY_RESOURCE_COST is maps from resource uids to the amount of them
// required to build a laboratory
var LABORATORY_RESOURCE_COST = map[string]int{
	"sapphire": 10,
}
