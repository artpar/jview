package transport

import "jview/protocol"

// Transport delivers A2UI messages to the engine.
type Transport interface {
	// Messages returns a channel of parsed messages.
	// The channel is closed when the transport is done.
	Messages() <-chan *protocol.Message

	// Errors returns a channel of transport errors.
	Errors() <-chan error

	// Start begins reading messages. Non-blocking.
	Start()

	// Stop terminates the transport.
	Stop()
}
