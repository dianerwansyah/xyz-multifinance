package logger

import (
	"os"
	"strings"

	"github.com/sirupsen/logrus"
)

var Log *logrus.Logger

func Setup() {
	Log = logrus.New()

	env := strings.ToLower(os.Getenv("APP_ENV"))
	if env == "local" {
		Log.SetFormatter(&logrus.TextFormatter{
			FullTimestamp: true,
		})
	}

	Log.SetOutput(os.Stdout)
	Log.SetLevel(logrus.InfoLevel)
}
