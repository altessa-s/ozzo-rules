package main

import (
	"github.com/go-ozzo/ozzo-validation/v4/is"

	ozzo_rules "github.com/altessa-s/ozzo-rules"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

// Struct validation with multiple rules
type User struct {
	Username    string
	Password    string
	Email       string
	Country     string
	Language    string
	Phone       string
	Birthdate   string
	Timezone    string
	ProfileSlug string
}

func (u User) Validate() error {
	return validation.ValidateStruct(&u,
		// Required fields with custom messages
		validation.Field(&u.Username,
			validation.Required.Error("Username is required"),
			ozzo_rules.Username()),

		// Conditional validation
		validation.Field(&u.Password,
			validation.Required.When(u.Password != ""),
			ozzo_rules.Password()),

		// Standard validation rules
		validation.Field(&u.Email,
			validation.Required,
			is.Email),

		// Geographic validation
		validation.Field(&u.Country,
			ozzo_rules.CountryCode2()),
		validation.Field(&u.Language,
			ozzo_rules.LangCode2()),

		// Phone validation with country context
		validation.Field(&u.Phone,
			ozzo_rules.Phone(u.Country).When(u.Phone != "")),

		// Date and timezone validation
		validation.Field(&u.Birthdate,
			ozzo_rules.Birthdate()),
		validation.Field(&u.Timezone,
			ozzo_rules.Timezone()),

		// Slug validation for SEO-friendly URLs
		validation.Field(&u.ProfileSlug,
			ozzo_rules.Slug()),
	)
}
