package logger

import (
	"Commentary/internal/config"
	"github.com/sirupsen/logrus"
	"os"
)

func InitLogger(cfg *config.Config) {
	logFile, err := os.OpenFile(cfg.Logger.FileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		logrus.WithError(err).Fatal("Failed to create or open log file")
	}

	logrus.SetOutput(logFile)
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})
	logrus.SetLevel(logrus.DebugLevel)
}
