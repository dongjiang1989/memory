package PUtil

import (
	"container/list"
)

type ArcList struct {
	l    *list.List
	keys map[interface{}]*list.Element
}

func NewARCList() *ArcList {
	return &ArcList{
		l:    list.New(),
		keys: make(map[interface{}]*list.Element),
	}
}

// has key func
// return bool
func (al *ArcList) Has(key interface{}) bool {
	_, ok := al.keys[key]
	return ok
}

// Lookup func : search list.element for key
func (al *ArcList) Lookup(key interface{}) *list.Element {
	elt := al.keys[key]
	return elt
}

// Move item to front
func (al *ArcList) MoveToFront(elt *list.Element) {
	al.l.MoveToFront(elt)
}

// push item to front
func (al *ArcList) PushFront(key interface{}) {
	elt := al.l.PushFront(key)
	al.keys[key] = elt
}

// delete item
func (al *ArcList) Remove(key interface{}, elt *list.Element) {
	if al.Has(key) {
		delete(al.keys, key)
	}
	al.l.Remove(elt)
}

// delete last
func (al *ArcList) RemoveTail() interface{} {
	elt := al.l.Back()
	al.l.Remove(elt)

	key := elt.Value
	if al.Has(key) {
		delete(al.keys, key)
	}

	return key
}

// list len
func (al *ArcList) Len() int {
	return al.l.Len()
}
