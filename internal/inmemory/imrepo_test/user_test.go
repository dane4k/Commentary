package imrepo_test

import (
	"Commentary/internal/inmemory/imrepo"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"sync"
	"testing"
)

func TestAddUser(t *testing.T) {
	repo := imrepo.NewInMemoryRepo()

	user, err := repo.AddUser("user1")
	require.NoError(t, err)
	assert.NotNil(t, user)
	assert.Greater(t, user.ID, 0)
	assert.Equal(t, "user1", user.Username)

	usr, err := repo.GetUser(user.ID)
	require.NoError(t, err)
	assert.Equal(t, user, usr)

	_, err = repo.AddUser("user1")
	assert.Error(t, err)
	assert.Equal(t, imrepo.ErrUserAlreadyExists, err)

	var wg sync.WaitGroup
	n := 10
	usernames := make([]string, n)
	for i := 0; i < n; i++ {
		usernames[i] = "user" + uuid.New().String()
	}

	for _, username := range usernames {
		wg.Add(1)
		go func(username string) {
			defer wg.Done()
			_, _ = repo.AddUser(username)
		}(username)
	}
	wg.Wait()

	for _, username := range usernames {
		user, err = repo.AddUser(username)
		assert.Error(t, err)
		assert.Equal(t, imrepo.ErrUserAlreadyExists, err)
	}
}

func TestGetUser(t *testing.T) {
	repo := imrepo.NewInMemoryRepo()

	_, err := repo.GetUser(111)
	assert.Error(t, err)
	assert.Equal(t, imrepo.ErrUserNotFound, err)

	user1, err := repo.AddUser("user1")
	require.NoError(t, err)
	user2, err := repo.AddUser("user2")
	require.NoError(t, err)

	retrievedUser1, err := repo.GetUser(user1.ID)
	require.NoError(t, err)
	assert.Equal(t, user1, retrievedUser1)

	retrievedUser2, err := repo.GetUser(user2.ID)
	require.NoError(t, err)
	assert.Equal(t, user2, retrievedUser2)
}
