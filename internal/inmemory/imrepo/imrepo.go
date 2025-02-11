package imrepo

import (
	"Commentary/internal/entity"
	"sync"
)

type InMemoryRepo struct {
	Posts       map[int]*entity.Post
	comments    map[int]*entity.Comment
	users       map[int]*entity.User
	nicknames   map[string]struct{}
	subscribers map[int]map[int]struct{}
	mu          sync.RWMutex
}

func NewInMemoryRepo() *InMemoryRepo {
	return &InMemoryRepo{
		Posts:       make(map[int]*entity.Post),
		comments:    make(map[int]*entity.Comment),
		users:       make(map[int]*entity.User),
		nicknames:   make(map[string]struct{}),
		subscribers: make(map[int]map[int]struct{}),
	}
}
