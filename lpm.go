// Package lpm implements thread-safe longest prefix match (LPM)
package lpm

type Matcher interface {
	Update(string, func(interface{}) interface{}, bool)
	UpdateAll(string, func(string, interface{}) interface{})
	Match(string, func(interface{}))
	Visit(func(string, interface{}) interface{})
}

func New() Matcher {
	return newThreadUnsafeMatcher()
}

func NewThreadSafe() Matcher {
	return newThreadSafeMatcher()
}
