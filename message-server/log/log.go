package log

import (
	"fmt"
	"strings"
	"time"
)

type LogLevel int

const (
	DEBUG LogLevel = iota
	INFO
	WARNING
	ERROR
)

var logLevelName = map[LogLevel]string{
	DEBUG:   "DEBUG",
	INFO:    "INFO",
	WARNING: "WARNING",
	ERROR:   "ERROR",
}

func (ll LogLevel) String() string {
	return logLevelName[ll]
}

func determineLogLevel(logLevel []LogLevel) LogLevel {
	if len(logLevel) == 0 {
		return INFO
	}

	if logLevel[0] > 3 {
		return INFO
	}

	return logLevel[0]
}

// Defaults to INFO
func Log(source string, msg string, logLevel ...LogLevel) {
	ll := determineLogLevel(logLevel)
	
	var builtMsg strings.Builder

	// Write datetime
	t := time.Now()
	t.UTC()
	fmt.Fprintf(&builtMsg, "%d-%d-%dT%d:%d:%dZ", t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second())

	// Write log level
	fmt.Fprintf(&builtMsg, " [%s] ", ll.String())

	// Write source
	fmt.Fprintf(&builtMsg, "%s: %s", source, msg)

	// Write out
	fmt.Println(builtMsg.String())
}
