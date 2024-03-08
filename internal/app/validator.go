package app

import "github.com/go-playground/validator/v10"

type ApiValidator struct {
	Validator *validator.Validate
}

type ValidationError struct {
	Namespace       string `json:"namespace,omitempty"`
	Field           string `json:"field"`
	StructNamespace string `json:"structNamespace,omitempty"`
	StructField     string `json:"structField,omitempty"`
	Tag             string `json:"tag"`
	ActualTag       string `json:"actualTag,omitempty"`
	Kind            string `json:"kind,omitempty"`
	Type            string `json:"type,omitempty"`
	Value           string `json:"value,omitempty"`
	Param           string `json:"param,omitempty"`
	Message         string `json:"message,omitempty"`
}

func InitValidator() *validator.Validate {
	validator := validator.New();

	// here you can register custom tag handlers

	return validator
}

func (av *ApiValidator) Validate(i interface{}) error {
  if err := av.Validator.Struct(i); err != nil {
    // Optionally, you could return the error to give each route more control over the status code
		return err
    // return echo.NewHTTPError(http.StatusBadRequest, err.Error())
  }
  return nil
}

