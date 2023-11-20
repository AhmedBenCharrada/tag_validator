package tag_validator_test

import (
	validator "tag_validator"
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

type MultiError interface {
	Error() string
	Unwrap() []error
}

func TestValidate(t *testing.T) {
	assert.Panics(t, func() {
		_ = validator.New[string]()
	})

	val := validator.New[User](validator.CustomValidator{
		Tag: "text",
		Validator: func(_ interface{}, _ []string) error {
			return nil
		},
	})

	validUUID := "ba6516aa-3cb8-4592-b3cf-ba3ad9e176ae"

	user := User{Id: validUUID, Name: "name", LastName: "lastName", Age: 19, Height: 180, Weight: 75}
	err := val.Validate(user)

	assert.NoError(t, err)

	// invalid ID
	user = User{Id: "id", Name: "name", LastName: "lastName", Age: 19, Height: 180, Weight: 75}
	err = val.Validate(user)

	multiErr := new(MultiError)

	assert.Error(t, err)
	assert.ErrorAs(t, err, multiErr)
	assert.Equal(t, 1, len(err.(MultiError).Unwrap()))

	// invalid name
	user = User{Id: validUUID, Name: "name007", LastName: "lastName", Age: 19, Height: 180, Weight: 75}
	err = val.Validate(user)

	assert.Error(t, err)
	assert.ErrorAs(t, err, multiErr)
	assert.Equal(t, 1, len(err.(MultiError).Unwrap()))

	// name is too long
	user = User{Id: validUUID, Name: "abcdefghijklmnopqrstuvwxyz", LastName: "lastName", Age: 19, Height: 180, Weight: 75}
	err = val.Validate(user)

	assert.Error(t, err)
	assert.ErrorAs(t, err, multiErr)
	assert.Equal(t, 1, len(err.(MultiError).Unwrap()))
	// empty name
	user = User{Id: validUUID, Name: "", LastName: "lastName", Age: 19}
	err = val.Validate(user)

	assert.Error(t, err)
	assert.ErrorAs(t, err, multiErr)
	assert.Equal(t, 1, len(err.(MultiError).Unwrap()))

	// missing required lastName
	user = User{Id: validUUID, Name: "name", Age: 19, Height: 180, Weight: 75}
	err = val.Validate(user)

	assert.Error(t, err)
	assert.ErrorAs(t, err, multiErr)
	assert.Equal(t, 1, len(err.(MultiError).Unwrap()))

	// too old
	user = User{Id: validUUID, Name: "name", LastName: "lastName", Age: 99, Height: 180, Weight: 75}
	err = val.Validate(user)

	assert.Error(t, err)
	assert.ErrorAs(t, err, multiErr)
	assert.Equal(t, 1, len(err.(MultiError).Unwrap()))

	// terribly wrong
	user = User{Id: "ba3ad9e176ae", Name: "y", LastName: "lastName", Age: 8}
	err = val.Validate(user)

	assert.Error(t, err)
	assert.ErrorAs(t, err, multiErr)
	assert.Equal(t, 3, len(err.(MultiError).Unwrap()))
}
