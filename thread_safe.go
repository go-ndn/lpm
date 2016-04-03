package lpm

import "sync"

type threadSafeMatcher struct {
	Matcher
	sync.RWMutex
}

// NewThreadSafe creates a new thread-safe matcher.
func NewThreadSafe() Matcher {
	return &threadSafeMatcher{Matcher: New()}
}

func (m *threadSafeMatcher) Update(key []Component, f func(interface{}) interface{}, exist bool) {
	m.Lock()
	m.Matcher.Update(key, f, exist)
	m.Unlock()
}

func (m *threadSafeMatcher) UpdateAll(key []Component, f func([]Component, interface{}) interface{}, exist bool) {
	m.Lock()
	m.Matcher.UpdateAll(key, f, exist)
	m.Unlock()
}

func (m *threadSafeMatcher) Match(key []Component, f func(interface{}), exist bool) {
	m.RLock()
	m.Matcher.Match(key, f, exist)
	m.RUnlock()
}

func (m *threadSafeMatcher) Visit(f func([]Component, interface{}) interface{}) {
	m.Lock()
	m.Matcher.Visit(f)
	m.Unlock()
}
