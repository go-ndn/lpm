package lpm

import (
	"fmt"
	"strings"
)

type threadUnsafeMatcher struct {
	table map[string]interface{}
}

func newThreadUnsafeMatcher() *threadUnsafeMatcher {
	return &threadUnsafeMatcher{table: make(map[string]interface{})}
}

func (this *threadUnsafeMatcher) Add(cs fmt.Stringer, i interface{}) {
	this.Update(cs, func(interface{}) interface{} { return i }, false)
}

func (this *threadUnsafeMatcher) Remove(cs fmt.Stringer) {
	this.Update(cs, func(interface{}) interface{} { return nil }, false)
}

func (this *threadUnsafeMatcher) prefix(s string, longest bool) (p []string) {
	for {
		if _, ok := this.table[s]; ok {
			p = append(p, s)
			if longest {
				break
			}
		}
		idx := strings.LastIndex(s, "/")
		if idx == -1 {
			break
		}
		s = s[:idx]
	}
	return
}

func (this *threadUnsafeMatcher) Update(cs fmt.Stringer, f func(interface{}) interface{}, isPrefix bool) {
	s := cs.String()
	if isPrefix {
		p := this.prefix(s, true)
		if len(p) == 0 {
			return
		}
		s = p[0]
	}
	this.table[s] = f(this.table[s])
	if this.table[s] == nil {
		delete(this.table, s)
	}
}

func (this *threadUnsafeMatcher) UpdateAll(cs fmt.Stringer, f func(string, interface{}) interface{}) {
	for _, s := range this.prefix(cs.String(), false) {
		this.table[s] = f(s, this.table[s])
		if this.table[s] == nil {
			delete(this.table, s)
		}
	}
}

func (this *threadUnsafeMatcher) Match(cs fmt.Stringer) interface{} {
	p := this.prefix(cs.String(), true)
	if len(p) == 0 {
		return nil
	}
	return this.table[p[0]]
}
