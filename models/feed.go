package models

import "time"

type Feed struct {
	CreatedAt   time.Time `json:"created_at"`
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
}
