package crawler

// ProxyEventKind categorizes proxy rotation lifecycle events.
type ProxyEventKind int

const (
	// ProxyAcquired is emitted when a proxy is successfully acquired from the pool.
	ProxyAcquired ProxyEventKind = iota
	// ProxyReleased is emitted when a proxy is returned to the pool.
	ProxyReleased
	// ProxyBlocked is emitted when a proxy is marked as blocked.
	ProxyBlocked
	// ProxyRotated is emitted when the active proxy is rotated for any reason.
	ProxyRotated
	// ProxyExhausted is emitted when all proxies are blocked or unavailable.
	ProxyExhausted
)

// ProxyEvent carries details about a proxy lifecycle change.
type ProxyEvent struct {
	Kind   ProxyEventKind
	Proxy  string
	Reason string
}

// proxyRotationEmitter is a small event emitter for proxy rotation.
type proxyRotationEmitter struct {
	listeners []func(ProxyEvent)
}

// newProxyRotationEmitter creates a new emitter.
func newProxyRotationEmitter() *proxyRotationEmitter {
	return &proxyRotationEmitter{}
}

// On registers a listener for proxy rotation events.
func (e *proxyRotationEmitter) On(fn func(ProxyEvent)) {
	if fn != nil {
		e.listeners = append(e.listeners, fn)
	}
}

// Emit sends an event to all registered listeners.
func (e *proxyRotationEmitter) Emit(ev ProxyEvent) {
	for _, fn := range e.listeners {
		fn(ev)
	}
}
