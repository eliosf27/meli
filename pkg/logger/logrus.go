package logger

import "github.com/sirupsen/logrus"

type Config interface {
	GetString(key string) string
}

func NewLogrus(config Config) (*logrus.Logger, error) {
	logger := logrus.New()
	logger.Formatter = &logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: config.GetString("logger.timestamp.format"),
	}
	level, err := logrus.ParseLevel(config.GetString("logger.level"))
	if err != nil {
		return nil, err
	}
	logger.Level = level
	return logger, nil
}
