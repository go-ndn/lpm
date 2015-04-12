package lpm

import "sync"

type threadSafeMatcher struct {
	u  *threadUnsafeMatcher
	mu sync.RWMutex
}

func newThreadSafeMatcher() *threadSafeMatcher {
	return &threadSafeMatcher{u: newThreadUnsafeMatcher()}
}

func (m *threadSafeMatcher) Update(s string, f func(interface{}) interface{}, lpm bool) {
	m.mu.Lock()
	m.u.Update(s, f, lpm)
	m.mu.Unlock()
}

func (m *threadSafeMatcher) UpdateAll(s string, f func(string, interface{}) interface{}) {
	m.mu.Lock()
	m.u.UpdateAll(s, f)
	m.mu.Unlock()
}

func (m *threadSafeMatcher) Match(s string, f func(interface{})) {
	m.mu.RLock()
	m.u.Match(s, f)
	m.mu.RUnlock()
}

func (m *threadSafeMatcher) Visit(f func(string, interface{}) interface{}) {
	m.mu.Lock()
	m.u.Visit(f)
	m.mu.Unlock()
}
