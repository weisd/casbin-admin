package handlers

import validator "gopkg.in/go-playground/validator.v9"

var validate *validator.Validate

// NewValidate NewValidate
func NewValidate() {
	validate = validator.New()
}

// ValidateStruct ValidateStruct
func ValidateStruct(i interface{}) error {
	return validate.Struct(i)
}
