package lpm

import (
	"testing"
)

func TestLPM(t *testing.T) {
	m := NewThreadSafe()
	m.Add(Key("/1/2/3"), "hello")
	m.Add(Key("/1/2"), "world")
	m.Add(Key("/1/2/4"), 124)
	m.Add(Key("/1/2/5/6"), 1256)

	hello := m.Match(Key("/1/2/3/4"))
	if hello == nil || hello.(string) != "hello" {
		t.Fatal("not hello")
	}
	world := m.Match(Key("/1/2"))
	if world == nil || world.(string) != "world" {
		t.Fatal("not world")
	}

	m.Remove(Key("/1/2/3"))
	world = m.Match(Key("/1/2/3"))
	if world == nil || world.(string) != "world" {
		t.Fatal("should be world", world)
	}
}
