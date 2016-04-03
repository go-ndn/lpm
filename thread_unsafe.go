package lpm

type node struct {
	val   interface{}
	table map[string]*node
}

func (n *node) empty() bool {
	return n.val == nil && len(n.table) == 0
}

func (n *node) update(key []Component, depth int, f func([]Component, interface{}) interface{}, exist, all bool) {
	try := func() {
		if depth == 0 {
			return
		}
		if !exist || n.val != nil {
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
		n.table = make(map[string]*node)
	}

	v, ok := n.table[string(key[depth])]
	if !ok {
		if exist {
			try()
			return
		}
		v = &node{}
		n.table[string(key[depth])] = v
	}

	if all {
		try()
	}

	v.update(key, depth+1, f, exist, all)
	if v.empty() {
		delete(n.table, string(key[depth]))
	}
}

func (n *node) match(key []Component, depth int, f func(interface{}), exist bool) {
	try := func() {
		if depth == 0 {
			return
		}
		if !exist || n.val != nil {
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

func (n *node) visit(key []Component, f func([]Component, interface{}) interface{}) {
	if n.val != nil {
		n.val = f(key, n.val)
	}
	for k, v := range n.table {
		v.visit(append(key, Component(k)), f)
		if v.empty() {
			delete(n.table, k)
		}
	}
}

// New creates a new thread-unsafe matcher.
func New() Matcher {
	return &node{}
}

func (n *node) Update(key []Component, f func(interface{}) interface{}, exist bool) {
	n.update(key, 0, func(_ []Component, v interface{}) interface{} {
		return f(v)
	}, exist, false)
}

func (n *node) UpdateAll(key []Component, f func([]Component, interface{}) interface{}, exist bool) {
	n.update(key, 0, f, exist, true)
}

func (n *node) Match(key []Component, f func(interface{}), exist bool) {
	n.match(key, 0, f, exist)
}

func (n *node) Visit(f func([]Component, interface{}) interface{}) {
	key := make([]Component, 0, 16)
	n.visit(key, f)
}
