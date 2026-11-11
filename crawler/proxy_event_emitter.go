package crawler

import (
	"errors"
	"sync"
)

// ProxyRotationEventType defines the types of events emitted by the ProxyRotator.
type ProxyRotationEventType string

const (
	ProxyRotationEventStart  ProxyRotationEventType = "start"
	ProxyRotationEventStop   ProxyRotationEventType = "stop"
	ProxyRotationEventRotate ProxyRotationEventType = "rotate"
	ProxyRotationEventError  ProxyRotationEventType = "error"
)

// ProxyRotationEvent represents an event emitted by the ProxyRotator.
type ProxyRotationEvent struct {
	Type      ProxyRotationEventType
	Proxy     *Proxy
	Error     error
}

// ProxyRotatorEventEmitter emits events related to proxy rotation.
type ProxyRotatorEventEmitter struct {
	mu        sync.Mutex
	listeners map[ProxyRotationEventType][]func(ProxyRotationEvent)
}

// NewProxyRotatorEventEmitter creates a new ProxyRotatorEventEmitter.
func NewProxyRotatorEventEmitter() *ProxyRotatorEventEmitter {
	return &ProxyRotatorEventEmitter{
		listeners: make(map[ProxyRotationEventType][]func(ProxyRotationEvent)),
	}
}

// Subscribe registers a listener function for the given event type.
func (e *ProxyRotatorEventEmitter) Subscribe(eventType ProxyRotationEventType, listener func(ProxyRotationEvent)) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	if _, ok := e.listeners[eventType]; !ok {
		return errors.New("invalid event type")
	}

	e.listeners[eventType] = append(e.listeners[eventType], listener)
	return nil
}

// Unsubscribe removes a listener function for the given event type.
func (e *ProxyRotatorEventEmitter) Unsubscribe(eventType ProxyRotationEventType, listener func(ProxyRotationEvent)) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	listeners, ok := e.listeners[eventType]
	if !ok || len(listeners) == 0 {
		return errors.New("no listeners registered for event type")
	}

	for i, l := range listeners {
		if l == listener {
			e.listeners[eventType] = append(listeners[:i], listeners[i+1:]...)
			return nil
		}
	}

	return errors.New("listener not found")
}

// Emit sends an event of the given type to all registered listeners.
func (e *ProxyRotatorEventEmitter) Emit(eventType ProxyRotationEventType, e Event) {
	e.mu.Lock()
	defer e.mu.Unlock()

	for _, listener := range e.listeners[eventType] {
		listener(e)
	}
}
