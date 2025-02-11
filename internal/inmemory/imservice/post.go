package imservice

import (
	"Commentary/internal/common"
	"Commentary/internal/entity"
	"Commentary/internal/graph/model"
	"Commentary/internal/inmemory/imrepo"
	"context"
)

type PostService struct {
	commentService common.CommentService
	repo           *imrepo.InMemoryRepo
}

func NewPostService(repo *imrepo.InMemoryRepo) common.PostService {
	return &PostService{repo: repo}
}

func (ps *PostService) SetCommentService(commentService common.CommentService) {
	ps.commentService = commentService
}

func (ps *PostService) CreatePost(ctx context.Context, input model.CreatePostInput) (*model.Post, error) {
	post, err := ps.repo.AddPost(input)
	if err != nil {
		return nil, err
	}

	user, err := ps.repo.GetUser(post.AuthorID)
	if err != nil {
		return nil, err
	}

	author := &model.User{
		ID:       post.AuthorID,
		Username: user.Username,
	}

	return &model.Post{
		ID:          post.ID,
		Author:      author,
		Title:       post.Title,
		Content:     post.Content,
		Created:     post.Created,
		Commentable: post.Commentable,
		Comments:    nil,
	}, nil
}

func (ps *PostService) GetPost(ctx context.Context, id int, limit, offset *int) (*model.Post, error) {
	post, err := ps.repo.GetPost(id)
	if err != nil {
		return nil, err
	}

	author, err := ps.repo.GetUser(post.AuthorID)
	if err != nil {
		return nil, err
	}

	authorModel := &model.User{
		ID:       author.ID,
		Username: author.Username,
	}

	comments, err := ps.commentService.GetComments(ctx, id, limit, offset)
	if err != nil {
		return nil, err
	}

	return &model.Post{
		ID:          id,
		Author:      authorModel,
		Title:       post.Title,
		Content:     post.Content,
		Created:     post.Created,
		Commentable: post.Commentable,
		Comments:    comments,
	}, nil
}

func (ps *PostService) ToggleComments(ctx context.Context, postID int) (*model.Post, error) {

	_, err := ps.repo.ToggleComments(postID)
	if err != nil {
		return nil, err
	}

	return ps.GetPost(ctx, postID, nil, nil)
}

func (ps *PostService) GetPosts(ctx context.Context, limit, offset *int) ([]*model.Post, error) {
	posts, err := ps.repo.GetPostsPag(limit, offset)
	if err != nil {
		return nil, err
	}
	authorIDs := getAuthorIDs(posts)

	authors, err := ps.repo.GetUsersByIDs(authorIDs)
	if err != nil {
		return nil, err
	}

	postIDs := make([]int, len(posts))
	for i, post := range posts {
		postIDs[i] = post.ID
	}

	comms, err := ps.commentService.GetCommentsForPosts(ctx, postIDs)
	if err != nil {
		return nil, err
	}

	commentsByPID := groupComments(comms)

	var modelPosts []*model.Post
	for _, post := range posts {
		modelPosts = append(modelPosts, &model.Post{
			ID:          post.ID,
			Author:      authors[post.AuthorID],
			Title:       post.Title,
			Content:     post.Content,
			Created:     post.Created,
			Commentable: post.Commentable,
			Comments:    commentsByPID[post.ID],
		})
	}

	return modelPosts, nil
}

func getAuthorIDs(posts []*entity.Post) []int {
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
