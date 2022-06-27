package tag_validator

import "fmt"

type stringValidator struct {
	min int
	max int
}

func (v stringValidator) validate(val interface{}) error {
	num := len(val.(string))
	if num == 0 {
		return fmt.Errorf("should not be empty")
	}
	if num < v.min {
		return fmt.Errorf("should be greater than %v", v.min)
	}
	// only max is defined (!= 0) and nun > v.max return an error
	if v.max >= v.min && num > v.max {
		return fmt.Errorf("should be less than %v", v.max)
	}
	return nil
}
