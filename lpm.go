// Package lpm implements thread-safe longest prefix match (LPM)
package lpm

import (
	"sync"
)

type Component string

type node struct {
	table map[Component]*node
	entry interface{}
}

func newNode() *node {
	return &node{
		table: make(map[Component]*node),
	}
}

type Matcher struct {
	root *node
	m    sync.RWMutex
}

func New() *Matcher {
	return &Matcher{
		root: newNode(),
	}
}

func update(n *node, cs []Component, f func(interface{}) interface{}, isPrefix bool) {
	if len(cs) == 0 || (isPrefix && len(n.table) == 0) {
		n.entry = f(n.entry)
		return
	}
	first, rest := cs[0], cs[1:]
	c, ok := n.table[first]
	if !ok {
		if isPrefix {
			return
		}
		c = newNode()
		n.table[first] = c
	}
	update(c, rest, f, isPrefix)
	if c.entry == nil && len(c.table) == 0 {
		delete(n.table, first)
	}
}

func (this *Matcher) Add(cs []Component, i interface{}) {
	this.m.Lock()
	update(this.root, cs, func(interface{}) interface{} { return i }, false)
	this.m.Unlock()
}

func (this *Matcher) Remove(cs []Component) {
	this.m.Lock()
	update(this.root, cs, func(interface{}) interface{} { return nil }, true)
	this.m.Unlock()
}

// Update provides atomic RW on longest prefix's entry with full name
func (this *Matcher) Update(cs []Component, f func(interface{}) interface{}, isPrefix bool) {
	this.m.Lock()
	update(this.root, cs, f, isPrefix)
	this.m.Unlock()
}

func findEntry(n *node) interface{} {
	if n.entry != nil {
		return n.entry
	}
	for _, c := range n.table {
		ce := findEntry(c)
		if ce != nil {
			return ce
		}
	}
	return nil
}

func match(n *node, cs []Component, isPrefix bool) interface{} {
	if len(cs) == 0 || (isPrefix && len(n.table) == 0) {
		if !isPrefix {
			return findEntry(n)
		}
		return n.entry
	}
	first, rest := cs[0], cs[1:]
	c, ok := n.table[first]
	if !ok {
		return nil
	}
	return match(c, rest, isPrefix)
}

// Match finds longest prefix's entry with full name
func (this *Matcher) Match(cs []Component) interface{} {
	this.m.RLock()
	defer this.m.RUnlock()
	return match(this.root, cs, true)
}

// Reverse Match finds full name's entry with longest prefix
func (this *Matcher) RMatch(cs []Component) interface{} {
	this.m.RLock()
	defer this.m.RUnlock()
	return match(this.root, cs, false)
}

func list(n *node, prefix string) (es []string) {
	if n.entry != nil {
		es = append(es, prefix)
	}
	for part, c := range n.table {
		ces := list(c, prefix+"/"+string(part))
		es = append(es, ces...)
	}
	return
}

func (this *Matcher) List() []string {
	this.m.RLock()
	defer this.m.RUnlock()
	return list(this.root, "")
}
