package helpers

import "github.com/go-playground/validator/v10"

type ErrorResponse struct {
	FailedField string
	Tag         string
	Value       string
}

func ValidateStruct(param any) []*ErrorResponse {
	var errors []*ErrorResponse
	validate := validator.New()
	err := validate.Struct(param)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			element := ErrorResponse{
				FailedField: err.StructNamespace(),
				Tag:         err.Tag(),
				Value:       err.Param(),
			}
			errors = append(errors, &element)
		}
	}
	return errors
}
