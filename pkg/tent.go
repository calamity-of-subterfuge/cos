package pkg

// MAX_TENT_PLACE_DISTANCE is the maximum distance between an economy AI
// trying to place a tent and the tent
const MAX_TENT_PLACE_DISTANCE = 1.0

// TENT_RESOURCE_COST is maps from resource uids to the amount of them
// required to build a tent
var TENT_RESOURCE_COST = map[string]int{
	"gold": 10,
}
