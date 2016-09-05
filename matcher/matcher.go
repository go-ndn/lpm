package matcher

import "github.com/go-ndn/lpm"

// Type is a placeholder.
type Type uint32

// TypeMatcher performs longest prefix match on components.
//
// If bool is false, exact match will be performed instead.
type TypeMatcher struct {
	node
}

var nodeValEmpty func(Type) bool

type node struct {
	val   Type
	table map[string]node
}

func (n *node) empty() bool {
	return nodeValEmpty(n.val) && len(n.table) == 0
}

func (n *node) update(key []lpm.Component, depth int, f func([]lpm.Component, Type) Type, exist, all bool) {
	try := func() {
		if !exist || !nodeValEmpty(n.val) {
			n.val = f(key[:depth], n.val)
		}
	}
	if len(key) == depth {
		try()
		return
	}

	if n.table == nil {
		if exist {
			try()
			return
		}
		n.table = make(map[string]node)
	}

	v, ok := n.table[string(key[depth])]
	if !ok {
		if exist {
			try()
			return
		}
	}

	if all {
		try()
	}

	v.update(key, depth+1, f, exist, all)
	if v.empty() {
		delete(n.table, string(key[depth]))
	} else {
		n.table[string(key[depth])] = v
	}
}

func (n *node) match(key []lpm.Component, depth int, f func(Type), exist bool) {
	try := func() {
		if !exist || !nodeValEmpty(n.val) {
			f(n.val)
		}
	}
	if len(key) == depth {
		try()
		return
	}

	if n.table == nil {
		if exist {
			try()
		}
		return
	}

	v, ok := n.table[string(key[depth])]
	if !ok {
		if exist {
			try()
		}
		return
	}

	v.match(key, depth+1, f, exist)
}

func (n *node) visit(key []lpm.Component, f func([]lpm.Component, Type) Type) {
	if !nodeValEmpty(n.val) {
		n.val = f(key, n.val)
	}
	for k, v := range n.table {
		v.visit(append(key, lpm.Component(k)), f)
		if v.empty() {
			delete(n.table, k)
		} else {
			n.table[k] = v
		}
	}
}

func (n *node) Update(key []lpm.Component, f func(Type) Type, exist bool) {
	n.update(key, 0, func(_ []lpm.Component, v Type) Type {
		return f(v)
	}, exist, false)
}

func (n *node) UpdateAll(key []lpm.Component, f func([]lpm.Component, Type) Type, exist bool) {
	n.update(key, 0, f, exist, true)
}

func (n *node) Match(key []lpm.Component, f func(Type), exist bool) {
	n.match(key, 0, f, exist)
}

func (n *node) Visit(f func([]lpm.Component, Type) Type) {
	key := make([]lpm.Component, 0, 16)
	n.visit(key, f)
}
