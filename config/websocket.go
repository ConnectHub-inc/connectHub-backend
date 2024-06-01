package config

import "time"

const (
	// Max wait time when writing a message to the peer.
	WriteWait = 10 * time.Second

	// Max wait time for the peer to read the next pong message.
	PongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	PingPeriod = (PongWait * 9) / 10

	// Max message size allowed from peer.
	MaxMessageSize = 10000
)
