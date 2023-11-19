package tag_validator

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type User struct {
	Id       string `validate:"uuid"`
	Name     string `validate:"string,min=2,max=10,pattern=^[a-zA-Z]+$"`
	LastName string `validate:"string,required,min=2,max=10"`
	Age      int    `validate:"number,min=18,max=20"`
	Height   int    `validate:"ignore,min=165,max=200"`
	Weight   int    `validate:"-"`
}

func TestValidateStruct(t *testing.T) {
	validUUID := "ba6516aa-3cb8-4592-b3cf-ba3ad9e176ae"

	user := User{Id: validUUID, Name: "name", LastName: "lastName", Age: 19, Height: 180, Weight: 75}
	errs := ValidateStruct(user, CustomValidator{
		Tag: "text",
		Validator: func(_ interface{}, _ []string) error {
			return nil
		},
	})

	assert.Empty(t, errs)

	// invalid ID
	user = User{Id: "id", Name: "name", LastName: "lastName", Age: 19, Height: 180, Weight: 75}
	errs = ValidateStruct(user)

	assert.NotEmpty(t, errs)
	assert.Equal(t, 1, len(errs))

	// invalid name
	user = User{Id: validUUID, Name: "name007", LastName: "lastName", Age: 19, Height: 180, Weight: 75}
	errs = ValidateStruct(user)

	assert.NotEmpty(t, errs)
	assert.Equal(t, 1, len(errs))

	// name is too long
	user = User{Id: validUUID, Name: "abcdefghijklmnopqrstuvwxyz", LastName: "lastName", Age: 19, Height: 180, Weight: 75}
	errs = ValidateStruct(user)

	assert.NotEmpty(t, errs)
	assert.Equal(t, 1, len(errs))

	// empty name
	user = User{Id: validUUID, Name: "", LastName: "lastName", Age: 19}
	errs = ValidateStruct(user)

	assert.NotEmpty(t, errs)
	assert.Equal(t, 1, len(errs))

	// missing required lastName
	user = User{Id: validUUID, Name: "name", Age: 19, Height: 180, Weight: 75}
	errs = ValidateStruct(user)

	assert.NotEmpty(t, errs)
	assert.Equal(t, 1, len(errs))

	// too old
	user = User{Id: validUUID, Name: "name", LastName: "lastName", Age: 99, Height: 180, Weight: 75}
	errs = ValidateStruct(user)

	assert.NotEmpty(t, errs)
	assert.Equal(t, 1, len(errs))

	// terribly wrong
	user = User{Id: "ba3ad9e176ae", Name: "y", LastName: "lastName", Age: 8}
	errs = ValidateStruct(user)

	assert.NotEmpty(t, errs)
	assert.Equal(t, 3, len(errs))
}
