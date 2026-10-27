package crawler

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

// ProxyState represents the state of a proxy.
type ProxyState int

const (
	Idle ProxyState = iota
	InUse
	Failed
	CoolingDown
)

// ProxyStateMachine manages the state of proxies.
type ProxyStateMachine struct {
	proxies map[string]ProxyState
	mutex   sync.Mutex
}

// NewProxyStateMachine initializes a new ProxyStateMachine.
func NewProxyStateMachine() *ProxyStateMachine {
	return &ProxyStateMachine{
		proxies: make(map[string]ProxyState),
	}
}

// SetProxyState updates the state of a proxy.
func (m *ProxyStateMachine) SetProxyState(proxy string, state ProxyState) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.proxies[proxy] = state
}

// GetProxyState returns the state of a proxy.
func (m *ProxyStateMachine) GetProxyState(proxy string) (ProxyState, bool) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	state, exists := m.proxies[proxy]
	return state, exists
}

// TransitionProxyState updates the state of a proxy based on a transition function.
func (m *ProxyStateMachine) TransitionProxyState(proxy string, transition TransitionFunc) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	state, exists := m.proxies[proxy]
	if !exists {
		return fmt.Errorf("proxy %s not found", proxy)
	}

	nextState, err := transition(state)
	if err != nil {
		return err
	}

	m.proxies[proxy] = nextState
	return nil
}

// TransitionFunc defines a function to transition from one ProxyState to another.
type TransitionFunc func(ProxyState) (ProxyState, error)

// idleToInUse is a TransitionFunc that transitions from Idle to InUse state.
func idleToInUse(state ProxyState) (ProxyState, error) {
	if state != Idle {
		return state, fmt.Errorf("invalid state transition from %d to %d", state, InUse)
	}
	return InUse, nil
}

// inUseToFailed is a TransitionFunc that transitions from InUse to Failed state.
func inUseToFailed(state ProxyState) (ProxyState, error) {
	if state != InUse {
		return state, fmt.Errorf("invalid state transition from %d to %d", state, Failed)
	}
	return Failed, nil
}

// failedToCoolingDown is a TransitionFunc that transitions from Failed to CoolingDown state.
func failedToCoolingDown(state ProxyState) (ProxyState, error) {
	if state != Failed {
		return state, fmt.Errorf("invalid state transition from %d to %d", state, CoolingDown)
	}
	return CoolingDown, nil
}

// coolingDownToIdle is a TransitionFunc that transitions from CoolingDown to Idle state.
func coolingDownToIdle(state ProxyState) (ProxyState, error) {
	if state != CoolingDown {
		return state, fmt.Errorf("invalid state transition from %d to %d", state, Idle)
	}
	return Idle, nil
}

// TransitionToIdle transitions a proxy to Idle state after a cooling down period.
func (m *ProxyStateMachine) TransitionToIdle(proxy string) {
	time.AfterFunc(time.Minute, func() {
		m.TransitionProxyState(proxy, coolingDownToIdle)
	})
}
