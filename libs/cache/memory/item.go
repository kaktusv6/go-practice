package memory

import "time"

type Item struct {
	Value      interface{}
	CreatedAd  time.Time
	Expiration int64
}
