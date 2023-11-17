package logger

import (
	"os"

	"github.com/sirupsen/logrus"
)

// OutputLog is used to output logs to a external .log file that ideally should
// not be in this root directory. Any other configuration to OutputLog can be made here.
var OutputLog = logrus.New()

func init() {
	// The logs directory path has been set to this root directory, but do not keep them here. Add your own path to a
	// directory outside of the root project, like /var/log/myapp for Unix/Linux systems, where you can separate the concerns
	// or create log rotation if needed.
	file, err := os.OpenFile("logs/api.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0664)
	if err != nil {
		logrus.Fatal("Failed to open log file: ", err)
	}
	OutputLog.Out = file
}

// Log is used to output the logs to the console in the development mode.
// Any other configuration to Log can be made here.
var Log = logrus.New()
