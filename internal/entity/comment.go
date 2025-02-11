package entity

import "time"

type Comment struct {
	ID       int       `json:"id" db:"id"`
	PostID   int       `json:"post" db:"post_id"`
	AuthorID int       `json:"author" db:"author_id"`
	Content  string    `json:"content" db:"content"`
	Created  time.Time `json:"created" db:"created"`
	ParentID *int      `json:"parent,omitempty" db:"parent_id"`
}
