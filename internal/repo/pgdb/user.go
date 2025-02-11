package pgdb

import (
	"Commentary/internal/graph/model"
	"context"
	"database/sql"
	"errors"
	"github.com/Masterminds/squirrel"
	"github.com/sirupsen/logrus"
)

type UserRepo interface {
	GetUsersByIDs(ctx context.Context, ids []int) (map[int]*model.User, error)
	GetUserByID(ctx context.Context, id int) (*model.User, error)
	AddUser(ctx context.Context, username string) (*model.User, error)
}

type userRepo struct {
	DB  *sql.DB
	SQL squirrel.StatementBuilderType
}

func NewUserRepo(DB *sql.DB) UserRepo {
	return &userRepo{
		DB:  DB,
		SQL: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
	}
}

func (ur *userRepo) AddUser(ctx context.Context, username string) (*model.User, error) {
	logrus.WithField("username", username).Debug("adding user")

	statement := ur.SQL.
		Insert("users").
		Columns("username").
		Values(username).
		Suffix("RETURNING id")

	query, args, err := statement.ToSql()
	if err != nil {
		logrus.WithError(err).Error("failed to generate SQL query")
		return nil, &RepositoryError{
			Operation: "adding user",
			Content:   "failed to generate SQL query",
			Err:       ErrGeneratingSQL,
		}
	}
	var id int
	err = ur.DB.QueryRowContext(ctx, query, args...).Scan(&id)
	if err != nil {
		logrus.WithError(err).Error("failed to add user")
		return nil, &RepositoryError{
			Operation: "adding user",
			Content:   "failed to add user",
			Err:       ErrAddingUser,
		}
	}

	logrus.WithField("userID", id).Debug("added user")

	return &model.User{ID: id, Username: username}, nil
}

func (ur *userRepo) GetUsersByIDs(ctx context.Context, ids []int) (map[int]*model.User, error) {
	logrus.WithField("ids", ids).Debug("getting users by IDs")
	statement := ur.SQL.
		Select("id", "username").
		From("users").
		Where(squirrel.Eq{"id": ids})
	query, args, err := statement.ToSql()
	if err != nil {
		logrus.WithError(err).Error("failed to get users")
		return nil, &RepositoryError{
			Operation: "getting users",
			Content:   "failed to get users",
			Err:       ErrGettingUsers,
		}
	}

	rows, err := ur.DB.QueryContext(ctx, query, args...)
	if err != nil {
		logrus.WithError(err).Error("failed to get users by IDs")
		return nil, &RepositoryError{
			Operation: "GetUsersByIDs",
			Content:   "failed to get users by IDs",
			Err:       ErrGettingUsers,
		}
	}
	defer rows.Close()

	users := make(map[int]*model.User)
	for rows.Next() {
		var user model.User
		var id int
		err = rows.Scan(&id, &user.Username)
		if err != nil {
			logrus.WithError(err).Error("rows error")
			return nil, &RepositoryError{
				Operation: "GetUsersByIDs",
				Content:   "failed to scan rows",
				Err:       ErrGettingUsers,
			}
		}
		user.ID = id
		users[id] = &user
	}

	logrus.Debug("got users by IDs")
	return users, nil
}

func (ur *userRepo) GetUserByID(ctx context.Context, id int) (*model.User, error) {
	logrus.Debugf("getting user by id=%d", id)
	statement := ur.SQL.
		Select("id", "username").
		From("users").
		Where(squirrel.Eq{"id": id})

	query, args, err := statement.ToSql()
	if err != nil {
		logrus.WithError(err).Error(ErrGeneratingSQL)
		return nil, ErrGettingUser
	}

	row := ur.DB.QueryRowContext(ctx, query, args...)

	var user model.User
	var dbID int
	err = row.Scan(&dbID, &user.Username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			logrus.WithError(err).Error("user not found")
			return nil, &RepositoryError{
				Operation: "getting user by id",
				Content:   "user not found",
			}
		} else {
			logrus.WithError(err).Error("rows error")
			return nil, &RepositoryError{
				Operation: "GetUsersByIDs",
				Content:   "failed to scan rows",
				Err:       ErrGettingUsers,
			}
		}
	}
	user.ID = dbID
	logrus.WithField("userID", dbID).Debug("successfully retrieved user")
	return &user, nil
}
