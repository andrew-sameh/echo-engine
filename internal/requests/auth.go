package requests

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

// use a single instance of Validate, it caches struct info
var validate *validator.Validate

const (
	minPathLength = 8
)

type BasicAuth struct {
	Email    string `json:"email" validate:"required" example:"john.doe@example.com"`
	Password string `json:"password" validate:"required" example:"11111111"`
}

func (ba BasicAuth) Validate() error {
	validate = validator.New()
	err := validate.Struct(&ba)
	if err != nil {
		return ValidationErrorHandling(err)
	}
	return nil
}

type LoginRequest struct {
	BasicAuth
}

type RegisterRequest struct {
	BasicAuth
	FirstName string `json:"first_name" validate:"required" example:"John"`
	LastName  string `json:"last_name" validate:"required" example:"Doe"`
	Role      string `json:"role" validate:"required" example:"admin"`
	Username  string `json:"username" validate:"required" example:"johndoe"`
}

func (rr RegisterRequest) Validate() error {
	validate = validator.New()
	err := validate.Struct(&rr)
	if err != nil {
		return ValidationErrorHandling(err)
	}
	return nil
}

type RefreshRequest struct {
	Token string `json:"token" validate:"required" example:"refresh_token"`
}

func ValidationErrorHandling(err error) error {
	for _, err := range err.(validator.ValidationErrors) {

		fmt.Println(err.Namespace())
		fmt.Println(err.Field())
		fmt.Println(err.StructNamespace())
		fmt.Println(err.StructField())
		fmt.Println(err.Tag())
		fmt.Println(err.ActualTag())
		fmt.Println(err.Kind())
		fmt.Println(err.Type())
		fmt.Println(err.Value())
		fmt.Println(err.Param())
		fmt.Println()
	}

	// from here you can create your own error messages in whatever language you wish
	return err
}
