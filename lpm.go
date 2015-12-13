// Package lpm implements thread-safe longest prefix match (LPM)
package lpm

type Matcher interface {
	Update(string, func(interface{}) interface{}, bool)
	UpdateAll(string, func(string, interface{}) interface{}, bool)
	Match(string, func(interface{}), bool)
	Visit(func(string, interface{}) interface{})
}
