package pkg

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
