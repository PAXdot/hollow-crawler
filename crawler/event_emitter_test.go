package crawler

import (
	"testing"
	"sync"

	"github.com/stretchr/testify/assert"
)

func TestEventEmitter(t *testing.T) {
	t.Run(".RegistersAndEmitsEvent", func(t *testing.T) {
		ee := NewEventEmitter()
		var eventTriggered bool

		ee.On("testEvent", func() {
			eventTriggered = true
		})

		ee.Emit("testEvent")

		assert.True(t, eventTriggered)
	})

	t.Run(".DoesNotTriggerUnregisteredEvents", func(t *testing.T) {
		ee := NewEventEmitter()
		ee.Emit("unregisteredEvent")

		assert.False(t, eventTriggered)
	})

	t.Run(".LocksRegistrationMutually", func(t *testing.T) {
		ee := NewEventEmitter()
		var wg sync.WaitGroup
		wg.Add(2)

		go func() {
			defer wg.Done()
			ee.On("testEvent", func() {})
		}()

		go func() {
			defer wg.Done()
			ee.On("testEvent", func() {})
		}()

		wg.Wait()
		ee.Emit("testEvent")
	})
}
