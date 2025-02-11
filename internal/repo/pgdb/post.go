package pgdb

import (
	"Commentary/internal/entity"
	"context"
	"database/sql"
	"errors"
	"github.com/Masterminds/squirrel"
	"github.com/sirupsen/logrus"
)

type PostRepo interface {
	GetPost(ctx context.Context, postID int) (*entity.Post, error)
	AddPost(ctx context.Context, post *entity.Post) (*entity.Post, error)
	ToggleComments(ctx context.Context, postID int) error
	GetPostsPag(ctx context.Context, limit *int, offset *int) ([]*entity.Post, error)
}

type postRepo struct {
	DB  *sql.DB
	SQL squirrel.StatementBuilderType
}

func NewPostRepo(DB *sql.DB) PostRepo {
	return &postRepo{
		DB:  DB,
		SQL: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
	}
}

func (pr *postRepo) GetPost(ctx context.Context, postID int) (*entity.Post, error) {
	logrus.Debugf("getting post by id=%d", postID)
	statement := pr.SQL.
		Select("id", "author_id", "title",
			"content", "created", "commentable").
		From("posts").
		Where(squirrel.Eq{"id": postID})

	query, args, err := statement.ToSql()
	if err != nil {
		logrus.WithError(err).Error("failed to generate SQL query for GetPost")
		return nil, &RepositoryError{
			Operation: "getting post",
			Content:   "failed to generate SQL query",
			Err:       ErrGeneratingSQL,
		}
	}

	var post entity.Post
	err = pr.DB.QueryRowContext(ctx, query, args...).Scan(&post.ID, &post.AuthorID, &post.Title,
		&post.Content, &post.Created, &post.Commentable)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			logrus.WithField("postID", postID).Error("post not found")
			return nil, &RepositoryError{
				Operation: "getting post",
				Content:   "post record not found",
				Err:       err,
			}
		}
		logrus.WithError(err).Error("error getting post")
		return nil, &RepositoryError{
			Operation: "getting post",
			Content:   "internal database error",
			Err:       ErrGettingPost,
		}
	}

	logrus.WithField("postID", post.ID).Debug("got post")
	return &post, nil
}

func (pr *postRepo) AddPost(ctx context.Context, post *entity.Post) (*entity.Post, error) {
	logrus.Debug("adding post")

	statement := pr.SQL.
		Insert("posts").
		Columns("author_id", "title", "content", "created", "commentable").
		Values(post.AuthorID, post.Title, post.Content, post.Created, post.Commentable).
		Suffix("RETURNING id")

	query, args, err := statement.ToSql()
	if err != nil {
		logrus.WithError(err).Error("failed to generate SQL query for AddPost")
		return nil, &RepositoryError{
			Operation: "adding post",
			Content:   "failed to generate SQL query",
			Err:       ErrGeneratingSQL,
		}
	}

	err = pr.DB.QueryRowContext(ctx, query, args...).Scan(&post.ID)
	if err != nil {
		logrus.WithError(err).Error(ErrAddingPost)
		return nil, &RepositoryError{
			Operation: "adding post",
			Content:   "internal database error",
			Err:       ErrAddingPost,
		}
	}
	logrus.WithField("postID", post.ID).Debug("added post")
	return post, nil
}

func (pr *postRepo) ToggleComments(ctx context.Context, postID int) error {
	logrus.WithField("postID", postID).Debug("toggling comments for post")

	statement := pr.SQL.
		Update("posts").
		Set("commentable", squirrel.Expr("NOT commentable")).
		Where(squirrel.Eq{"id": postID})

	query, args, err := statement.ToSql()
	if err != nil {
		logrus.WithError(err).Error("failed to generate SQL query")
		return &RepositoryError{
			Operation: "toggling comments",
			Content:   "failed to generate SQL query",
			Err:       ErrGeneratingSQL,
		}
	}

	_, err = pr.DB.ExecContext(ctx, query, args...)
	if err != nil {
		logrus.WithError(err).Error(ErrTogglingComments)
		return &RepositoryError{
			Operation: "toggling comments",
			Content:   "failed to toggle comments",
			Err:       ErrTogglingComments,
		}
	}
	logrus.WithField("postID", postID).Debug("toggled comments for post")

	return nil
}

func (pr *postRepo) GetPostsPag(ctx context.Context, limit *int, offset *int) ([]*entity.Post, error) {
	logrus.Debugf("getting posts paginated")

	statement := pr.SQL.
		Select("id", "author_id", "title", "content", "created", "commentable").
		From("posts")

	if limit != nil && *limit <= 0 {
		return nil, errors.New("wrong limit parameter value")
	}
	if offset != nil && *offset < 0 {
		return nil, errors.New("wrong offset parameter value")
	}

	if limit != nil {
		statement = statement.Limit(uint64(*limit))
	}
	if offset != nil {
		statement = statement.Offset(uint64(*offset))
	}

	query, args, err := statement.ToSql()
	if err != nil {
		logrus.WithError(err).Error("failed to generate SQL query for AddPost")
		return nil, &RepositoryError{
			Operation: "getting posts",
			Content:   "failed to generate SQL query",
			Err:       ErrGeneratingSQL,
		}
	}

	rows, err := pr.DB.QueryContext(ctx, query, args...)
	if err != nil {
		logrus.WithError(err).Error("failed to get paginated posts")
		return nil, &RepositoryError{
			Operation: "getting posts",
			Content:   "failed to get posts",
			Err:       ErrGettingPosts,
		}
	}
	defer rows.Close()

	var posts []*entity.Post
	for rows.Next() {
		var post entity.Post
		err = rows.Scan(&post.ID, &post.AuthorID, &post.Title,
			&post.Content, &post.Created, &post.Commentable,
		)
		if err != nil {
			logrus.WithError(err).Error("failed to scan data")
			return nil, &RepositoryError{
				Operation: "getting posts",
				Content:   "failed to scan row",
				Err:       err,
			}
		}
		posts = append(posts, &post)
	}

	if err = rows.Err(); err != nil {
		logrus.WithError(err).Error("rows error")
		return nil, &RepositoryError{
			Operation: "getting posts",
			Content:   "rows error",
			Err:       ErrGettingPosts,
		}
	}

	if len(posts) == 0 {
		return nil, &RepositoryError{
			Operation: "getting posts",
			Content:   "no posts found",
		}
	}

	logrus.Debug("got posts paginated")

	return posts, nil
}
