package imrepo

import "errors"

var (
	ErrUserAlreadyExists = errors.New("user with this nickname already exists")
	ErrUserNotFound      = errors.New("user not found")
	ErrCommentNotFound   = errors.New("comment not found")
	ErrPostNotFound      = errors.New("post not found")
)
