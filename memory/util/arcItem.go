package PUtil

import (
	"time"
)

type ArcItem struct {
	Key        interface{}
	Value      interface{}
	Expiration *time.Time
}

// returns boolean value whether this item is expired or not.
func (it *ArcItem) IsExpired(now *time.Time) bool {
	if it.Expiration == nil {
		return false
	}
	if now == nil {
		t := time.Now()
		now = &t
	}
	return it.Expiration.Before(*now)
}

func (it *ArcItem) Expire() *time.Time {
	return it.Expiration
}
