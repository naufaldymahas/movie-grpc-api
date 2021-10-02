package entity

import "time"

type MovieLog struct {
	ID        int64
	CreatedAt time.Time
	EventType string
	Params    string
}
