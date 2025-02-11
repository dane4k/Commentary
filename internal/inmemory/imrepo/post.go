package imrepo

import (
	"Commentary/internal/entity"
	"Commentary/internal/graph/model"
	"github.com/sirupsen/logrus"
	"sort"
	"time"
)

func (imr *InMemoryRepo) AddPost(input model.CreatePostInput) (*entity.Post, error) {
	logrus.Debug("adding post")

	imr.mu.Lock()
	defer imr.mu.Unlock()

	if _, ok := imr.users[input.AuthorID]; !ok {
		logrus.Error(ErrUserNotFound)
		return nil, ErrUserNotFound
	}

	postID := len(imr.Posts) + 1
	post := &entity.Post{
		ID:          postID,
		AuthorID:    input.AuthorID,
		Title:       input.Title,
		Content:     input.Content,
		Created:     time.Now(),
		Commentable: input.Commentable,
	}

	logrus.Debug("added post")

	imr.Posts[postID] = post

	return post, nil
}

func (imr *InMemoryRepo) GetPost(id int) (*entity.Post, error) {
	logrus.Debug("getting post")

	imr.mu.RLock()
	defer imr.mu.RUnlock()

	post, ok := imr.Posts[id]
	if !ok {
		logrus.Error(ErrPostNotFound)
		return nil, ErrPostNotFound
	}

	logrus.Debug("got post")

	return post, nil
}

func (imr *InMemoryRepo) ToggleComments(postID int) (*entity.Post, error) {
	logrus.Debug("toggling comments")

	imr.mu.Lock()
	defer imr.mu.Unlock()

	post, ok := imr.Posts[postID]
	if !ok {
		logrus.Error(ErrPostNotFound)
		return nil, ErrPostNotFound
	}

	imr.Posts[postID].Commentable = !imr.Posts[postID].Commentable
	logrus.Debug("toggled comments")

	return post, nil
}

func (imr *InMemoryRepo) GetPostsPag(limit, offset *int) ([]*entity.Post, error) {
	logrus.Debug("getting posts paginated")

	imr.mu.RLock()
	defer imr.mu.RUnlock()

	allPosts := make([]*entity.Post, 0, len(imr.Posts))
	for _, post := range imr.Posts {
		allPosts = append(allPosts, post)
	}

	sort.Slice(allPosts, func(i, j int) bool {
		return allPosts[i].Created.Before(allPosts[j].Created)
	})

	length := len(allPosts)
	startIdx, endIdx := 0, length
	if offset != nil {
		startIdx = *offset
		if startIdx < 0 || startIdx > length {
			startIdx = 0
		}
	}
	if limit != nil {
		endIdx = *limit + startIdx
		if endIdx < 0 || endIdx > length {
			endIdx = length
		}
	}
	logrus.Debug("getting posts paginated")

	return allPosts[startIdx:endIdx], nil

}

func (imr *InMemoryRepo) GetComments(postID int) ([]*entity.Comment, error) {
	logrus.Debug("getting comments")

	imr.mu.RLock()
	defer imr.mu.RUnlock()

	var comments []*entity.Comment
	for _, comment := range imr.comments {
		if comment.PostID == postID {
			comments = append(comments, comment)
		}
	}
	logrus.Debug("got comments")

	return comments, nil
}

func (imr *InMemoryRepo) GetRootCommentsPag(postID int, limit, offset *int) ([]*entity.Comment, error) {
	logrus.Debug("getting root comments paginated")

	imr.mu.RLock()
	defer imr.mu.RUnlock()

	var roots []*entity.Comment
	for _, comment := range imr.comments {
		if comment.PostID == postID && comment.ParentID == nil {
			roots = append(roots, comment)
		}
	}

	length := len(roots)
	startIdx, endIdx := 0, length
	if offset != nil {
		startIdx = *offset
		if startIdx < 0 || startIdx > length {
			startIdx = 0
		}
	}
	if limit != nil {
		endIdx = *limit + startIdx
		if endIdx < 0 || endIdx > length {
			endIdx = length
		}
	}

	logrus.Debug("got root comments paginated")
	return roots[startIdx:endIdx], nil
}
