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

func set(n *node, cs []Component, i interface{}) {
	if len(cs) == 0 {
		n.entry = i
		return
	}
	first, rest := cs[0], cs[1:]
	c, ok := n.table[first]
	if !ok {
		c = newNode()
		n.table[first] = c
	}
	set(c, rest, i)
}

// Set changes full name's entry
func (this *Matcher) Set(cs []Component, i interface{}) {
	this.m.Lock()
	set(this.root, cs, i)
	this.m.Unlock()
}

func update(n *node, cs []Component, f func(interface{}) interface{}) {
	if len(cs) == 0 || len(n.table) == 0 {
		n.entry = f(n.entry)
		return
	}
	first, rest := cs[0], cs[1:]
	c, ok := n.table[first]
	if !ok {
		return
	}
	update(c, rest, f)
	if c.entry == nil && len(c.table) == 0 {
		delete(n.table, first)
	}
}

// Update changes longest prefix's entry with full name
func (this *Matcher) Update(cs []Component, f func(interface{}) interface{}) {
	this.m.Lock()
	update(this.root, cs, f)
	this.m.Unlock()
}

func match(n *node, cs []Component) interface{} {
	if len(cs) == 0 || len(n.table) == 0 {
		return n.entry
	}
	first, rest := cs[0], cs[1:]
	c, ok := n.table[first]
	if !ok {
		return nil
	}
	return match(c, rest)
}

// Match finds longest prefix's entry with full name
func (this *Matcher) Match(cs []Component) interface{} {
	this.m.RLock()
	defer this.m.RUnlock()
	return match(this.root, cs)
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

func rmatch(n *node, cs []Component) interface{} {
	if len(cs) == 0 {
		return findEntry(n)
	}
	first, rest := cs[0], cs[1:]
	c, ok := n.table[first]
	if !ok {
		return nil
	}
	return rmatch(c, rest)
}

// Reverse Match finds full name's entry with longest prefix
func (this *Matcher) RMatch(cs []Component) interface{} {
	this.m.RLock()
	defer this.m.RUnlock()
	return rmatch(this.root, cs)
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
