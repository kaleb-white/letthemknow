package utils

import (
	"fmt"
	"reflect"
	"strings"
)

func CheckRequiredFieldsArentDefault(i any, fields *[]string, result *[]error) bool {
	wasError := false
	wasPanic := false
	var nonexistentFields strings.Builder
	nonexistentFields.WriteString("Nonexistent fields passed: ")

	s := reflect.ValueOf(i)
	for _, field := range *fields {
		val := s.FieldByName(field)
		if !val.IsValid() {
			wasPanic = true
			fmt.Fprintf(&nonexistentFields, "%s ", field)
		}
		
		if !wasPanic && val.IsZero() {
			wasError = true
			*result = append(*result, fmt.Errorf("%s is required, found an empty string.", field))	
		}
	}

	if wasPanic {
		panic(nonexistentFields.String())
	}

	return wasError
}
