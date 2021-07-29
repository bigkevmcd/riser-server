package model

import (
	validation "github.com/go-ozzo/ozzo-validation/v3"
	"github.com/google/uuid"
)

type AppName string

// App is an app within a namespace.
type App struct {
	Id        uuid.UUID     `json:"id"`
	Name      AppName       `json:"name"`
	Namespace NamespaceName `json:"namespace"`
}

// Validate performs some checks and returns an error if the App is not valid.
func (v App) Validate() error {
	return validation.ValidateStruct(&v,
		validation.Field(&v.Name),
		validation.Field(&v.Namespace))
}

// NewApp is an app creation request received in the API endpoint.
type NewApp struct {
	Name      AppName       `json:"name"`
	Namespace NamespaceName `json:"namespace"`
}

// Validate performs some checks and returns an error if the NewApp is not
// valid.
func (v NewApp) Validate() error {
	return validation.ValidateStruct(&v,
		validation.Field(&v.Name),
		validation.Field(&v.Namespace))
}

func (v AppName) Validate() error {
	return validation.Validate(string(v), RulesAppName()...)
}
