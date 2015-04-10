package lpm

import (
	"strings"
	"testing"
)

type test struct {
	in   string
	want interface{}
}

func add(m Matcher, s string, v int) {
	m.Update(s, func(interface{}) interface{} { return v }, false)
}

func remove(m Matcher, s string) {
	m.Update(s, func(interface{}) interface{} { return nil }, false)
}

func match(m Matcher, s string) (r interface{}) {
	m.Match(s, func(v interface{}) { r = v })
	return
}

func TestLPM(t *testing.T) {
	m := NewThreadSafe()
	add(m, "1", 1)
	add(m, "1/2", 12)
	add(m, "1/2/3", 123)
	add(m, "1/2/4", 124)
	add(m, "1/2/4/5", 1245)

	for _, test := range []test{
		{"2", nil},
		{"1/2/3/4", 123},
	} {
		got := match(m, test.in)
		if got != test.want {
			t.Fatalf("Match(%s) == %v, got %v", test.in, test.want, got)
		}
	}

	remove(m, "1/2/3")
	for _, test := range []test{
		{"1/2/3", 12},
	} {
		got := match(m, test.in)
		if got != test.want {
			t.Fatalf("Match(%s) == %v, got %v", test.in, test.want, got)
		}
	}

	m.Update("1/2/5", func(interface{}) interface{} {
		return 125
	}, true)
	for _, test := range []test{
		{"1/2", 125},
	} {
		got := match(m, test.in)
		if got != test.want {
			t.Fatalf("Match(%s) == %v, got %v", test.in, test.want, got)
		}
	}

	m.UpdateAll("1/2/4/5", func(s string, i interface{}) interface{} {
		if strings.Count(s, "/")%2 == 0 {
			return 2
		}
		return 1
	})
	for _, test := range []test{
		{"1/2/4/5", 1},
		{"1/2/4", 2},
		{"1/2/3", 1},
		{"1/2", 1},
		{"1", 2},
	} {
		got := match(m, test.in)
		if got != test.want {
			t.Fatalf("Match(%s) == %v, got %v", test.in, test.want, got)
		}
	}
}
