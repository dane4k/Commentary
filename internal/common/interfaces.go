package common

import (
	"Commentary/internal/graph/model"
	"context"
)

type PostService interface {
	SetCommentService(commentService CommentService)
	CreatePost(ctx context.Context, input model.CreatePostInput) (*model.Post, error)
	GetPost(ctx context.Context, id int, limit, offset *int) (*model.Post, error)
	ToggleComments(ctx context.Context, postID int) (*model.Post, error)
	GetPosts(ctx context.Context, limit, offset *int) ([]*model.Post, error)
}

type CommentService interface {
	GetComments(ctx context.Context, postID int, limit, offset *int) ([]*model.Comment, error)
	CreateComment(ctx context.Context, input model.CreateCommentInput) (*model.Comment, error)
	GetCommentsForPosts(ctx context.Context, postIDs []int) ([]*model.Comment, error)
}

type UserService interface {
	CreateUser(ctx context.Context, username string) (*model.User, error)
}
