package lpm

import "strings"

type threadUnsafeMatcher struct {
	table map[string]interface{}
}

func newThreadUnsafeMatcher() *threadUnsafeMatcher {
	return &threadUnsafeMatcher{table: make(map[string]interface{})}
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

func (m *threadUnsafeMatcher) Update(s string, f func(interface{}) interface{}, lpm bool) {
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

func (m *threadUnsafeMatcher) UpdateAll(s string, f func(string, interface{}) interface{}) {
	for _, s := range m.findPrefix(s, true) {
		m.table[s] = f(s, m.table[s])
		if m.table[s] == nil {
			delete(m.table, s)
		}
	}
}

func (m *threadUnsafeMatcher) Match(s string, f func(interface{})) {
	prefix := m.findPrefix(s, false)
	if len(prefix) == 0 {
		return
	}
	f(m.table[prefix[0]])
}

func (m *threadUnsafeMatcher) Visit(f func(string, interface{}) interface{}) {
	for s := range m.table {
		m.table[s] = f(s, m.table[s])
		if m.table[s] == nil {
			delete(m.table, s)
		}
	}
}
