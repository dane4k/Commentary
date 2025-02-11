package pgdb_test

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

func TestAddComment(t *testing.T) {
	db, mock, _ := sqlmock.New()
	repo := pgdb.NewCommentRepo(db)

	mock.ExpectQuery(`INSERT INTO comments`).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	comment := &entity.Comment{
		PostID:   1,
		AuthorID: 1,
		Content:  "comment1",
		Created:  time.Now(),
	}
	id, err := repo.AddComment(context.Background(), comment)
	require.NoError(t, err)
	assert.Equal(t, 1, id)
	assert.Equal(t, "comment1", comment.Content)
}

func TestGetComments(t *testing.T) {
	db, mock, _ := sqlmock.New()
	repo := pgdb.NewCommentRepo(db)

	mock.ExpectQuery("SELECT id, post_id, author_id, content, created, parent_id FROM comments").
		WillReturnRows(sqlmock.NewRows([]string{"id", "post_id", "author_id",
			"content", "created", "parent_id"}).
			AddRow(1, 1, 1, "comment1", time.Now(), nil))

	comments, err := repo.GetComments(context.Background(), 1)
	require.NoError(t, err)
	assert.Len(t, comments, 1)
}

func TestGetCommentByID(t *testing.T) {
	db, mock, _ := sqlmock.New()
	repo := pgdb.NewCommentRepo(db)

	mock.ExpectQuery(`SELECT id, post_id, author_id, content, created, parent_id FROM comments WHERE id = \$1`).
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "post_id", "author_id", "content", "created", "parent_id"}).
			AddRow(1, 1, 1, "comment1", time.Now(), nil))

	comment, err := repo.GetCommentByID(context.Background(), 1)
	require.NoError(t, err)
	assert.NotNil(t, comment)
	assert.Equal(t, "comment1", comment.Content)
}
