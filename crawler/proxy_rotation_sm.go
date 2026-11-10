package crawler

import (
	"sync"
	"time"
)

const (
	Idle State = iota
	Active
	Exhausted
)

type State int

type ProxyRotationSM struct {
	state       State
	proxyCount  int
	proxyTick   time.Duration
	currentProxy Proxy
	mtx         sync.Mutex
}

type ProxyRotatorInterface interface {
	NextProxy() Proxy
}

func NewProxyRotationSM(initialProxy Proxy, count int, tick time.Duration) *ProxyRotationSM {
	return &ProxyRotationSM{
		state:      Idle,
		proxyCount: count,
		proxyTick:  tick,
		currentProxy: initialProxy,
	}
}

func (s *ProxyRotationSM) SetState(state State) {
	s.mtx.Lock()
	defer s.mtx.Unlock()
	s.state = state
}

func (s *ProxyRotationSM) GetState() State {
	s.mtx.Lock()
	defer s.mtx.Unlock()
	return s.state
}

func (s *ProxyRotationSM) SetProxy(proxy Proxy) {
	s.mtx.Lock()
	defer s.mtx.Unlock()
	s.currentProxy = proxy
}

func (s *ProxyRotationSM) GetProxy() Proxy {
	s.mtx.Lock()
	defer s.mtx.Unlock()
	return s.currentProxy
}

func (s *ProxyRotationSM) Rotate(proxyRotator ProxyRotatorInterface) {
	switch s.state {
	case Idle:
		s.currentProxy = proxyRotator.NextProxy()
		s.state = Active
		time.AfterFunc(s.proxyTick, func() {
			s.Rotate(proxyRotator)
		})
	case Active:
		// Continue rotating every 'tick' duration
		s.currentProxy = proxyRotator.NextProxy()
		time.AfterFunc(s.proxyTick, func() {
			s.Rotate(proxyRotator)
		})
	case Exhausted:
		// Optionally start rotating again after a delay
		time.AfterFunc(5*time.Minute, func() {
			s.state = Idle
			s.Rotate(proxyRotator)
		})
	}
}
