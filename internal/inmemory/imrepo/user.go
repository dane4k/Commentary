package imrepo

import (
	"Commentary/internal/entity"
	"Commentary/internal/graph/model"
	"github.com/sirupsen/logrus"
)

func (imr *InMemoryRepo) AddUser(username string) (*entity.User, error) {
	logrus.Debug("adding user")

	imr.mu.Lock()
	defer imr.mu.Unlock()

	if _, ok := imr.nicknames[username]; ok {
		return nil, ErrUserAlreadyExists
	}

	id := len(imr.users) + 1
	user := &entity.User{
		ID:       id,
		Username: username,
	}

	imr.users[id] = user
	imr.nicknames[username] = struct{}{}

	logrus.Debug("added user")

	return user, nil
}

func (imr *InMemoryRepo) GetUser(userID int) (*entity.User, error) {
	logrus.Debug("getting user")

	imr.mu.RLock()
	defer imr.mu.RUnlock()

	if user, ok := imr.users[userID]; ok {
		return user, nil
	}
	logrus.Debug("got user")

	return nil, ErrUserNotFound
}

func (imr *InMemoryRepo) GetUsersByIDs(ids []int) (map[int]*model.User, error) {
	logrus.Debug("getting users by ids")

	imr.mu.RLock()
	defer imr.mu.RUnlock()

	users := make(map[int]*model.User)
	for _, id := range ids {
		if user, ok := imr.users[id]; ok {
			users[id] = &model.User{
				ID:       user.ID,
				Username: user.Username,
			}
		}
	}
	logrus.Debug("got users")

	return users, nil
}
