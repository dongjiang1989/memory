package PUtil

import (
	"time"
)

type SimpleItem struct {
	Value      interface{}
	Expiration *time.Time
}

// returns boolean value whether this item is expired or not.
func (si *SimpleItem) IsExpired(now *time.Time) bool {
	if si.Expiration == nil {
		return false
	}
	if now == nil {
		t := time.Now()
		now = &t
	}
	return si.Expiration.Before(*now)
}

func (it *SimpleItem) Expire() *time.Time {
	return it.Expiration
}
