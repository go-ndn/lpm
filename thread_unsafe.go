package lpm

import "strings"

type threadUnsafeMatcher struct {
	table map[string]interface{}
}

func newThreadUnsafeMatcher() *threadUnsafeMatcher {
	return &threadUnsafeMatcher{table: make(map[string]interface{})}
}

func prefixOf(s string) string {
	i := strings.LastIndex(s, "/")
	if i == -1 {
		return ""
	}
	return s[:i]
}

func (m *threadUnsafeMatcher) update(s string, f func(string, interface{}) interface{}, exist, all bool) {
	for s != "" {
		v, ok := m.table[s]
		if exist && !ok {
			s = prefixOf(s)
			continue
		}
		v = f(s, v)
		if v == nil {
			delete(m.table, s)
		} else {
			m.table[s] = v
		}
		if all {
			s = prefixOf(s)
		} else {
			break
		}
	}
}

func (m *threadUnsafeMatcher) Update(s string, f func(interface{}) interface{}, exist bool) {
	m.update(s, func(_ string, v interface{}) interface{} {
		return f(v)
	}, exist, false)
}

func (m *threadUnsafeMatcher) UpdateAll(s string, f func(string, interface{}) interface{}, exist bool) {
	m.update(s, f, exist, true)
}

func (m *threadUnsafeMatcher) Match(s string, f func(interface{}), exist bool) {
	for s != "" {
		v, ok := m.table[s]
		if exist && !ok {
			s = prefixOf(s)
			continue
		}
		f(v)
		break
	}
}

func (m *threadUnsafeMatcher) Visit(f func(string, interface{}) interface{}) {
	for s, v := range m.table {
		v = f(s, v)
		if v == nil {
			delete(m.table, s)
		} else {
			m.table[s] = v
		}
	}
}
