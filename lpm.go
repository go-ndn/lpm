// Package lpm implements thread-safe longest prefix match (LPM).
package lpm

import "net/url"

// Component is an arbitrary byte sequence.
type Component []byte

func (c Component) String() string {
	return url.QueryEscape(string(c))
}

// Matcher performs longest prefix match on percent-encoded string
// (https://en.wikipedia.org/wiki/Percent-encoding).
// It accepts both string key and raw component key.
//
// If func returns nil, the entry will be removed.
// If bool is false, exact matching will be performed instead.
type Matcher interface {
	Update(string, func(interface{}) interface{}, bool)
	UpdateRaw([]Component, func(interface{}) interface{}, bool)

	UpdateAll(string, func([]byte, interface{}) interface{}, bool)
	UpdateAllRaw([]Component, func([]byte, interface{}) interface{}, bool)

	Match(string, func(interface{}), bool)
	MatchRaw([]Component, func(interface{}), bool)

	Visit(func(string, interface{}) interface{})
}
