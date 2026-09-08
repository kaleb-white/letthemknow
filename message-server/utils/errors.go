package utils

import (
	"strings"
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
