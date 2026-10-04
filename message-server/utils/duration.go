package utils

import (
	"strconv"
	"time"
)

// Converts env durations to time.Durations, with defaults and string conversion.
func ConfiguredDuration(key string, getenv func(string) string) time.Duration {
	strD := getenv(key)
	dur, err := strconv.ParseUint(strD, 10, 64)
	if strD != "" && err == nil {
		return time.Duration(dur) * time.Millisecond
	}

	switch key {
	case "SQLITE_READ_TIMEOUT":
		dur = 5000
	case "SQLITE_WRITE_TIMEOUT":
		dur = 7500
	case "SQLITE_DELETE_TIMEOUT":
		dur = 7500
	default:
		dur = 10000
	}

	return time.Duration(dur) * time.Millisecond	
}
