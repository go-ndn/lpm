package lpm

import "sync"

type threadSafeMatcher struct {
	Matcher
	sync.RWMutex
}

func NewThreadSafe() Matcher {
	return &threadSafeMatcher{Matcher: New()}
}

func (m *threadSafeMatcher) Update(s string, f func(interface{}) interface{}, exist bool) {
	m.Lock()
	m.Matcher.Update(s, f, exist)
	m.Unlock()
}

func (m *threadSafeMatcher) UpdateAll(s string, f func(string, interface{}) interface{}, exist bool) {
	m.Lock()
	m.Matcher.UpdateAll(s, f, exist)
	m.Unlock()
}

func (m *threadSafeMatcher) Match(s string, f func(interface{}), exist bool) {
	m.RLock()
	m.Matcher.Match(s, f, exist)
	m.RUnlock()
}

func (m *threadSafeMatcher) Visit(f func(string, interface{}) interface{}) {
	m.Lock()
	m.Matcher.Visit(f)
	m.Unlock()
}
