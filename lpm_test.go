package lpm

import (
	"strings"
	"testing"
)

type test struct {
	in   string
	want interface{}
}

func TestLPM(t *testing.T) {
	m := NewThreadSafe()
	m.Add(Key("1"), 1)
	m.Add(Key("1/2"), 12)
	m.Add(Key("1/2/3"), 123)
	m.Add(Key("1/2/4"), 124)
	m.Add(Key("1/2/4/5"), 1245)

	for _, test := range []test{
		{"2", nil},
		{"1/2/3/4", 123},
	} {
		got := m.Match(Key(test.in))
		if got != test.want {
			t.Fatalf("Match(%s) == %v, got %v", test.in, test.want, got)
		}
	}

	m.Remove(Key("1/2/3"))
	for _, test := range []test{
		{"1/2/3", 12},
	} {
		got := m.Match(Key(test.in))
		if got != test.want {
			t.Fatalf("Match(%s) == %v, got %v", test.in, test.want, got)
		}
	}

	m.Update(Key("1/2/5"), func(interface{}) interface{} {
		return 125
	}, true)
	for _, test := range []test{
		{"1/2", 125},
	} {
		got := m.Match(Key(test.in))
		if got != test.want {
			t.Fatalf("Match(%s) == %v, got %v", test.in, test.want, got)
		}
	}

	m.UpdateAll(Key("1/2/4/5"), func(s string, i interface{}) interface{} {
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
		got := m.Match(Key(test.in))
		if got != test.want {
			t.Fatalf("Match(%s) == %v, got %v", test.in, test.want, got)
		}
	}
}
