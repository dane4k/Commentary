package app

import (
	"Commentary/internal/common"
	"Commentary/internal/config"
	"Commentary/internal/db"
	"Commentary/internal/graph"
	"database/sql"
	"github.com/sirupsen/logrus"
)

type App struct {
	Resolver       *graph.Resolver
	CommentService common.CommentService
	PostService    common.PostService
	UserService    common.UserService
}

func InitApp(cfg *config.Config) *App {
	var DB *sql.DB
	var err error
	if cfg.Database.StoreInDB {
		DB, err = db.InitDB(cfg)
		if err != nil {
			logrus.Fatal(err)
		}
	}

	factory := NewServiceFactory(cfg, DB)

	postService := factory.CreatePostService()
	commentService, broker := factory.CreateCommentService()
	postService.SetCommentService(commentService)
	userService := factory.CreateUserService()

	resolver := graph.NewResolver(postService, commentService, userService, broker)

	logrus.Info("Initialized App")

	return &App{
		Resolver:       resolver,
		CommentService: commentService,
		PostService:    postService,
		UserService:    userService,
	}
}
