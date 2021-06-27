package tag_validator

import (
	"errors"
	"regexp"
)

type uuidValidator struct {
	min int
	max int
}

func (v uuidValidator) validate(cid string, val interface{}) error{
	id := val.(string)
	r := regexp.MustCompile("" +
		"^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-4[a-fA-F0-9]{3}-[8|9|aA|bB][a-fA-F0-9]{3}-[a-fA-F0-9]{12}$",
	)
	if !r.MatchString(id) {
		return errors.New("invalid uuid")
	}
	return nil
}
