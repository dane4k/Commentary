package imservice

import (
	"Commentary/internal/common"
	"Commentary/internal/graph/model"
	"Commentary/internal/inmemory/imrepo"
	"context"
)

type userService struct {
	repo *imrepo.InMemoryRepo
}

func NewUserService(repo *imrepo.InMemoryRepo) common.UserService {
	return &userService{repo: repo}
}

func (us *userService) CreateUser(ctx context.Context, username string) (*model.User, error) {
	addedUser, err := us.repo.AddUser(username)
	if err != nil {
		return nil, err
	}
	return &model.User{
		ID:       addedUser.ID,
		Username: username,
	}, nil
}
