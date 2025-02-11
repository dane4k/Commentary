package model

import (
	"time"
)

type Comment struct {
	ID      int        `json:"id"`
	Post    *Post      `json:"post"`
	Author  *User      `json:"author"`
	Content string     `json:"content"`
	Created time.Time  `json:"created"`
	Parent  *Comment   `json:"parent,omitempty"`
	Replies []*Comment `json:"replies"`
}
