package utils_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/kaleb-white/letthemknow/message-server/models"
	"github.com/kaleb-white/letthemknow/message-server/utils"
)

func TestCheckFieldEqualityNoErrors(t *testing.T) {
	c1 := models.Contact{
		Id: 123,
		Phone: 456,
	}
	c2 := models.Contact{
		Id: 123,
		Phone: 456,
	}

	fields := []string{"Id", "Phone"}
	errs := []error{}
	wasError := utils.CheckFieldEquality(c1, c2, &fields, &errs)

	if len(errs) > 0 || wasError {
		t.Errorf("Should've gotten 0 errors, received more than 0")
	}
}

func TestCheckFieldEqualityNoErrorsWithTimeField(t *testing.T) {
	now := time.Now()
	c1 := models.Contact{
		Id: 789,
		Phone: 101112,
		LastUpdatedAtDateTime: now,
	}
	c2 := models.Contact{
		Id: 789,
		Phone: 101112,
		LastUpdatedAtDateTime: now,
	}

	fields := []string{"Id", "Phone", "LastUpdatedAtDateTime"}
	errs := []error{}
	wasError := utils.CheckFieldEquality(c1, c2, &fields, &errs)

	if len(errs) > 0 || wasError {
		t.Errorf("Should've gotten 0 errors, received more than 0")
		fmt.Println(utils.PrettyPrintErrors(errs))
	}
}

func TestCheckFieldEqualityTimeWrong(t *testing.T) {
	now := time.Now()
	c1 := models.Contact{
		Id: 789,
		Phone: 101112,
		LastUpdatedAtDateTime: now,
	}
	c2 := models.Contact{
		Id: 789,
		Phone: 101112,
		LastUpdatedAtDateTime: time.Date(2025, time.January, 0, 0, 0, 0, 0, time.UTC),
	}

	fields := []string{"Id", "Phone", "LastUpdatedAtDateTime"}
	errs := []error{}
	wasError := utils.CheckFieldEquality(c1, c2, &fields, &errs)

	if len(errs) != 1 || !wasError {
		t.Errorf("Should've gotten 1 error, received %d", len(errs))
		fmt.Println(utils.PrettyPrintErrors(errs))
	}
}

type HasSlice struct {
	S []string
}

func TestCheckFieldEqualitySlicePass(t *testing.T) {
	c1 := HasSlice{
		S: []string{"a", "b"},
	}
	c2 := HasSlice{
		S: []string{"a", "b"},
	}

	fields := []string{"S"}
	errs := []error{}
	wasError := utils.CheckFieldEquality(c1, c2, &fields, &errs)

	if len(errs) != 0 || wasError {
		t.Errorf("Should've gotten 0 errors, received %d", len(errs))
		fmt.Println(utils.PrettyPrintErrors(errs))
	}
}

func TestCheckFieldEqualitySliceFail(t *testing.T) {
	c1 := HasSlice{
		S: []string{"a", "c"},
	}
	c2 := HasSlice{
		S: []string{"a", "b"},
	}

	fields := []string{"S"}
	errs := []error{}
	wasError := utils.CheckFieldEquality(c1, c2, &fields, &errs)

	if len(errs) != 1 || !wasError {
		t.Errorf("Should've gotten 1 error, received %d", len(errs))
		fmt.Println(utils.PrettyPrintErrors(errs))
	}
}

func TestCheckFieldEqualityDifferentStructs(t *testing.T) {
	c1 := models.Contact{
		Id: 123,
		Phone: 123,
	}
	c2 := models.ContactDetail{
		Id: 123,
		Phone: 123,
	}

	fields := []string{"Id", "Phone"}
	errs := []error{}
	wasError := utils.CheckFieldEquality(c1, c2, &fields, &errs)

	if len(errs) != 1 || !wasError {
		t.Errorf("Should've gotten 1 error, received %d", len(errs))
	}
}

func TestCheckFieldEqualityShouldPanicNonexistent(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Should've panic, defer fired without error to recover")
		}
	}()
	c1 := models.Contact{
		Id: 123,
		Phone: 123,
	}
	c2 := models.Contact{
		Id: 123,
		Phone: 123,
	}

	fields := []string{"Abc"}
	errs := []error{}
	utils.CheckFieldEquality(c1, c2, &fields, &errs)
}

func TestCheckRequiredFieldsArentDefaultNoErrors(t *testing.T) {
	c := models.Contact{
		Id: 123,
		Phone: 123,
	}

	errs := []error{}
	fields := []string{"Id", "Phone"}

	utils.CheckRequiredFieldsArentDefault(c, &fields, &errs)

	if len(errs) > 0 {
		t.Errorf("Should've gotten 0 errors, received more than 0")
		fmt.Println(utils.PrettyPrintErrors(errs))
	}	
}

func TestCheckRequiredFieldsArentDefaultSomeErrors(t *testing.T) {
	c := models.Contact{
		Id: 123,
		Phone: 123,
	}

	errs := []error{}
	fields := []string{"Id", "Phone", "Phone1"}

	utils.CheckRequiredFieldsArentDefault(c, &fields, &errs)

	if len(errs) != 1 {
		t.Errorf("Should've gotten 1 error, received %d.", len(errs))
		fmt.Println(utils.PrettyPrintErrors(errs))
	}	
}

