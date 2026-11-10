package crawler

import (
	"sync"
)

type EventType int

const (
	// DeduplicationEvent emitted when deduplication occurs
	DeduplicationEvent EventType = iota

	// ExportEvent emitted when an item is exported
	ExportEvent
)

type Event struct {
	EventType EventType
	Data      interface{}
}

type EventEmitter struct {
	events chan Event
	sync.Mutex
	done chan struct{}
}

func NewEventEmitter() *EventEmitter {
	return &EventEmitter{
		events: make(chan Event),
		done:   make(chan struct{}),
	}
}

func (ee *EventEmitter) Emit(eventType EventType, data interface{}) {
	ee.Lock()
	defer ee.Unlock()
	ee.events <- Event{eventType, data}
}

func (ee *EventEmitter) Listen() <-chan Event {
	return ee.events
}

func (ee *EventEmitter) Close() {
	close(ee.done)
	close(ee.events)
}
