package service

import (
	"Commentary/internal/common"
	"Commentary/internal/entity"
	"Commentary/internal/graph/model"
	"Commentary/internal/pubsub"
	"Commentary/internal/repo/pgdb"
	"context"
	"fmt"
	"sync"
	"time"
)

type commentService struct {
	commentRepo pgdb.CommentRepo
	postRepo    pgdb.PostRepo
	userRepo    pgdb.UserRepo
	postService common.PostService
	broker      *pubsub.Broker
}

func NewCommentService(commentRepo pgdb.CommentRepo, postRepo pgdb.PostRepo, userRepo pgdb.UserRepo,
	postService common.PostService, broker *pubsub.Broker) common.CommentService {
	return &commentService{
		commentRepo: commentRepo, postRepo: postRepo, userRepo: userRepo,
		postService: postService, broker: broker,
	}
}

func (cs *commentService) GetComments(ctx context.Context, postID int, limit *int, offset *int) ([]*model.Comment, error) {
	roots, comments, authors, post, err := cs.getCommentsData(ctx, postID, limit, offset)
	if err != nil {
		return nil, err
	}

	return FillComments(roots, comments, authors, post), nil
}

func (cs *commentService) getCommentsData(ctx context.Context, postID int, limit, offset *int) ([]*entity.Comment,
	[]*entity.Comment, map[int]*model.User, *model.Post, error) {

	var wg sync.WaitGroup
	var errs [3]error
	var roots, allComments []*entity.Comment
	var post *entity.Post

	wg.Add(3)

	go func() {
		defer wg.Done()
		allComments, errs[0] = cs.commentRepo.GetComments(ctx, postID)
	}()

	go func() {
		defer wg.Done()
		roots, errs[1] = cs.commentRepo.GetRootCommentsPag(ctx, postID, limit, offset)
	}()

	go func() {
		defer wg.Done()
		post, errs[2] = cs.postRepo.GetPost(ctx, postID)
	}()

	wg.Wait()

	for _, err := range errs {
		if err != nil {
			return nil, nil, nil, nil, err
		}
	}

	authorIDs := GetAuthorIDs(allComments)
	authorIDs = append(authorIDs, post.AuthorID)

	authors, err := cs.userRepo.GetUsersByIDs(ctx, authorIDs)
	if err != nil {
		return nil, nil, nil, nil, err
	}
	modelPost := &model.Post{
		ID:          post.ID,
		Author:      authors[post.AuthorID],
		Title:       post.Title,
		Content:     post.Content,
		Created:     post.Created,
		Commentable: post.Commentable,
		Comments:    []*model.Comment{},
	}
	return roots, allComments, authors, modelPost, nil
}

func GetAuthorIDs(comments []*entity.Comment) []int {
	uniqueIDs := make(map[int]struct{})
	for _, comment := range comments {
		uniqueIDs[comment.AuthorID] = struct{}{}
	}

	ids := make([]int, 0, len(uniqueIDs))
	for id := range uniqueIDs {
		ids = append(ids, id)
	}
	return ids
}

func (cs *commentService) CreateComment(ctx context.Context,
	comment model.CreateCommentInput) (*model.Comment, error) {

	commentToAdd := &entity.Comment{
		PostID:   comment.PostID,
		AuthorID: comment.AuthorID,
		Content:  comment.Content,
		Created:  time.Now(),
		ParentID: comment.Parent,
	}

	post, err := cs.postService.GetPost(ctx, comment.PostID, nil, nil)
	if err != nil {
		return nil, err
	}

	if !post.Commentable {
		return nil, fmt.Errorf("comments are disabled")
	}

	id, err := cs.commentRepo.AddComment(ctx, commentToAdd)
	if err != nil {
		return nil, err
	}

	author, err := cs.userRepo.GetUserByID(ctx, comment.AuthorID)
	if err != nil {
		return nil, err
	}

	var parent *model.Comment
	if comment.Parent != nil {
		parentComment, err := cs.GetCommentByID(ctx, *comment.Parent)
		if err != nil {
			return nil, err
		}
		parent = parentComment
	}

	added := &model.Comment{
		ID:      id,
		Post:    post,
		Author:  author,
		Content: comment.Content,
		Created: commentToAdd.Created,
		Parent:  parent,
		Replies: []*model.Comment{},
	}

	if cs.broker != nil {
		go func() {
			cs.broker.Publish(comment.PostID, added)
		}()
	}

	return added, nil
}

func FillComments(roots, fullComments []*entity.Comment, users map[int]*model.User,
	post *model.Post) []*model.Comment {
	commentMap := make(map[int]*model.Comment)

	for _, comment := range fullComments {
		commentMap[comment.ID] = commentToModel(comment, users, post)
	}

	for _, comment := range fullComments {
		if comment.ParentID != nil {
			parent := commentMap[*comment.ParentID]
			child := commentMap[comment.ID]
			child.Parent = parent
			parent.Replies = append(parent.Replies, child)
		}
	}

	result := make([]*model.Comment, 0, len(roots))
	for _, root := range roots {
		result = append(result, commentMap[root.ID])
	}

	return result
}

func commentToModel(comment *entity.Comment, authors map[int]*model.User,
	post *model.Post) *model.Comment {
	return &model.Comment{
		ID:      comment.ID,
		Post:    post,
		Author:  authors[comment.AuthorID],
		Content: comment.Content,
		Created: comment.Created,
		Replies: make([]*model.Comment, 0),
	}
}

func (cs *commentService) GetCommentByID(ctx context.Context, id int) (*model.Comment, error) {
	comment, err := cs.commentRepo.GetCommentByID(ctx, id)
	if err != nil {
		return nil, err
	}

	author, err := cs.userRepo.GetUserByID(ctx, comment.AuthorID)
	if err != nil {
		return nil, err
	}

	post, err := cs.postService.GetPost(ctx, comment.PostID, nil, nil)
	if err != nil {
		return nil, err
	}

	return &model.Comment{
		ID:      comment.ID,
		Post:    post,
		Author:  author,
		Content: comment.Content,
		Created: comment.Created,
	}, nil
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
