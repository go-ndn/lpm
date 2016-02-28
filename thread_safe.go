package lpm

import "sync"

type threadSafeMatcher struct {
	Matcher
	sync.Mutex
}

// NewThreadSafe creates a new thread-safe matcher.
func NewThreadSafe() Matcher {
	return &threadSafeMatcher{Matcher: New()}
}

func (m *threadSafeMatcher) Update(s string, f func(interface{}) interface{}, exist bool) {
	m.Lock()
	m.Matcher.Update(s, f, exist)
	m.Unlock()
}

func (m *threadSafeMatcher) UpdateRaw(key []Component, f func(interface{}) interface{}, exist bool) {
	m.Lock()
	m.Matcher.UpdateRaw(key, f, exist)
	m.Unlock()
}

func (m *threadSafeMatcher) UpdateAll(s string, f func([]byte, interface{}) interface{}, exist bool) {
	m.Lock()
	m.Matcher.UpdateAll(s, f, exist)
	m.Unlock()
}

func (m *threadSafeMatcher) UpdateAllRaw(key []Component, f func([]byte, interface{}) interface{}, exist bool) {
	m.Lock()
	m.Matcher.UpdateAllRaw(key, f, exist)
	m.Unlock()
}

func (m *threadSafeMatcher) Match(s string, f func(interface{}), exist bool) {
	m.Lock()
	m.Matcher.Match(s, f, exist)
	m.Unlock()
}

func (m *threadSafeMatcher) MatchRaw(key []Component, f func(interface{}), exist bool) {
	m.Lock()
	m.Matcher.MatchRaw(key, f, exist)
	m.Unlock()
}

func (m *threadSafeMatcher) Visit(f func(string, interface{}) interface{}) {
	m.Lock()
	m.Matcher.Visit(f)
	m.Unlock()
}
