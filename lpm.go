// Package lpm implements thread-safe longest prefix match (LPM)
package lpm

import "net/url"

type Component []byte

func (c Component) String() string {
	return url.QueryEscape(string(c))
}

type Matcher interface {
	Update(string, func(interface{}) interface{}, bool)
	UpdateRaw([]Component, func(interface{}) interface{}, bool)

	UpdateAll(string, func([]byte, interface{}) interface{}, bool)
	UpdateAllRaw([]Component, func([]byte, interface{}) interface{}, bool)

	Match(string, func(interface{}), bool)
	MatchRaw([]Component, func(interface{}), bool)

	Visit(func(string, interface{}) interface{})
}
