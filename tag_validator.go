package tag_validator

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strings"
)

// The validator tag name
//
// Fields to be validated need to be tagged with this tag.
// For example having the struct:
//
//	type User struct {
//		ID   	 string `validate:"uuid"`
//		Name     string `validate:"string,min=2,max=10,pattern=^[a-zA-Z]+$"`
//		LastName string `validate:"string,required,min=2,max=10"`
//		Age  	 int    `validate:"number,min=18,max=20"`
//	}
//
// a uuid validator will be applied on the ID field;
// a string validator will be applied to the Name field and
// a number validator for Age field where min and max are parameter being passed to the validator
//
// Validator parameters should be added according to the validator type
//   - number:
//     min={value}: define the minimum allowed value.
//     max={value}: define the maximum allowed value.
//   - string:
//     min={value}: define the minimum string length.
//     max={value}: define the maximum string length.
//     required: define whether this filed is required and cannot be empty.
//     pattern: define the patter the value should match.
const tagName = "validate"

type number interface {
	int | int64
}

// Validator ..
type Validator func(interface{}, []string) error

// CustomValidator ..
type CustomValidator struct {
	Tag       string
	Validator Validator
}

// ValidateStruct traverse all the struct fields and validates attributes marked to be validated
func ValidateStruct(s interface{}, opts ...CustomValidator) []error {
	if reflect.TypeOf(s).Kind() != reflect.Struct {
		panic("input should be a struct")
	}

	validators := map[string]Validator{
		"number": numberValidator,
		"string": stringValidator,
		"uuid":   uuidValidator,
	}

	for _, opt := range opts {
		validators[opt.Tag] = opt.Validator
	}

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
		err := check(validators, tag, v.Field(i).Interface())
		// Append error to results
		if err != nil {
			errs = append(errs, fmt.Errorf("%s %s", v.Type().Field(i).Name, err.Error()))
		}
	}
	if len(errs) == 0 {
		return nil
	}
	return errs
}

func numberValidator(val interface{}, args []string) error {
	var min int
	var max int
	for _, arg := range args[1:] {
		_, _ = fmt.Sscanf(arg, "max=%v", &max)
		_, _ = fmt.Sscanf(arg, "min=%v", &min)
	}

	return validateNumber(min, max, val.(int))
}

func stringValidator(val interface{}, args []string) error {
	var min int
	var max int
	var required bool
	var pattern string

	for _, arg := range args[1:] {
		_, _ = fmt.Sscanf(arg, "max=%v", &max)
		_, _ = fmt.Sscanf(arg, "min=%v", &min)
		_, _ = fmt.Sscanf(arg, "pattern=%s", &pattern)

		if !required {
			required = arg == "required"
		}
	}

	return validateString(min, max, required, pattern, val.(string))
}

func uuidValidator(val interface{}, _ []string) error {
	return validateUUID(val.(string))
}

// gets the appropriate validator based on field type
func check(validators map[string]Validator, tag string, item interface{}) error {
	args := strings.Split(tag, ",")
	val, ok := validators[args[0]]
	if !ok || val == nil {
		return nil
	}

	return val(item, args)
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

func validateString(min int, max int, required bool, pattern string, val string) error {
	num := len(val)

	if num == 0 && required {
		return fmt.Errorf("should not be empty")
	}

	if num < min {
		return fmt.Errorf("should be greater than %v", min)
	}

	// only max is defined (!= 0) and nun > v.max return an error
	if max >= min && num > max {
		return fmt.Errorf("should be less than %v", max)
	}

	if len(pattern) > 0 {
		return validateRegex(val, pattern)
	}

	return nil
}

func validateUUID(id string) error {
	return validateRegex(
		id,
		"^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-4[a-fA-F0-9]{3}-[8|9|aA|bB][a-fA-F0-9]{3}-[a-fA-F0-9]{12}$",
	)
}

func validateRegex(val, pattern string) error {
	r := regexp.MustCompile("\\b" + pattern + "\\b")

	if !r.MatchString(val) {
		return errors.New("pattern does not match")
	}

	return nil
}
