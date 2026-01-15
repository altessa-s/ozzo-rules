package main

import (
	"fmt"

	ozzo_rules "github.com/altessa-s/ozzo-rules"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

func main() {
	// Validate a country code
	err := validation.Validate("US", ozzo_rules.CountryCode2())
	if err != nil {
		fmt.Printf("Invalid country: %v\n", err)
	}

	// Validate a phone number
	err = validation.Validate("+12125551234", ozzo_rules.Phone("US"))
	if err != nil {
		fmt.Printf("Invalid phone: %v\n", err)
	}

	// Validate with custom error message
	err = validation.Validate("invalid",
		ozzo_rules.Username().Error("Username must be alphanumeric"))
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}
}
