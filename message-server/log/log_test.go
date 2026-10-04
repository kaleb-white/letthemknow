package log_test

import (
	"os"
	"strings"
	"testing"

	"github.com/kaleb-white/letthemknow/message-server/log"
)

func getenv(k string) string {
	switch k {
	case "LOG_LEVEL":
		return "WARNING"
	default:
		return ""
	}
}

func TestConfigureConfigurationPreventsWriting(t *testing.T) {
	buf := strings.Builder{}
	log.Configure(&buf, os.Stderr, getenv)
	log.Log("test", "test", log.INFO)

	if buf.Len() > 0 {
		t.Errorf("Log logged when it shouldn't have: %s", buf.String())
	}
}

func TestLogWritesToStdOut(t *testing.T) {
	buf := strings.Builder{}
	log.Configure(os.Stdout, &buf, getenv)
	log.Log("test", "test", log.ERROR)

	if buf.Len() == 0 {
		t.Errorf("Log didn't log when it should have: %s", buf.String())
	}
}

func TestLogWritesAtMatchingLogLevel(t *testing.T) {
	buf := strings.Builder{}
	log.Configure(&buf, os.Stderr, getenv)
	log.Log("test", "test", log.WARNING)

	if buf.Len() == 0 {
		t.Errorf("Log didn't log when it should have: %s", buf.String())
	}
}
