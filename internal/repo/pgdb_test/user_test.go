package pgdb

import (
	"Commentary/internal/repo/pgdb"
	"context"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestAddUser(t *testing.T) {
	db, mock, _ := sqlmock.New()
	repo := pgdb.NewUserRepo(db)

	mock.ExpectQuery("INSERT INTO users").WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	user, err := repo.AddUser(context.Background(), "user1")
	require.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, "user1", user.Username)
	assert.Equal(t, 1, user.ID)
}

func TestGetUserByID(t *testing.T) {
	db, mock, _ := sqlmock.New()
	repo := pgdb.NewUserRepo(db)

	mock.ExpectQuery(`SELECT id, username FROM users WHERE id = \$1`).
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "username"}).AddRow(1, "user1"))

	user, err := repo.GetUserByID(context.Background(), 1)
	require.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, "user1", user.Username)
	assert.Equal(t, 1, user.ID)
}
