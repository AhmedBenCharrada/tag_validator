package tag_validator

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateStruct(t *testing.T) {
	type User struct {
		Id     string `validate:"uuid"`
		Name   string `validate:"string,min=2,max=10"`
		Age    int    `validate:"number,min=18,max=20"`
		Height int    `validate:"ignore,min=165,max=200"`
		Weight int    `validate:"-"`
	}

	user := User{Id: "ba6516aa-3cb8-4592-b3cf-ba3ad9e176ae", Name: "name", Age: 19, Height: 180, Weight: 75}
	errs := ValidateStruct("", user)

	assert.Empty(t, errs)

	user = User{Id: "ba3ad9e176ae", Name: "y", Age: 8}
	errs = ValidateStruct("", user)

	assert.NotEmpty(t, errs)
	assert.Equal(t, 3, len(errs))

	user = User{Id: "ba6516aa-3cb8-4592-b3cf-ba3ad9e176ae", Name: "name", Age: 29}
	errs = ValidateStruct("", user)

	assert.NotEmpty(t, errs)
	assert.Equal(t, 1, len(errs))

	user = User{Id: "ba6516aa-3cb8-4592-b3cf-ba3ad9e176ae", Name: "", Age: 19}
	errs = ValidateStruct("", user)

	assert.NotEmpty(t, errs)
	assert.Equal(t, 1, len(errs))

	user = User{Id: "ba6516aa-3cb8-4592-b3cf-ba3ad9e176ae", Name: "ba6516aa-3cb8-4592-b3cf-ba3ad9e176ae", Age: 19}
	errs = ValidateStruct("", user)

	assert.NotEmpty(t, errs)
	assert.Equal(t, 1, len(errs))
}
