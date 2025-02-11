package imrepo

import (
	"Commentary/internal/entity"
	"Commentary/internal/graph/model"
	"Commentary/internal/inmemory/imrepo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestAddComment(t *testing.T) {
	repo := imrepo.NewInMemoryRepo()
	user, _ := repo.AddUser("user1")
	repo.Posts[1] = &entity.Post{ID: 1, Title: "Post1"}

	input := model.CreateCommentInput{
		PostID:   1,
		AuthorID: user.ID,
		Content:  "Comment1",
	}

	comment, err := repo.AddComment(input)
	require.NoError(t, err)
	assert.NotNil(t, comment)
	assert.Equal(t, input.Content, comment.Content)
	assert.Equal(t, input.AuthorID, comment.AuthorID)
}

func TestGetComment(t *testing.T) {
	repo := imrepo.NewInMemoryRepo()

	_, err := repo.GetComment(5)
	assert.Error(t, err)
	assert.Equal(t, imrepo.ErrCommentNotFound, err)
}

func TestGetCommentsByPostID(t *testing.T) {
	repo := imrepo.NewInMemoryRepo()
	user, _ := repo.AddUser("testuser")
	repo.Posts[1] = &entity.Post{ID: 1, Title: "Post1"}

	input := model.CreateCommentInput{
		PostID:   1,
		AuthorID: user.ID,
		Content:  "Comment1",
		Parent:   nil,
	}
	comment, _ := repo.AddComment(input)

	comments, err := repo.GetCommentsByPostID(1)
	require.NoError(t, err)
	assert.Len(t, comments, 1)
	assert.Equal(t, comment.Content, comments[0].Content)
}
