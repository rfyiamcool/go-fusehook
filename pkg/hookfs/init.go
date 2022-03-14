package hookfs

import (
	"os"

	log "github.com/sirupsen/logrus"
)

// LogLevelMin is the minimum log level
const LogLevelMin = 0

// LogLevelMax is the maximum log level
const LogLevelMax = 2

var logLevel int

func initLog() {
	// log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stderr)
}

// LogLevel gets the log level.
func LogLevel() int {
	return logLevel
}

// SetLogLevel sets the log level. newLevel must be >= LogLevelMin, and <= LogLevelMax.
func SetLogLevel(level log.Level) {
	log.SetLevel(level)
}

func init() {
	initLog()
	SetLogLevel(0)
}
