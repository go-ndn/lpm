package lpm

import (
	"fmt"
	"sync"
)

type threadSafeMatcher struct {
	u threadUnsafeMatcher
	sync.RWMutex
}

func newThreadSafeMatcher() *threadSafeMatcher {
	return &threadSafeMatcher{u: *newThreadUnsafeMatcher()}
}

func (m *threadSafeMatcher) Add(key fmt.Stringer, i interface{}) {
	m.Lock()
	m.u.Add(key, i)
	m.Unlock()
}

func (m *threadSafeMatcher) Remove(key fmt.Stringer) {
	m.Lock()
	m.u.Remove(key)
	m.Unlock()
}

func (m *threadSafeMatcher) Update(key fmt.Stringer, f func(interface{}) interface{}, lpm bool) {
	m.Lock()
	m.u.Update(key, f, lpm)
	m.Unlock()
}

func (m *threadSafeMatcher) UpdateAll(key fmt.Stringer, f func(string, interface{}) interface{}) {
	m.Lock()
	m.u.UpdateAll(key, f)
	m.Unlock()
}

func (m *threadSafeMatcher) Match(key fmt.Stringer) interface{} {
	m.RLock()
	defer m.RUnlock()
	return m.u.Match(key)
}
