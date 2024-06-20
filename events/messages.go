package events

import "time"

type Message interface {
	Type() string
}

type CreatedFeedMessage struct {
	CreatedAt   time.Time `json:"created_at"`
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
}

func (m CreatedFeedMessage) Type() string {
	return "created_feed"
}
