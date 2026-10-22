package utils_test

import (
	"errors"
	"regexp"
	"testing"

	"github.com/kaleb-white/letthemknow/message-server/utils"
)

var hasMultiErrorCopy = regexp.MustCompile(`Errors:`)

func TestPrettyPrintErrorsOneError(t *testing.T) {
	var oneErr []error = make([]error, 0, 1)
	oneErr = append(oneErr, errors.New("Test"))	

	errs := utils.PrettyPrintErrors(oneErr)
	if errs == "" {
		t.Errorf("Pretty print errors returned empty string when called with one error")
	} else if hasMultiErrorCopy.MatchString(errs) {
		t.Errorf("Pretty print errors returned multi error copy when called with one error")
	}
}

func TestPrettyPrintErrorsMultiError(t *testing.T) {
	var multiErr []error = make([]error, 1)
	multiErr = append(multiErr, errors.New("Test"), errors.New("Test2"))	

	errs := utils.PrettyPrintErrors(multiErr)
	if errs == "" {
		t.Errorf("Pretty print errors returned empty string when called with multi error")
	} else if !hasMultiErrorCopy.MatchString(errs) {
		t.Errorf("Pretty print errors did not return multi error copy when called with multi error")
	}
}


