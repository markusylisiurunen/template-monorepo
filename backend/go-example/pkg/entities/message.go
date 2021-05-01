package entities

import "time"

type Message struct {
	ID        int64
	CreatedAt time.Time
	Text      string
}
