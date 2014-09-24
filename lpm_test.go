package lpm

import (
	"testing"
)

type key string

func (this key) String() string {
	return string(this)
}

func TestLPM(t *testing.T) {
	m := New()
	m.Add(key("/1/2/3"), "hello")
	m.Add(key("/1/2"), "world")
	m.Add(key("/1/2/4"), 124)
	m.Add(key("/1/2/5/6"), 1256)

	hello := m.Match(key("/1/2/3/4"))
	if hello == nil || hello.(string) != "hello" {
		t.Fatal("not hello")
	}
	world := m.Match(key("/1/2"))
	if world == nil || world.(string) != "world" {
		t.Fatal("not world")
	}

	m.Remove(key("/1/2/3/4"))
	null := m.Match(key("/1/2/3/4"))
	if null != nil {
		t.Fatal("should be nil")
	}

	n := m.RMatch(key("/1/2/5"))
	if n == nil || n.(int) != 1256 {
		t.Fatal("should be 1256")
	}

	t.Logf("%v\n", m.List())
}
