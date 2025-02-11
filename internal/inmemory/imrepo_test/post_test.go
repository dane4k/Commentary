package imrepo

import (
	"Commentary/internal/entity"
	"Commentary/internal/graph/model"
	"Commentary/internal/inmemory/imrepo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestAddPost(t *testing.T) {
	repo := imrepo.NewInMemoryRepo()
	user, _ := repo.AddUser("user1")

	input := model.CreatePostInput{
		AuthorID:    user.ID,
		Title:       "Post1",
		Content:     "Post1",
		Commentable: true,
	}

	post, err := repo.AddPost(input)
	require.NoError(t, err)
	assert.NotNil(t, post)
	assert.Equal(t, input.Title, post.Title)
	assert.Equal(t, input.Content, post.Content)
	assert.True(t, post.Commentable)
}

func TestGetPost(t *testing.T) {
	repo := imrepo.NewInMemoryRepo()
	_, err := repo.GetPost(5)
	assert.Error(t, err)
	assert.Equal(t, imrepo.ErrPostNotFound, err)
}

func TestToggleComments(t *testing.T) {
	repo := imrepo.NewInMemoryRepo()
	user, _ := repo.AddUser("user1")
	repo.Posts[1] = &entity.Post{ID: 1, AuthorID: user.ID, Title: "Post1", Commentable: true}

	post, err := repo.ToggleComments(1)
	require.NoError(t, err)
	assert.False(t, post.Commentable)
}
