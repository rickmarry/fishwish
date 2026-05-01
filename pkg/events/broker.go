package events

import (
	"context"
	"encoding/json"
)

type EventType string

const (
	EventSpotCreated     EventType = "spot.created"
	EventSpotUpdated     EventType = "spot.updated"
	EventReviewPosted    EventType = "review.posted"
	EventCatchLogged     EventType = "catch.logged"
	EventWeatherUpdated  EventType = "weather.updated"
)

type Event struct {
	Type      EventType       `json:"type"`
	Source    string          `json:"source"`
	Timestamp string          `json:"timestamp"`
	Data      json.RawMessage `json:"data"`
}

type Publisher interface {
	Publish(ctx context.Context, event Event) error
}

type Subscriber interface {
	Subscribe(ctx context.Context, eventType EventType, handler func(Event) error) error
}

type MemoryBroker struct {
	handlers map[EventType][]func(Event) error
}

func NewMemoryBroker() *MemoryBroker {
	return &MemoryBroker{
		handlers: make(map[EventType][]func(Event) error),
	}
}

func (b *MemoryBroker) Publish(ctx context.Context, event Event) error {
	for _, handler := range b.handlers[event.Type] {
		go func(h func(Event) error) {
			_ = h(event)
		}(handler)
	}
	return nil
}

func (b *MemoryBroker) Subscribe(_ context.Context, eventType EventType, handler func(Event) error) error {
	b.handlers[eventType] = append(b.handlers[eventType], handler)
	return nil
}
