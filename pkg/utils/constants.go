package utils

import "time"

// API_BASE contains the base URL of the server instance.
const API_BASE = "https://calamityofsubterfuge.com"

// WEBSOCKET_ORIGIN is the value of the Origin header on websockets. This is
// required or the server will reject the websocket, although it's mainly to
// protect against browsers starting websockets on different websites.
const WEBSOCKET_ORIGIN = "https://calamityofsubterfuge.com"

// CONN_WRITE_TIMEOUT is the amount of time we can spend waiting for
// a write to be acknowledged on the socket before we close the socket.
var CONN_WRITE_TIMEOUT = 20 * time.Second

// CONN_READ_TIMEOUT is the amount of time we spend waiting for a read
// on the socket before we close the socket.
var CONN_READ_TIMEOUT = 20 * time.Second

// MAP_HEX_RADIUS is the radius of each of the hexes on the map which correspond
// to the player bases in world units. One world unit is 64 pixels at standard
// zoom.
const MAP_HEX_RADIUS = 10.0

// MAP_HEX_WALL_THICKNESS is the thickness of the walls on the edge of the map
// hexes.
const MAP_HEX_WALL_THICKNESS = 1.0

// VISION_DISTANCE describes the view distance of players. Note that view distance
// is not calculated as a circle - it is a square where VISION_DISTANCE is half
// the length of each side. This improves performance and better matches the
// rectangular nature of clients.
const VISION_DISTANCE = 10.0
