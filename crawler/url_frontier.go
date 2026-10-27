package crawler

import (
	"context"
	"sync"
)

// URLFrontierEventType defines the types of events emitted by the URLFrontier.
type URLFrontierEventType int

const (
	URLDiscovered     URLFrontierEventType = iota
	URLProcessed
	URLError
)

// URLFrontierEvent holds the details of an URL frontier event.
type URLFrontierEvent struct {
	EventType URLFrontierEventType
	URL       string
	Err       error
}

// URLFrontierEventSubscriber is an interface for objects that want to receive URLFrontier events.
type URLFrontierEventSubscriber interface {
	OnURLFrontierEvent(context.Context, *URLFrontierEvent)
}

// URLFrontier is a high-throughput concurrent crawler and scraper's URL frontier.
type URLFrontier struct {
	events chan *URLFrontierEvent
	mutex  sync.Mutex
	closed bool
}

// NewURLFrontier initializes a new URLFrontier with a given buffer size for events.
func NewURLFrontier(bufferSize int) *URLFrontier {
	return &URLFrontier{
		events: make(chan *URLFrontierEvent, bufferSize),
	}
}

// Subscribe adds an event subscriber to the URLFrontier.
func (f *URLFrontier) Subscribe(subscriber URLFrontierEventSubscriber, ctx context.Context) {
	go func() {
		for {
			select {
			case event := <-f.events:
				subscriber.OnURLFrontierEvent(ctx, event)
			case <-ctx.Done():
				return
			}
		}
	}()
}

// Emit sends an event on the URLFrontier channel.
func (f *URLFrontier) Emit(eventType URLFrontierEventType, url string, err error) {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	if f.closed {
		return
	}

	f.events <- &URLFrontierEvent{
		EventType: eventType,
		URL:       url,
		Err:       err,
	}
}

// Close stops accepting new events and waits for the event consumers to finish.
func (f *URLFrontier) Close() {
	f.mutex.Lock()
	f.closed = true
	f.mutex.Unlock()
	close(f.events)
}
