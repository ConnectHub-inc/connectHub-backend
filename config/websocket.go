package config

import "time"

const (
	// Max wait time when writing a message to the peer.
	WriteWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	PongWait       = 60 * time.Second
	PingMultiplier = 9

	// Send pings to peer with this period. Must be less than pongWait.
	PingPeriod = (PongWait * PingMultiplier) / 10

	// Max message size allowed from peer.
	MaxMessageSize = 10000

	// Max buffer size for messages.
	BufferSize = 4096

	// ChannelBufferSize is the buffer size for the channel.
	ChannelBufferSize = 256

	// PubSubGeneralChannel is the general channel for pubsub.
	PubSubGeneralChannel = "general"

	// PubSubChannelPrefix is the prefix for the channel channel.
	WelcomeMessage = "%s joined the channel"
	GoodbyeMessage = "%s left the channel"
)
