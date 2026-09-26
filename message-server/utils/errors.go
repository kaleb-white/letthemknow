package utils

import (
	"fmt"
	"strings"

	"github.com/kaleb-white/letthemknow/message-server/log"
)

func PrettyPrintErrors(errs []error) string {
	numErrors := len(errs)

	switch numErrors {
	case 0: 
		return ""
	case 1:
		return errs[0].Error()
	}

	var res strings.Builder
	res.WriteString("Errors:\n")
	 
	for idx, err := range errs {
		if err == nil {
			continue
		}

		res.WriteString("\t- ")
		res.WriteString(err.Error())

		if idx >= numErrors - 1 {
			break
		}
	  res.WriteString("\n")
	}

	return res.String()
}

func DebugLogUnexpectedValue(source string, expected string, actual string) {
	log.Log(source, fmt.Sprintf("Got: %s, Expected: %s", actual, expected), log.DEBUG)
}

