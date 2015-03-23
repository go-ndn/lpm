package lpm

import (
	"testing"
)

func TestLPM(t *testing.T) {
	m := NewThreadSafe()
	m.Add(Key("/1/2/3"), 123)
	m.Add(Key("/1/2"), 12)
	m.Add(Key("/1/2/4"), 124)
	m.Add(Key("/1/2/5/6"), 1256)

	num := m.Match(Key("/1/2/3/4"))
	if num == nil || num.(int) != 123 {
		t.Fatal("want 123")
	}

	m.Remove(Key("/1/2/3"))
	num = m.Match(Key("/1/2/3"))
	if num == nil || num.(int) != 12 {
		t.Fatal("want 12")
	}

	m.Update(Key("/1/2/6"), func(interface{}) interface{} {
		return 126
	}, true)
	num = m.Match(Key("/1/2"))
	if num == nil || num.(int) != 126 {
		t.Fatal("want 126")
	}

	m.UpdateAll(Key("/1/2/3"), func(string, interface{}) interface{} {
		return 1
	})
	num = m.Match(Key("/1/2"))
	if num == nil || num.(int) != 1 {
		t.Fatal("want 1")
	}
}
