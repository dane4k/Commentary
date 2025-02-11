package pgdb

import (
	"Commentary/internal/entity"
	"context"
	"database/sql"
	"errors"
	"github.com/Masterminds/squirrel"
	"github.com/sirupsen/logrus"
)

type CommentRepo interface {
	GetComments(ctx context.Context, postID int) ([]*entity.Comment, error)
	GetRootCommentsPag(ctx context.Context, postID int, limit *int, offset *int) ([]*entity.Comment, error)
	AddComment(ctx context.Context, comment *entity.Comment) (int, error)
	GetCommentByID(ctx context.Context, id int) (*entity.Comment, error)
}

type commentRepo struct {
	DB  *sql.DB
	SQL squirrel.StatementBuilderType
}

func NewCommentRepo(DB *sql.DB) CommentRepo {
	return &commentRepo{
		DB:  DB,
		SQL: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
	}
}

func (cr *commentRepo) AddComment(ctx context.Context, comment *entity.Comment) (int, error) {
	logrus.WithField("postID", comment.PostID).Debug("adding comment for post")

	statement := cr.SQL.Insert("comments").
		Columns("post_id", "author_id", "content", "created", "parent_id").
		Values(comment.PostID, comment.AuthorID, comment.Content, comment.Created, comment.ParentID).
		Suffix("RETURNING id")

	query, args, err := statement.ToSql()
	if err != nil {
		logrus.WithError(err).Error("failed to generate SQL query for AddComment")
		return 0, &RepositoryError{
			Operation: "adding comment",
			Content:   "failed to generate SQL query",
			Err:       ErrGeneratingSQL,
		}
	}

	var id int
	err = cr.DB.QueryRowContext(ctx, query, args...).Scan(&id)
	if err != nil {
		logrus.WithError(err).Error(ErrAddingComment)
		return 0, &RepositoryError{
			Operation: "adding comment",
			Content:   "failed to add comment",
			Err:       ErrAddingComment,
		}
	}
	logrus.WithField("commentID", id).Debug("added comment")
	return id, nil
}

func (cr *commentRepo) GetComments(ctx context.Context, postID int) ([]*entity.Comment, error) {
	logrus.WithField("postID", postID).Debug("getting comments for post")

	statement := cr.SQL.Select("id", "post_id", "author_id",
		"content", "created", "parent_id").
		From("comments").
		Where(squirrel.Eq{"post_id": postID})

	query, args, err := statement.ToSql()
	if err != nil {
		logrus.WithError(err).Error("failed to generate SQL query")
		return nil, &RepositoryError{
			Operation: "getting comments",
			Content:   "failed to generate SQL query",
			Err:       ErrGeneratingSQL,
		}
	}

	rows, err := cr.DB.QueryContext(ctx, query, args...)
	if err != nil {
		logrus.WithError(err).Error("failed to get comments")
		return nil, &RepositoryError{
			Operation: "getting comments",
			Content:   "failed to get comments",
			Err:       ErrGettingComments,
		}
	}
	defer rows.Close()

	var comments []*entity.Comment
	for rows.Next() {
		var comment entity.Comment
		err = rows.Scan(&comment.ID, &comment.PostID, &comment.AuthorID,
			&comment.Content, &comment.Created, &comment.ParentID)
		if err != nil {
			logrus.WithError(err).Error("failed to scan row")
			return nil, &RepositoryError{
				Operation: "getting comments",
				Content:   "failed to scan row",
				Err:       ErrGettingComments,
			}
		}
		comments = append(comments, &comment)
	}

	if err = rows.Err(); err != nil {
		logrus.WithError(err).Error("rows error")
		return nil, &RepositoryError{
			Operation: "getting comments",
			Content:   "rows error",
			Err:       ErrGettingComments,
		}
	}
	logrus.WithField("postID", postID).Debug("successfully retrieved comments")

	return comments, nil
}

func (cr *commentRepo) GetRootCommentsPag(ctx context.Context, postID int, limit *int, offset *int) ([]*entity.Comment, error) {
	logrus.Debugf("getting root comments paginated for post id=%d", postID)
	statement := cr.SQL.Select("id", "post_id", "author_id",
		"content", "created", "parent_id").
		From("comments").
		Where(squirrel.Eq{"post_id": postID}).
		Where(squirrel.Expr("parent_id IS NULL"))

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
		logrus.WithError(err).Error("failed to generate SQL query")
		return nil, &RepositoryError{
			Operation: "getting root comments",
			Content:   "failed to generate SQL query",
			Err:       ErrGeneratingSQL,
		}
	}

	rows, err := cr.DB.QueryContext(ctx, query, args...)
	if err != nil {
		logrus.WithError(err).Error("failed to get root comments paginated")
		return nil, &RepositoryError{
			Operation: "getting root comments",
			Content:   "failed to get root comments",
			Err:       ErrGettingComments,
		}
	}
	defer rows.Close()

	var comments []*entity.Comment
	for rows.Next() {
		var comment entity.Comment
		err = rows.Scan(&comment.ID, &comment.PostID, &comment.AuthorID,
			&comment.Content, &comment.Created, &comment.ParentID)
		if err != nil {
			logrus.WithError(err).Error("failed to scan comment data")
			return nil, &RepositoryError{
				Operation: "GetRootCommentsPag",
				Content:   "failed to scan row",
				Err:       ErrGettingComments,
			}
		}
		comments = append(comments, &comment)
	}
	if err = rows.Err(); err != nil {
		logrus.WithError(err).Error("error iterating over rows")
		return nil, &RepositoryError{
			Operation: "GetRootCommentsPag",
			Content:   "rows error",
			Err:       ErrGettingComments,
		}
	}
	logrus.Debug("got root comments")
	return comments, nil
}

func (cr *commentRepo) GetCommentByID(ctx context.Context, id int) (*entity.Comment, error) {
	logrus.WithField("commentID", id).Debug("getting comment by id")

	statement := cr.SQL.Select("id", "post_id", "author_id", "content", "created", "parent_id").
		From("comments").
		Where(squirrel.Eq{"id": id})

	query, args, err := statement.ToSql()
	logrus.WithError(err).Error(ErrGeneratingSQL)
	if err != nil {
		logrus.WithError(err).Error("failed to generate SQL query")
		return nil, &RepositoryError{
			Operation: "getting comment by id",
			Content:   "failed to generate SQL query",
			Err:       ErrGeneratingSQL,
		}
	}

	row := cr.DB.QueryRowContext(ctx, query, args...)
	var comment entity.Comment
	err = row.Scan(&comment.ID, &comment.PostID, &comment.AuthorID, &comment.Content, &comment.Created, &comment.ParentID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			logrus.WithError(err).Error("comment not found")
			return nil, &RepositoryError{
				Operation: "getting comment by id",
				Content:   "comment not found",
			}
		}
		logrus.WithError(err).Error("failed to get comment by id")
		return nil, &RepositoryError{
			Operation: "getting comment by id",
			Content:   "failed to get comment by id",
			Err:       ErrGettingComment,
		}
	}

	logrus.WithField("commentID", comment.ID).Debug("got comment")
	return &comment, nil
}
