package lpm

import (
	"testing"
)

func TestLPM(t *testing.T) {
	m := New()
	m.Set([]Component{"1", "2", "3"}, "hello")
	m.Set([]Component{"1", "2"}, "world")
	m.Set([]Component{"1", "2", "4"}, 124)
	m.Set([]Component{"1", "2", "5", "6"}, 1256)

	hello := m.Match([]Component{"1", "2", "3", "4"})
	if hello == nil || hello.(string) != "hello" {
		t.Fatal("not hello")
	}
	world := m.Match([]Component{"1", "2"})
	if world == nil || world.(string) != "world" {
		t.Fatal("not world")
	}

	m.Update([]Component{"1", "2", "3", "4"}, func(interface{}) interface{} { return nil })
	null := m.Match([]Component{"1", "2", "3", "4"})
	if null != nil {
		t.Fatal("should be nil")
	}

	n := m.RMatch([]Component{"1", "2", "5"})
	if n == nil || n.(int) != 1256 {
		t.Fatal("should be 1256")
	}

	t.Logf("%v\n", m.List())
}
