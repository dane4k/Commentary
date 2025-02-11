package db

import (
	"Commentary/internal/config"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"time"
)

func InitDB(cfg *config.Config) (*sql.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s dbname=%s port=%d user=%s password=%s sslmode=disable",
		cfg.Database.Host,
		cfg.Database.Name,
		cfg.Database.Port,
		cfg.Database.User,
		cfg.Database.Password,
	)

	DB, err := sql.Open("postgres", dsn)
	if err != nil {
		logrus.WithError(err).Fatal("error connecting to database")
		return nil, err
	}
	DB.SetMaxOpenConns(20)
	DB.SetMaxIdleConns(10)
	DB.SetConnMaxLifetime(time.Minute * 30)

	return DB, nil
}
