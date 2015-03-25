package lpm

import (
	"fmt"
	"strings"
)

type threadUnsafeMatcher struct {
	table map[string]interface{}
}

func newThreadUnsafeMatcher() *threadUnsafeMatcher {
	return &threadUnsafeMatcher{table: make(map[string]interface{})}
}

func (m *threadUnsafeMatcher) Add(key fmt.Stringer, i interface{}) {
	m.Update(key, func(interface{}) interface{} { return i }, false)
}

func (m *threadUnsafeMatcher) Remove(key fmt.Stringer) {
	m.Update(key, func(interface{}) interface{} { return nil }, false)
}

func (m *threadUnsafeMatcher) findPrefix(s string, all bool) (prefix []string) {
	for {
		if _, ok := m.table[s]; ok {
			prefix = append(prefix, s)
			if !all {
				break
			}
		}
		i := strings.LastIndex(s, "/")
		if i == -1 {
			break
		}
		s = s[:i]
	}
	return
}

func (m *threadUnsafeMatcher) Update(key fmt.Stringer, f func(interface{}) interface{}, lpm bool) {
	s := key.String()
	if lpm {
		prefix := m.findPrefix(s, false)
		if len(prefix) == 0 {
			return
		}
		s = prefix[0]
	}
	m.table[s] = f(m.table[s])
	if m.table[s] == nil {
		delete(m.table, s)
	}
}

func (m *threadUnsafeMatcher) UpdateAll(key fmt.Stringer, f func(string, interface{}) interface{}) {
	for _, s := range m.findPrefix(key.String(), true) {
		m.table[s] = f(s, m.table[s])
		if m.table[s] == nil {
			delete(m.table, s)
		}
	}
}

func (m *threadUnsafeMatcher) Match(key fmt.Stringer) interface{} {
	prefix := m.findPrefix(key.String(), false)
	if len(prefix) == 0 {
		return nil
	}
	return m.table[prefix[0]]
}
