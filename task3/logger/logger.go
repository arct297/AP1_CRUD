package logger

import (
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
)

var Log *logrus.Logger

func InitLogger() error {
	Log = logrus.New()

	Log.SetFormatter(&logrus.JSONFormatter{})

	file, err := os.OpenFile("activity_logs.json", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		Log.SetOutput(os.Stdout)
		Log.Warn("Failed to log to file, using default stdout")
		return fmt.Errorf("failed to open log file: %v", err)
	}

	Log.SetOutput(file)
	return nil
}
