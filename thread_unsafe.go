package lpm

import "strings"

type threadUnsafeMatcher struct {
	table map[string]interface{}
}

func newThreadUnsafeMatcher() *threadUnsafeMatcher {
	return &threadUnsafeMatcher{table: make(map[string]interface{})}
}

func prefixOf(s string) (prefix []string) {
	for {
		if s == "" {
			break
		}
		prefix = append(prefix, s)
		i := strings.LastIndex(s, "/")
		if i == -1 {
			break
		}
		s = s[:i]
	}
	return
}

func (m *threadUnsafeMatcher) update(s string, f func(string, interface{}) interface{}, exist, all bool) {
	for _, prefix := range prefixOf(s) {
		if _, ok := m.table[prefix]; exist && !ok {
			continue
		}
		m.table[prefix] = f(prefix, m.table[prefix])
		if m.table[prefix] == nil {
			delete(m.table, prefix)
		}
		if !all {
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
	for _, prefix := range prefixOf(s) {
		if _, ok := m.table[prefix]; exist && !ok {
			continue
		}
		f(m.table[prefix])
		break
	}
}

func (m *threadUnsafeMatcher) Visit(f func(string, interface{}) interface{}) {
	for s := range m.table {
		m.table[s] = f(s, m.table[s])
		if m.table[s] == nil {
			delete(m.table, s)
		}
	}
}
