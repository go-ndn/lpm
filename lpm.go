// Package lpm implements thread-safe longest prefix match (LPM)
package lpm

import "fmt"

type Key string

func (this Key) String() string {
	return string(this)
}

type Matcher interface {
	Add(fmt.Stringer, interface{})
	Remove(fmt.Stringer)
	Update(fmt.Stringer, func(interface{}) interface{}, bool)
	UpdateAll(fmt.Stringer, func(string, interface{}) interface{})
	Match(fmt.Stringer) interface{}
	Visit(func(string, interface{}) interface{})
}

func New() Matcher {
	return newThreadUnsafeMatcher()
}

func NewThreadSafe() Matcher {
	return newThreadSafeMatcher()
}
