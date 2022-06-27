package tag_validator

import (
	"fmt"
	"reflect"
	"strings"
)

// the validator tag name
// fields to be validated need to be tagged with this tag
const tagName = "validate"

// validator an interface that every tag validator should implement
type validator interface {
	// Validate checks whether the validated field is valid or not
	validate(in interface{}) error
}

// ValidateStruct travers all the struct fields and validates attributes marked to be validated
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
		validator := getValidator(tag)
		if validator == nil {
			continue
		}
		// Perform validation
		err := validator.validate(v.Field(i).Interface())
		// Append error to results
		if err != nil {
			errs = append(errs, fmt.Errorf("%s %s", v.Type().Field(i).Name, err.Error()))
		}
	}
	return errs
}

// gets the appropriate validator based on field type
func getValidator(tag string) validator {
	args := strings.Split(tag, ",")
	switch args[0] {
	case "number":
		validator := numberValidator{}
		_, _ = fmt.Sscanf(strings.Join(args[1:], ","), "min=%d,max=%d", &validator.min, &validator.max)
		return validator
	case "string":
		validator := stringValidator{}
		_, _ = fmt.Sscanf(strings.Join(args[1:], ","), "min=%d,max=%d", &validator.min, &validator.max)
		return validator
	case "uuid":
		return uuidValidator{}

	}
	return nil
}
