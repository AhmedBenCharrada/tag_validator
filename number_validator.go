package tag_validator

import "fmt"

type numberValidator struct {
	min int
	max int
}

func (v numberValidator) validate(cid string, val interface{}) error{
	num := val.(int)
	if num < v.min {
		return fmt.Errorf("should be greater than %v", v.min)
	}
	// only max is defined (!= 0) and nun > v.max return an error
	if v.max >= v.min && num > v.max {
		return fmt.Errorf("should be less than %v", v.max)
	}
	return nil
}
