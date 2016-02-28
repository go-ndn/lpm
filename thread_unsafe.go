package lpm

import (
	"bytes"

	"github.com/go-ndn/tlv"
)

type threadUnsafeMatcher struct {
	table map[string]interface{}
	b     []byte
}

// New creates a new thread-unsafe matcher.
func New() Matcher {
	return &threadUnsafeMatcher{
		table: make(map[string]interface{}),
		b:     make([]byte, tlv.MaxSize),
	}
}

func prefixB(b []byte) []byte {
	i := bytes.LastIndexByte(b, '/')
	if i == -1 {
		return nil
	}
	return b[:i]
}

func mustEscape(c byte) bool {
	switch {
	case 'A' <= c && c <= 'Z':
		fallthrough
	case 'a' <= c && c <= 'z':
		fallthrough
	case '0' <= c && c <= '9':
		fallthrough
	case c == '-':
		fallthrough
	case c == '_':
		fallthrough
	case c == '.':
		fallthrough
	case c == '~':
		return false
	default:
		return true
	}
}

const (
	hex = "0123456789ABCDEF"
)

func escape(b []byte, key []Component) []byte {
	var n int
	for _, raw := range key {
		b[n] = '/'
		n++
		for _, c := range raw {
			if mustEscape(c) {
				b[n] = '%'
				b[n+1] = hex[c>>4]
				b[n+2] = hex[c&15]
				n += 3
			} else {
				b[n] = c
				n++
			}
		}
	}
	return b[:n]
}

func cpy(b []byte, s string) []byte {
	n := copy(b, s)
	return b[:n]
}

func (m *threadUnsafeMatcher) updateB(b []byte, f func([]byte, interface{}) interface{}, exist, all bool) {
	for len(b) != 0 {
		v, ok := m.table[string(b)]
		if exist && !ok {
			b = prefixB(b)
			continue
		}
		v = f(b, v)
		if v == nil {
			delete(m.table, string(b))
		} else {
			m.table[string(b)] = v
		}
		if all {
			b = prefixB(b)
		} else {
			break
		}
	}
}

func (m *threadUnsafeMatcher) Update(s string, f func(interface{}) interface{}, exist bool) {
	m.updateB(cpy(m.b, s), func(_ []byte, v interface{}) interface{} {
		return f(v)
	}, exist, false)
}

func (m *threadUnsafeMatcher) UpdateRaw(key []Component, f func(interface{}) interface{}, exist bool) {
	m.updateB(escape(m.b, key), func(_ []byte, v interface{}) interface{} {
		return f(v)
	}, exist, false)
}

func (m *threadUnsafeMatcher) UpdateAll(s string, f func([]byte, interface{}) interface{}, exist bool) {
	m.updateB(cpy(m.b, s), f, exist, true)
}

func (m *threadUnsafeMatcher) UpdateAllRaw(key []Component, f func([]byte, interface{}) interface{}, exist bool) {
	m.updateB(escape(m.b, key), f, exist, true)
}

func (m *threadUnsafeMatcher) matchB(b []byte, f func(interface{}), exist bool) {
	for len(b) != 0 {
		v, ok := m.table[string(b)]
		if exist && !ok {
			b = prefixB(b)
			continue
		}
		f(v)
		break
	}
}

func (m *threadUnsafeMatcher) Match(s string, f func(interface{}), exist bool) {
	m.matchB(cpy(m.b, s), f, exist)
}

func (m *threadUnsafeMatcher) MatchRaw(key []Component, f func(interface{}), exist bool) {
	m.matchB(escape(m.b, key), f, exist)
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
