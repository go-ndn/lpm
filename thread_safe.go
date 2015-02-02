package lpm

import (
	"fmt"
	"sync"
)

type threadSafeMatcher struct {
	m threadUnsafeMatcher
	sync.RWMutex
}

func newThreadSafeMatcher() *threadSafeMatcher {
	return &threadSafeMatcher{m: *newThreadUnsafeMatcher()}
}

func (this *threadSafeMatcher) Add(cs fmt.Stringer, i interface{}) {
	this.Lock()
	this.m.Add(cs, i)
	this.Unlock()
}

func (this *threadSafeMatcher) Remove(cs fmt.Stringer) {
	this.Lock()
	this.m.Remove(cs)
	this.Unlock()
}

func (this *threadSafeMatcher) Update(cs fmt.Stringer, f func(interface{}) interface{}, isPrefix bool) {
	this.Lock()
	this.m.Update(cs, f, isPrefix)
	this.Unlock()
}

func (this *threadSafeMatcher) UpdateAll(cs fmt.Stringer, f func(string, interface{}) interface{}) {
	this.Lock()
	this.m.UpdateAll(cs, f)
	this.Unlock()
}

func (this *threadSafeMatcher) Match(cs fmt.Stringer) interface{} {
	this.RLock()
	defer this.RUnlock()
	return this.m.Match(cs)
}

func (this *threadSafeMatcher) Visit(f func(string, interface{}) interface{}) {
	this.Lock()
	this.m.Visit(f)
	this.Unlock()
}
