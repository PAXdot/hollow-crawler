package crawler

type RobotsHandlerEvents struct {
	OnNewPolicy     func(Policy)
	OnPolicyUpdate  func(Policy)
	OnPolicyDiscard func(Policy)
}

type Policy struct {
	URL       string
	Allow     bool
	Rules     map[string]PolicyRule
	LastUpdated time.Time
}

type PolicyRule struct {
	Name  string
	Value string
	Allow bool
}

var DefaultRobotsHandlerEvents = RobotsHandlerEvents{
	OnNewPolicy:     func(Policy) {},
	OnPolicyUpdate:  func(Policy) {},
	OnPolicyDiscard: func(Policy) {},
}
