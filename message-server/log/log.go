package log

import (
	"fmt"
	"io"
	"os"
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

const DEFAULT_LOG_LEVEL LogLevel = INFO

var logLevelName = map[LogLevel]string{
	DEBUG:   "DEBUG",
	INFO:    "INFO",
	WARNING: "WARNING",
	ERROR:   "ERROR",
}

var nameLogLevel = map[string]LogLevel{
	"DEBUG":   DEBUG,
	"INFO":    INFO,
	"WARNING": WARNING,
	"ERROR":   ERROR,
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


var stdOut io.Writer = os.Stdout
var stdErr io.Writer = os.Stdin

// Unless overwritten use INFO as default log level
var getEnv = func(k string) string {
	switch k {
	case "LOG_LEVEL":
		return "INFO"
	default:
		return ""
	}
}

func logLevelShouldLog(logLevel LogLevel) bool {
	// Use env configured string or package default
	configuredLogLevelString := getEnv("LOG_LEVEL")
	configuredLogLevel := DEFAULT_LOG_LEVEL
	if configuredLogLevelString != "" {
		configuredLogLevel = nameLogLevel[configuredLogLevelString]
	}
	
	return logLevel >= configuredLogLevel
}

func Configure(stdout, stderr io.Writer, getenv func(string) string) {
	if stdout != nil {
		stdOut = stdout
	}
	if stderr != nil {
		stdErr = stderr
	}
	getEnv = getenv
}

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
	if determineLogLevel(logLevel) == ERROR {
		fmt.Fprint(stdErr, builtMsg.String())
	} else if logLevelShouldLog(determineLogLevel(logLevel)) {
		fmt.Fprint(stdOut, builtMsg.String())
	}
}
