package matcher

import (
	"testing"

	"github.com/go-ndn/lpm"
)

func init() {
	nodeValEmpty = func(t Type) bool {
		return t == 0
	}
}

func TestMatcher(t *testing.T) {
	var m TypeMatcher

	update := func(key string, f func(Type) Type, exist bool) {
		m.Update(lpm.NewComponents(key), f, exist)
	}

	updateAll := func(key string, f func([]lpm.Component, Type) Type, exist bool) {
		m.UpdateAll(lpm.NewComponents(key), f, exist)
	}

	match := func(key string, f func(Type), exist bool) {
		m.Match(lpm.NewComponents(key), f, exist)
	}

	for _, test := range []struct {
		key   string
		value Type
	}{
		{"/1", 1},
		{"/1/2", 12},
		{"/1/2/3", 123},
		{"/1/2/4", 124},
		{"/1/2/4/5", 1245},
	} {
		// add
		update(test.key, func(Type) Type { return test.value }, false)
	}

	for _, test := range []struct {
		in   string
		want Type
	}{
		{"/2", 0},
		{"/1/2/3/4", 123},
	} {
		match(test.in, func(got Type) {
			if got != test.want {
				t.Fatalf("Match(%v) == %v, got %v", test.in, test.want, got)
			}
		}, true)
	}

	update("/1/2/3", func(Type) Type { return 0 }, false)
	for _, test := range []struct {
		in   string
		want Type
	}{
		{"/2", 0},
		{"/1/2/3/4", 12},
		{"/1/2/3", 12},
	} {
		match(test.in, func(got Type) {
			if got != test.want {
				t.Fatalf("Match(%v) == %v, got %v", test.in, test.want, got)
			}
		}, true)
	}

	update("/1/2/5", func(Type) Type { return 125 }, true)
	for _, test := range []struct {
		in   string
		want Type
	}{
		{"/2", 0},
		{"/1/2/3/4", 125},
		{"/1/2", 125},
	} {
		match(test.in, func(got Type) {
			if got != test.want {
				t.Fatalf("Match(%v) == %v, got %v", test.in, test.want, got)
			}
		}, true)
	}

	updateAll("/1/2/4/5", func(key []lpm.Component, i Type) Type {
		if len(key)%2 == 0 {
			return 2
		}
		return 1
	}, true)
	for _, test := range []struct {
		in   string
		want Type
	}{
		{"/2", 0},
		{"/1/2/4/5", 2},
		{"/1/2/4", 1},
		{"/1/2/3", 2},
		{"/1/2", 2},
		{"/1", 1},
	} {
		match(test.in, func(got Type) {
			if got != test.want {
				t.Fatalf("Match(%v) == %v, got %v", test.in, test.want, got)
			}
		}, true)
	}

	var count int
	m.Visit(func(_ []lpm.Component, v Type) Type {
		count++
		return v
	})
	if count != 4 {
		t.Fatalf("expect entry count to be 4, got %v", count)
	}
}
