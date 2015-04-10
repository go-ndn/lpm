package lpm

import "sync"

type threadSafeMatcher struct {
	u threadUnsafeMatcher
	sync.RWMutex
}

func newThreadSafeMatcher() *threadSafeMatcher {
	return &threadSafeMatcher{u: *newThreadUnsafeMatcher()}
}

func (m *threadSafeMatcher) Update(s string, f func(interface{}) interface{}, lpm bool) {
	m.Lock()
	m.u.Update(s, f, lpm)
	m.Unlock()
}

func (m *threadSafeMatcher) UpdateAll(s string, f func(string, interface{}) interface{}) {
	m.Lock()
	m.u.UpdateAll(s, f)
	m.Unlock()
}

func (m *threadSafeMatcher) Match(s string, f func(interface{})) {
	m.RLock()
	m.u.Match(s, f)
	m.RUnlock()
}

func (m *threadSafeMatcher) Visit(f func(string, interface{}) interface{}) {
	m.Lock()
	m.u.Visit(f)
	m.Unlock()
}
