package pgdb

import (
	"errors"
	"fmt"
)

type RepositoryError struct {
	Operation string
	Content   string
	Err       error
}

func (re *RepositoryError) Error() string {
	return fmt.Sprintf("repository error occured while %s, %s: %v", re.Operation, re.Content, re.Err)
}

var (
	ErrAddingComment    = errors.New("error adding comment")
	ErrGeneratingSQL    = errors.New("failed to generate SQL query")
	ErrGettingComments  = errors.New("error getting comments")
	ErrGettingPost      = errors.New("error getting post")
	ErrAddingPost       = errors.New("error adding post")
	ErrTogglingComments = errors.New("error toggling comments")
	ErrGettingPosts     = errors.New("error getting posts")
	ErrGettingUser      = errors.New("error getting user")
	ErrGettingUsers     = errors.New("error getting users")
	ErrGettingComment   = errors.New("error getting comment")
	ErrAddingUser       = errors.New("error adding user")
)
