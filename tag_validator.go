package tag_validator

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strings"
)

// the validator tag name
// fields to be validated need to be tagged with this tag
const tagName = "validate"

type number interface {
	int | int64
}

// ValidateStruct traverse all the struct fields and validates attributes marked to be validated
func ValidateStruct(cid string, s interface{}) []error {
	// ValueOf returns a Value representing the run-time data
	v := reflect.ValueOf(s)
	errs := make([]error, 0, v.NumField())

	for i := 0; i < v.NumField(); i++ {
		// Get the field tag value
		tag := v.Type().Field(i).Tag.Get(tagName)
		// Skip if tag is not defined or ignored
		if tag == "" || tag == "-" {
			continue
		}
		// Get a validator that corresponds to a tag
		err := check(tag, v.Field(i).Interface())
		// Append error to results
		if err != nil {
			errs = append(errs, fmt.Errorf("%s %s", v.Type().Field(i).Name, err.Error()))
		}
	}
	return errs
}

// gets the appropriate validator based on field type
func check(tag string, item interface{}) error {
	args := strings.Split(tag, ",")
	switch args[0] {
	case "number":
		var min int
		var max int
		_, _ = fmt.Sscanf(strings.Join(args[1:], ","), "min=%d,max=%d", &min, &max)

		return validateNumber(min, max, item.(int))
	case "string":
		var min int
		var max int
		_, _ = fmt.Sscanf(strings.Join(args[1:], ","), "min=%d,max=%d", &min, &max)

		return validateString(min, max, item.(string))
	case "uuid":
		return validateUUID(item.(string))
	}

	return nil
}

func validateNumber[T number](min, max, val T) error {
	if val < min {
		return fmt.Errorf("should be greater than %v", min)
	}
	// only max is defined (!= 0) and nun > v.max return an error
	if max >= min && val > max {
		return fmt.Errorf("should be less than %v", max)
	}

	return nil
}

func validateString(min int, max int, val string) error {
	num := len(val)
	if num == 0 {
		return fmt.Errorf("should not be empty")
	}

	if num < min {
		return fmt.Errorf("should be greater than %v", min)
	}

	// only max is defined (!= 0) and nun > v.max return an error
	if max >= min && num > max {
		return fmt.Errorf("should be less than %v", max)
	}

	return nil
}

func validateUUID(id string) error {
	r := regexp.MustCompile("" +
		"^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-4[a-fA-F0-9]{3}-[8|9|aA|bB][a-fA-F0-9]{3}-[a-fA-F0-9]{12}$",
	)
	if !r.MatchString(id) {
		return errors.New("invalid uuid")
	}
	return nil
}
