package graph

import (
	"Commentary/internal/common"
	"Commentary/internal/pubsub"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	PostService    common.PostService
	CommentService common.CommentService
	UserService    common.UserService
	broker         *pubsub.Broker
}

func NewResolver(PostService common.PostService, CommentService common.CommentService,
	UserService common.UserService, broker *pubsub.Broker) *Resolver {
	return &Resolver{
		PostService:    PostService,
		CommentService: CommentService,
		UserService:    UserService,
		broker:         broker,
	}
}
