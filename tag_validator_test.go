package tag_validator

import (
	"fmt"
	"github.com/google/uuid"
	"testing"
)

func TestValidateStruct(t *testing.T) {
	type User struct {
		Id   string `validate:"uuid"`
		Name string `validate:"string,min=2,max=10"`
		Age  int    `validate:"number,min=18,max=20"`
	}
	user := User{Id: uuid.New().String(),Name: "y", Age: 5}
	fmt.Println("Errors:")
	for i, err := range ValidateStruct("", user) {
		fmt.Printf("\t%d. %s\n", i+1, err.Error())
	}
}
