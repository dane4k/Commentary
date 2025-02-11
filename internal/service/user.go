package service

import (
	"Commentary/internal/common"
	"Commentary/internal/graph/model"
	"Commentary/internal/repo/pgdb"
	"context"
)

type userService struct {
	userRepo pgdb.UserRepo
}

func NewUserService(userRepo pgdb.UserRepo) common.UserService {
	return &userService{userRepo: userRepo}
}

func (us *userService) CreateUser(ctx context.Context, username string) (*model.User, error) {
	return us.userRepo.AddUser(ctx, username)
}
