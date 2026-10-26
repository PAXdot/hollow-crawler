package crawler

import (
	"fmt"
	"sync"
)

// EventEmitter emits events related to URL frontier
type EventEmitter struct {
	listeners      map[string][]func()
	registrationMU sync.Mutex
}

// NewEventEmitter creates a new event emitter
func NewEventEmitter() *EventEmitter {
	return &EventEmitter{
		listeners: make(map[string][]func()),
	}
}

// On registers a callback function for a specific event
func (ee *EventEmitter) On(event string, callback func()) {
	ee.registrationMU.Lock()
	defer ee.registrationMU.Unlock()

	ee.listeners[event] = append(ee.listeners[event], callback)
}

// Emit triggers registered callback functions for a specific event
func (ee *EventEmitter) Emit(event string, args ...interface{}) {
	ee.registrationMU.Lock()
	callbacks, exists := ee.listeners[event]
	ee.registrationMU.Unlock()

	if !exists {
		return
	}

	for _, callback := range callbacks {
		callback()
	}

	fmt.Printf("Emitted '%s' event with args: %+v\n", event, args)
}
