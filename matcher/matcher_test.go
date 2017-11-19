package matcher

import (
	"testing"

	"github.com/go-ndn/lpm"
)

func TestMatcher(t *testing.T) {
	var m TypeMatcher

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
		m.Update(lpm.NewComponents(test.key), test.value)
	}

	for _, test := range []struct {
		in   string
		want Type
	}{
		{"/2", 0},
		{"/1/2/3/4", 123},
	} {
		got, _ := m.Match(lpm.NewComponents(test.in))
		if got != test.want {
			t.Fatalf("Match(%v) == %v, got %v", test.in, test.want, got)
		}
	}

	m.Delete(lpm.NewComponents("/1/2/3"))
	for _, test := range []struct {
		in   string
		want Type
	}{
		{"/2", 0},
		{"/1/2/3/4", 12},
		{"/1/2/3", 12},
	} {
		got, _ := m.Match(lpm.NewComponents(test.in))
		if got != test.want {
			t.Fatalf("Match(%v) == %v, got %v", test.in, test.want, got)
		}
	}

	m.Update(lpm.NewComponents("/1/2"), 125)
	for _, test := range []struct {
		in   string
		want Type
	}{
		{"/2", 0},
		{"/1/2/3/4", 125},
		{"/1/2", 125},
	} {
		got, _ := m.Match(lpm.NewComponents(test.in))
		if got != test.want {
			t.Fatalf("Match(%v) == %v, got %v", test.in, test.want, got)
		}
	}

	m.UpdateAll(lpm.NewComponents("/1/2/4/5"), func(key []lpm.Component, i Type) (Type, bool) {
		if len(key)%2 == 0 {
			return 2, false
		}
		return 1, false
	})
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
		got, _ := m.Match(lpm.NewComponents(test.in))
		if got != test.want {
			t.Fatalf("Match(%v) == %v, got %v", test.in, test.want, got)
		}
	}

	var count int
	m.Visit(func(_ []lpm.Component, v Type) (Type, bool) {
		count++
		return v, false
	})
	if want := 4; count != want {
		t.Fatalf("expect entry count to be %d, got %d", want, count)
	}
}
