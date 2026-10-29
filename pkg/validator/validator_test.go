package validator_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"golang-starter-kit/pkg/validator"
)

type TestUserDTO struct {
	Name     string `json:"name" validate:"required,min=3"`
	Email    string `json:"email" validate:"required,email"`
	Age      int    `json:"age" validate:"min=18"`
	Optional string `json:"optional"`
}

func TestValidator_ValidStruct(t *testing.T) {
	v := validator.New()
	dto := TestUserDTO{
		Name:  "Budi Santoso",
		Email: "budi@example.com",
		Age:   25,
	}

	errs := v.Validate(dto)
	assert.Nil(t, errs)
}

func TestValidator_InvalidStruct(t *testing.T) {
	v := validator.New()
	dto := TestUserDTO{
		Name:  "Al",
		Email: "bukan-email",
		Age:   15,
	}

	errs := v.Validate(dto)
	assert.NotNil(t, errs)
	assert.Len(t, errs, 3)

	fields := map[string]string{}
	for _, e := range errs {
		fields[e.Field] = e.Message
	}

	assert.Contains(t, fields["name"], "must be at least 3")
	assert.Contains(t, fields["email"], "must be a valid email")
	assert.Contains(t, fields["age"], "must be at least 18")
}

func TestValidator_ValidateStructGlobal(t *testing.T) {
	dto := TestUserDTO{}
	errs := validator.ValidateStruct(dto)
	assert.NotNil(t, errs)
	assert.NotEmpty(t, errs)
}
