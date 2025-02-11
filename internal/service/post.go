package service

import (
	"Commentary/internal/common"
	"Commentary/internal/entity"
	"Commentary/internal/graph/model"
	"Commentary/internal/repo/pgdb"
	"context"
	"errors"
	"fmt"
	"log"
	"sync"
	"time"
)

type PostService struct {
	postRepo       pgdb.PostRepo
	userRepo       pgdb.UserRepo
	commentService common.CommentService
}

func NewPostService(postRepo pgdb.PostRepo, userRepo pgdb.UserRepo) common.PostService {
	return &PostService{postRepo: postRepo, userRepo: userRepo}
}

func (ps *PostService) SetCommentService(commentService common.CommentService) {
	ps.commentService = commentService
}

func (ps *PostService) CreatePost(ctx context.Context, input model.CreatePostInput) (*model.Post, error) {
	newPost := &entity.Post{
		AuthorID:    input.AuthorID,
		Title:       input.Title,
		Content:     input.Content,
		Created:     time.Now(),
		Commentable: input.Commentable,
	}

	post, err := ps.postRepo.AddPost(ctx, newPost)
	if err != nil {
		return nil, err
	}
	author, err := ps.userRepo.GetUserByID(ctx, input.AuthorID)
	if err != nil {
		return nil, err
	}

	return &model.Post{
		ID:          post.ID,
		Author:      author,
		Title:       post.Title,
		Content:     post.Content,
		Created:     post.Created,
		Commentable: post.Commentable,
		Comments:    []*model.Comment{},
	}, nil
}

func (ps *PostService) GetPost(ctx context.Context, id int, limit *int, offset *int) (*model.Post, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	var wg sync.WaitGroup
	var postErr, commentsErr, authorErr error
	var post *entity.Post
	var comments []*model.Comment
	var author *model.User

	postChan := make(chan *entity.Post, 1)
	wg.Add(3)

	go func() {
		defer wg.Done()
		post, postErr = ps.postRepo.GetPost(ctx, id)
		if postErr != nil {
			cancel()
			return
		}
		postChan <- post
	}()

	go func() {
		defer wg.Done()
		comments, commentsErr = ps.commentService.GetComments(ctx, id, limit, offset)
		if commentsErr != nil {
			cancel()
			return
		}
	}()

	go func() {
		defer wg.Done()
		select {
		case p := <-postChan:
			if p == nil {
				cancel()
				return
			}
			if post != nil {
				a, authorErr := ps.userRepo.GetUserByID(ctx, post.AuthorID)
				if authorErr != nil {
					cancel()
					return
				}
				author = a
			}
		case <-ctx.Done():
			authorErr = errors.New("context canceled")
			return
		}
	}()

	wg.Wait()

	if postErr != nil || commentsErr != nil || authorErr != nil {
		log.Println(postErr)
		log.Println(commentsErr)
		log.Println(authorErr)
		return nil, fmt.Errorf("error getting post")
	}

	return &model.Post{
		ID:          id,
		Author:      author,
		Title:       post.Title,
		Content:     post.Content,
		Created:     post.Created,
		Commentable: post.Commentable,
		Comments:    comments,
	}, nil
}

func (ps *PostService) ToggleComments(ctx context.Context, postID int) (*model.Post, error) {
	err := ps.postRepo.ToggleComments(ctx, postID)
	if err != nil {
		return nil, err
	}

	return ps.GetPost(ctx, postID, nil, nil)
}

func (ps *PostService) GetPosts(ctx context.Context, limit *int, offset *int) ([]*model.Post, error) {
	posts, err := ps.postRepo.GetPostsPag(ctx, limit, offset)
	if err != nil {
		return nil, err
	}

	authorIDs := extractAuthorIDs(posts)
	users, err := ps.userRepo.GetUsersByIDs(ctx, authorIDs)
	if err != nil {
		return nil, err
	}

	postIDs := make([]int, len(posts))
	for i, post := range posts {
		postIDs[i] = post.ID
	}

	comments, err := ps.commentService.GetCommentsForPosts(ctx, postIDs)
	if err != nil {
		return nil, err
	}

	commentsByPID := groupComments(comments)

	var modelPosts []*model.Post
	for _, post := range posts {
		modelPosts = append(modelPosts, &model.Post{
			ID:          post.ID,
			Author:      users[post.AuthorID],
			Title:       post.Title,
			Content:     post.Content,
			Created:     post.Created,
			Commentable: post.Commentable,
			Comments:    commentsByPID[post.ID],
		})
	}
	return modelPosts, nil
}

func extractAuthorIDs(posts []*entity.Post) []int {
	ids := make(map[int]bool)
	for _, post := range posts {
		ids[post.AuthorID] = true
	}
	result := make([]int, 0, len(ids))
	for id := range ids {
		result = append(result, id)
	}
	return result
}

func groupComments(comments []*model.Comment) map[int][]*model.Comment {
	grouped := make(map[int][]*model.Comment)
	for _, comment := range comments {
		grouped[comment.Post.ID] = append(grouped[comment.Post.ID], comment)
	}
	return grouped
}
