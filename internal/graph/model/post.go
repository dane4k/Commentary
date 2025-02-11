package model

import (
	"time"
)

type Post struct {
	ID          int        `json:"id"`
	Author      *User      `json:"author"`
	Title       string     `json:"title"`
	Content     string     `json:"content"`
	Created     time.Time  `json:"created"`
	Commentable bool       `json:"commentable"`
	Comments    []*Comment `json:"comments"`
}
