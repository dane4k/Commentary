package pgdb

import (
	"Commentary/internal/entity"
	"Commentary/internal/repo/pgdb"
	"context"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestGetPost(t *testing.T) {
	db, mock, _ := sqlmock.New()
	repo := pgdb.NewPostRepo(db)

	mock.ExpectQuery(`SELECT id, author_id, title, content, created, commentable FROM posts WHERE id = \$1`).
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "author_id", "title", "content", "created", "commentable"}).
			AddRow(1, 1, "title", "content", time.Now(), true))

	post, err := repo.GetPost(context.Background(), 1)
	require.NoError(t, err)
	assert.NotNil(t, post)
	assert.Equal(t, "title", post.Title)
	assert.Equal(t, "content", post.Content)
}

func TestAddPost(t *testing.T) {
	db, mock, _ := sqlmock.New()
	repo := pgdb.NewPostRepo(db)

	mock.ExpectQuery(`INSERT INTO posts`).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	post := &entity.Post{
		AuthorID:    1,
		Title:       "title",
		Content:     "content",
		Created:     time.Now(),
		Commentable: true,
	}

	newPost, err := repo.AddPost(context.Background(), post)
	require.NoError(t, err)
	assert.NotNil(t, newPost)
	assert.Equal(t, 1, newPost.ID)
}

func TestToggleComments(t *testing.T) {
	db, mock, _ := sqlmock.New()
	repo := pgdb.NewPostRepo(db)

	mock.ExpectExec(`UPDATE posts SET commentable = NOT commentable WHERE id = \$1`).
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err := repo.ToggleComments(context.Background(), 1)
	require.NoError(t, err)
}
