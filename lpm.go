// Package lpm implements longest prefix match (LPM).
package lpm

import (
	"net/url"
	"strings"
)

// Component is an arbitrary byte sequence.
type Component []byte

// NewComponents creates components from percent-encoded form.
//
// See https://en.wikipedia.org/wiki/Percent-encoding.
func NewComponents(s string) (cs []Component) {
	s = strings.Trim(s, "/")
	if s == "" {
		return
	}
	parts := strings.Split(s, "/")
	cs = make([]Component, len(parts))
	for i := range parts {
		parts[i], _ = url.QueryUnescape(parts[i])
		cs[i] = Component(parts[i])
	}
	return
}

func (c Component) String() string {
	return url.QueryEscape(string(c))
}

// Matcher performs longest prefix match on components.
//
// If func returns nil, the entry will be removed.
// If bool is false, exact match will be performed instead.
type Matcher interface {
	Update([]Component, func(interface{}) interface{}, bool)
	UpdateAll([]Component, func([]Component, interface{}) interface{}, bool)
	Match([]Component, func(interface{}), bool)
	Visit(func([]Component, interface{}) interface{})
}
