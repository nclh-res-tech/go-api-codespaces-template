// Package messagebroker represents a stub for producing events to a Kafka topic,
// a real implementation would contain logic for retrying failed events etc
package messagebroker

import "context"

// Topic represents a topic in Kafka.
type Topic string

// Events represents an implementation that can produce events.
type Events struct{}

// NewEvents will instantiate a new instance of Events.
func NewEvents() *Events {
	return &Events{}
}

// Produce will produce an event on the given topic using the supplied payload.
func (e *Events) Produce(ctx context.Context, topic Topic, payload interface{}) {
	// TODO implement a Kafka producer implementation
	// ideally the payload would be a protobuf message from a shared schema package
	// for this exercise we will pass a struct that can be marshalled to JSON

	// The main implementation would happen asynchronously to avoid blocking the producing call
}
