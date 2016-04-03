package lpm

import "testing"

func TestMatcher(t *testing.T) {
	m := NewThreadSafe()

	update := func(key string, f func(interface{}) interface{}, exist bool) {
		m.Update(NewComponents(key), f, exist)
	}

	updateAll := func(key string, f func([]Component, interface{}) interface{}, exist bool) {
		m.UpdateAll(NewComponents(key), f, exist)
	}

	match := func(key string, f func(interface{}), exist bool) {
		m.Match(NewComponents(key), f, exist)
	}

	for _, test := range []struct {
		key   string
		value int
	}{
		{"/1", 1},
		{"/1/2", 12},
		{"/1/2/3", 123},
		{"/1/2/4", 124},
		{"/1/2/4/5", 1245},
	} {
		// add
		update(test.key, func(interface{}) interface{} { return test.value }, false)
	}

	for _, test := range []struct {
		in   string
		want interface{}
	}{
		{"/2", nil},
		{"/1/2/3/4", 123},
	} {
		match(test.in, func(got interface{}) {
			if got != test.want {
				t.Fatalf("Match(%v) == %v, got %v", test.in, test.want, got)
			}
		}, true)
	}

	update("/1/2/3", func(interface{}) interface{} { return nil }, false)
	for _, test := range []struct {
		in   string
		want interface{}
	}{
		{"/2", nil},
		{"/1/2/3/4", 12},
		{"/1/2/3", 12},
	} {
		match(test.in, func(got interface{}) {
			if got != test.want {
				t.Fatalf("Match(%v) == %v, got %v", test.in, test.want, got)
			}
		}, true)
	}

	update("/1/2/5", func(interface{}) interface{} { return 125 }, true)
	for _, test := range []struct {
		in   string
		want interface{}
	}{
		{"/2", nil},
		{"/1/2/3/4", 125},
		{"/1/2", 125},
	} {
		match(test.in, func(got interface{}) {
			if got != test.want {
				t.Fatalf("Match(%v) == %v, got %v", test.in, test.want, got)
			}
		}, true)
	}

	updateAll("/1/2/4/5", func(key []Component, i interface{}) interface{} {
		if len(key)%2 == 0 {
			return 2
		}
		return 1
	}, true)
	for _, test := range []struct {
		in   string
		want interface{}
	}{
		{"/2", nil},
		{"/1/2/4/5", 2},
		{"/1/2/4", 1},
		{"/1/2/3", 2},
		{"/1/2", 2},
		{"/1", 1},
	} {
		match(test.in, func(got interface{}) {
			if got != test.want {
				t.Fatalf("Match(%v) == %v, got %v", test.in, test.want, got)
			}
		}, true)
	}

	var count int
	m.Visit(func(_ []Component, v interface{}) interface{} {
		count++
		return v
	})
	if count != 4 {
		t.Fatalf("expect entry count to be 4, got %v", count)
	}
}
