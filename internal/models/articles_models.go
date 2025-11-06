package models

import "time"

// Article represents structure of existing article
//
//	@Description	Article structure
type Article struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	Body      string    `json:"body"`
	CreatedAt time.Time `json:"created_at"`
}
