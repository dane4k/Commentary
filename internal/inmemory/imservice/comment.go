package imservice

import (
	"Commentary/internal/common"
	"Commentary/internal/entity"
	"Commentary/internal/graph/model"
	"Commentary/internal/inmemory/imrepo"
	"Commentary/internal/pubsub"
	"Commentary/internal/service"
	"context"
	"fmt"
)

type commentService struct {
	repo        *imrepo.InMemoryRepo
	postService common.PostService
	broker      *pubsub.Broker
}

func NewCommentService(repo *imrepo.InMemoryRepo, postService common.PostService,
	broker *pubsub.Broker) common.CommentService {
	return &commentService{repo: repo, postService: postService, broker: broker}
}

func (cs *commentService) CreateComment(ctx context.Context, input model.CreateCommentInput) (*model.Comment, error) {

	post, err := cs.postService.GetPost(ctx, input.PostID, nil, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get post: %w", err)
	}

	if !post.Commentable {
		return nil, fmt.Errorf("comments are disabled")
	}

	comment, err := cs.repo.AddComment(input)
	if err != nil {
		return nil, err
	}

	author, err := cs.repo.GetUser(input.AuthorID)
	if err != nil {
		return nil, fmt.Errorf("failed to get author: %w", err)
	}

	authorModel := &model.User{
		ID:       author.ID,
		Username: author.Username,
	}

	added := &model.Comment{
		ID:      comment.ID,
		Post:    post,
		Author:  authorModel,
		Content: comment.Content,
		Created: comment.Created,
	}
	cs.broker.Publish(input.PostID, added)
	return added, nil
}

func (cs *commentService) GetComments(ctx context.Context, postID int, limit *int, offset *int) ([]*model.Comment, error) {
	comments, err := cs.repo.GetComments(postID)
	if err != nil {
		return nil, err
	}

	roots, err := cs.repo.GetRootCommentsPag(postID, limit, offset)
	if err != nil {
		return nil, err
	}
	authorIDs := service.GetAuthorIDs(comments)
	authors, err := cs.repo.GetUsersByIDs(authorIDs)
	if err != nil {
		return nil, err
	}
	commentsMap := make(map[int]*model.Comment)
	for _, comment := range comments {
		commentsMap[comment.ID] = commentToModel(comment, authors, postID)
	}

	for _, comment := range comments {
		if comment.ParentID != nil {
			parent := commentsMap[*comment.ParentID]
			child := commentsMap[comment.ID]
			child.Parent = parent
			parent.Replies = append(parent.Replies, child)
		}
	}

	var res []*model.Comment
	for _, root := range roots {
		res = append(res, commentsMap[root.ID])
	}

	return res, nil
}

func commentToModel(comment *entity.Comment, authors map[int]*model.User, postID int) *model.Comment {
	return &model.Comment{
		ID:      comment.ID,
		Post:    &model.Post{ID: postID},
		Author:  authors[comment.AuthorID],
		Content: comment.Content,
		Created: comment.Created,
		Replies: make([]*model.Comment, 0),
	}
}

func (cs *commentService) GetCommentsForPosts(ctx context.Context, postIDs []int) ([]*model.Comment, error) {
	var allComments []*model.Comment
	for _, postID := range postIDs {
		comments, err := cs.GetComments(ctx, postID, nil, nil)
		if err != nil {
			return nil, err
		}
		allComments = append(allComments, comments...)
	}
	return allComments, nil
}
