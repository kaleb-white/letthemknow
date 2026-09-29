package utils

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"
)

func isTimeStruct(val reflect.Value) bool {
	return val.Type() == reflect.ValueOf(time.Now()).Type()
}

func doTimeValsMatch(val1 reflect.Value, val2 reflect.Value) bool {
	val1AsTime, ok := (val1.Interface()).(time.Time)
	if !ok {
		return false
	}
	val2AsTime, ok := (val2.Interface()).(time.Time)
	if !ok {
		return false
	}
	return val1AsTime.Equal(val2AsTime)
}

// Does not conduct recursive checkoing itself, so deeply nested times will be compared using ==, not time.Equal.
func CheckFieldEquality(this any, other any, fields *[]string, result *[]error) bool {
	if reflect.TypeOf(this) != reflect.TypeOf(other) {
		*result = append(*result, errors.New("This and other have mismatched types."))
		return true 
	}

	wasError := false

	wasUnsupportedPanic := false
	var unsupportedKinds strings.Builder
	unsupportedKinds.WriteString("Unsupported kinds tested: ")

	thisValue := reflect.ValueOf(this)
	otherValue := reflect.ValueOf(other)
	for _, field := range *fields {
		thisField := thisValue.FieldByName(field)
		otherField := otherValue.FieldByName(field)

		// Switch on kind
		switch thisField.Kind() {
		case reflect.Invalid, reflect.Chan, reflect.Func, reflect.Interface: 
			wasUnsupportedPanic = true
			fmt.Fprintf(&unsupportedKinds, "%s ", thisField.Type())
		
		case reflect.Struct:
			if isTimeStruct(thisField) && !doTimeValsMatch(thisField, otherField) {
				wasError = true
				*result = append(*result, errors.New(field))
			} else if !reflect.DeepEqual(thisField.Interface(), otherField.Interface()){
				wasError = true
				*result = append(*result, errors.New(field))
			}
		default:
			if !reflect.DeepEqual(fmt.Sprint(thisField), fmt.Sprint(otherField)){
				wasError = true
				*result = append(*result, errors.New(field))
			}
		}
	}

	if wasUnsupportedPanic {
		panic(unsupportedKinds.String())
	}

	return wasError
}

func CheckRequiredFieldsArentDefault(i any, fields *[]string, result *[]error) bool {
	wasError := false

	s := reflect.ValueOf(i)
	for _, field := range *fields {
		val := s.FieldByName(field)
		
		if !val.IsValid() || val.IsZero() {
			wasError = true
			*result = append(*result, fmt.Errorf("%s is required, found an empty string.", field))	
		}
	}

	return wasError
}
