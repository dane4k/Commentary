package model

type CreateCommentInput struct {
	AuthorID int    `json:"authorID"`
	PostID   int    `json:"postID"`
	Content  string `json:"content"`
	Parent   *int   `json:"parent,omitempty"`
}

type CreatePostInput struct {
	AuthorID    int    `json:"authorID"`
	Title       string `json:"title"`
	Content     string `json:"content"`
	Commentable bool   `json:"commentable"`
}

type Mutation struct{}

type Query struct {
}

type Subscription struct {
}
