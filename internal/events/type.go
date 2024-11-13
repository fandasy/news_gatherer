package events

import "context"

type Fetcher interface {
	Fetch(ctx context.Context, limit int) ([]Event, error)
}

type Processor interface {
	Process(ctx context.Context, e Event) error
}

const (
	Unknown int = iota
	Message
	Callback
)

type Event struct {
	Type int
	Text string
	Meta interface{}
}
