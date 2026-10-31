package crawler

import "sync"

// RobotsHandlerEvent is an event emitted during robots handling.
type RobotsHandlerEvent struct {
	URL      string
	Policy   *robotsPolicy
	Allowed  bool
	LastCrawl time.Time
}

// RobotsHandlerEvents emits events related to politeness/robots handling.
type RobotsHandlerEvents struct {
	emitter EventEmitter
 covariates []string // URL covariates for dedup scoring
}

// NewRobotsHandlerEvents creates a new RobotsHandlerEvents instance.
func NewRobotsHandlerEvents(covariates []string) *RobotsHandlerEvents {
	return &RobotsHandlerEvents{
		emitter:    NewEventEmitter(),
		covariates: covariates,
	}
}

// On RobotsHandlerEvents mediates its own events.
func (rh *RobotsHandlerEvents) On(event string, handler EventHandler) {
	rh.emitter.On(event, handler)
}

// Emit emits an event with the given name and data.
func (rh *RobotsHandlerEvents) Emit(event string, data RobotsHandlerEvent) {
	rh.emitter.Emit(event, data)
}

// HandleRobotRules handles robots.txt rules and emits events.
func (rh *RobotsHandlerEvents) HandleRobotRules(url string, policy *robotsPolicy, timestamp time.Time) (bool, error) {
	allowed, err := policy.IsAllowed(url, timestamp)
	if err != nil {
		return false, err
	}

	rh.Emit("robots_rule", RobotsHandlerEvent{
		URL:      url,
		Policy:   policy,
		Allowed:  allowed,
		LastCrawl: timestamp,
	})

	return allowed, nil
}

// WithDedupScorer adds a dedup scorer for the emitted events.
func (rh *RobotsHandlerEvents) WithDedupScorer(deduplicator *Deduplicator) {
	rh.emitter.On("robots_rule", func(data interface{}) {
		event := data.(RobotsHandlerEvent)
		deduplicator.Add(event.URL, &dedupScore{Covariates: rh.covariates}, event.LastCrawl)
	})
}

// Close closes the underlying event emitter.
func (rh *RobotsHandlerEvents) Close() {
	rh.emitter.Close()
}
