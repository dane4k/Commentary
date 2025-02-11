package imrepo

import (
	"Commentary/internal/entity"
	"Commentary/internal/graph/model"
	"errors"
	"github.com/sirupsen/logrus"
	"time"
)

func (imr *InMemoryRepo) AddComment(input model.CreateCommentInput) (*entity.Comment, error) {
	logrus.Debug("adding comment")

	if len(input.Content) > 2000 {
		return nil, errors.New("comment content exceeds the limitation of 2000 symbols")
	}

	imr.mu.Lock()
	defer imr.mu.Unlock()

	if _, ok := imr.users[input.AuthorID]; !ok {
		logrus.Error(ErrUserNotFound)
		return nil, ErrUserNotFound
	}
	if _, ok := imr.Posts[input.PostID]; !ok {
		logrus.Error(ErrPostNotFound)
		return nil, ErrPostNotFound
	}

	id := len(imr.comments) + 1
	comment := &entity.Comment{
		ID:       id,
		PostID:   input.PostID,
		AuthorID: input.AuthorID,
		Content:  input.Content,
		Created:  time.Now(),
		ParentID: input.Parent,
	}

	logrus.Debug("added comment")
	imr.comments[id] = comment
	return comment, nil
}

func (imr *InMemoryRepo) GetComment(id int) (*entity.Comment, error) {
	logrus.Debug("getting comment")

	imr.mu.RLock()
	defer imr.mu.RUnlock()

	comment, ok := imr.comments[id]
	if !ok {
		logrus.Error(ErrCommentNotFound)
		return nil, ErrCommentNotFound
	}
	logrus.Debug("got comment")
	return comment, nil
}

func (imr *InMemoryRepo) GetCommentsByPostID(postID int) ([]*entity.Comment, error) {
	logrus.Debug("getting comments by post id")

	imr.mu.RLock()
	defer imr.mu.RUnlock()

	if _, ok := imr.Posts[postID]; !ok {
		logrus.Error(ErrPostNotFound)
		return nil, ErrPostNotFound
	}

	comments := make([]*entity.Comment, 0, len(imr.comments))
	for _, comment := range imr.comments {
		if comment.PostID == postID {
			comments = append(comments, comment)
		}
	}
	logrus.Debug("got comments by post id")

	return comments, nil
}
