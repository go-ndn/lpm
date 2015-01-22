// Package lpm implements thread-safe longest prefix match (LPM)
package lpm

import (
	"fmt"
	"strings"
	"sync"
)

type Key string

func (this Key) String() string {
	return string(this)
}

type Matcher struct {
	table map[string]interface{}
	m     sync.RWMutex
}

func New() *Matcher {
	return &Matcher{
		table: make(map[string]interface{}),
	}
}

func (this *Matcher) Add(cs fmt.Stringer, i interface{}) {
	this.Update(cs, func(interface{}) interface{} { return i }, false)
}

func (this *Matcher) Remove(cs fmt.Stringer) {
	this.Update(cs, func(interface{}) interface{} { return nil }, false)
}

func (this *Matcher) prefix(s string, longest bool) (p []string) {
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

func (this *Matcher) Update(cs fmt.Stringer, f func(interface{}) interface{}, isPrefix bool) {
	this.m.Lock()
	defer this.m.Unlock()
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

func (this *Matcher) UpdateAll(cs fmt.Stringer, f func(string, interface{}) interface{}) {
	this.m.Lock()
	for _, s := range this.prefix(cs.String(), false) {
		this.table[s] = f(s, this.table[s])
		if this.table[s] == nil {
			delete(this.table, s)
		}
	}
	this.m.Unlock()
}

func (this *Matcher) Match(cs fmt.Stringer) interface{} {
	this.m.RLock()
	defer this.m.RUnlock()
	p := this.prefix(cs.String(), true)
	if len(p) == 0 {
		return nil
	}
	return this.table[p[0]]
}

func (this *Matcher) Visit(f func(string, interface{}) interface{}) {
	this.m.Lock()
	for e := range this.table {
		this.table[e] = f(e, this.table[e])
		if this.table[e] == nil {
			delete(this.table, e)
		}
	}
	this.m.Unlock()
}
