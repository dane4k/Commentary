package entity

import "time"

type Post struct {
	ID          int       `json:"id" db:"id"`
	AuthorID    int       `json:"author" db:"author_id"`
	Title       string    `json:"title" db:"title"`
	Content     string    `json:"content" db:"content"`
	Created     time.Time `json:"created" db:"created"`
	Commentable bool      `json:"commentable" db:"commentable"`
}
