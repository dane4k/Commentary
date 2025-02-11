package app

import (
	"Commentary/internal/common"
	"Commentary/internal/config"
	"Commentary/internal/inmemory/imrepo"
	"Commentary/internal/inmemory/imservice"
	"Commentary/internal/pubsub"
	"Commentary/internal/repo/pgdb"
	"Commentary/internal/service"
	"database/sql"
)

type ServiceFactory struct {
	cfg         *config.Config
	db          *sql.DB
	imRepo      *imrepo.InMemoryRepo
	postService common.PostService
}

func NewServiceFactory(cfg *config.Config, db *sql.DB) *ServiceFactory {
	return &ServiceFactory{
		cfg:    cfg,
		db:     db,
		imRepo: imrepo.NewInMemoryRepo(),
	}
}

func (f *ServiceFactory) CreatePostService() common.PostService {
	if f.postService == nil {
		if f.cfg.Database.StoreInDB {
			f.postService = service.NewPostService(
				pgdb.NewPostRepo(f.db),
				pgdb.NewUserRepo(f.db),
			)
		} else {
			f.postService = imservice.NewPostService(f.imRepo)
		}
	}
	return f.postService
}

func (f *ServiceFactory) CreateCommentService() (common.CommentService, *pubsub.Broker) {
	broker := pubsub.NewBroker()
	if f.cfg.Database.StoreInDB {
		return service.NewCommentService(
			pgdb.NewCommentRepo(f.db),
			pgdb.NewPostRepo(f.db),
			pgdb.NewUserRepo(f.db),
			f.CreatePostService(),
			broker), broker
	}
	return imservice.NewCommentService(
		f.imRepo,
		f.CreatePostService(),
		broker), broker
}

func (f *ServiceFactory) CreateUserService() common.UserService {
	if f.cfg.Database.StoreInDB {
		return service.NewUserService(pgdb.NewUserRepo(f.db))
	}
	return imservice.NewUserService(f.imRepo)
}
