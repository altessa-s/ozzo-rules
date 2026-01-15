// Copyright 2021-2026 Altessa Solutions Inc. All rights reserved.
// Use of this source code is governed by license that can be found in
// the LICENSE file.

package ozzo_rules

import (
	"regexp"
	"strings"

	"github.com/adrg/postcode"
	"github.com/go-ozzo/ozzo-validation/v4"
)

// ErrZipCodeInvalid is the error that returns when a value is not a valid zip code
var ErrZipCodeInvalid = validation.NewError("validation_zip_code", "invalid zip code")

// ZipCode is a validation rule that checks if a value is a valid zip code.
func ZipCode() ZipCodeRule {
	return ZipCodeRule{condition: true, err: ErrZipCodeInvalid}
}

// ZipCodeRule is a rule that checks if a value is a valid zip code
type ZipCodeRule struct {
	err       validation.Error
	condition bool
}

// Validate checks if the given value is valid or not.
func (r ZipCodeRule) Validate(v any) error {
	if !r.condition {
		return nil
	}

	value, isNil := validation.Indirect(v)
	if isNil {
		return r.err
	}

	str, err := validation.EnsureString(value)
	if err != nil {
		return r.err
	}

	// Trim spaces for validation
	trimmed := strings.TrimSpace(str)
	if trimmed != str {
		// If the original string had leading/trailing spaces, it's invalid
		return r.err
	}

	// Empty string is invalid
	if str == "" {
		return r.err
	}

	// Special handling for UK postcodes - convert to uppercase
	if isUKPostcodeFormat(str) {
		str = strings.ToUpper(str)
	}

	if err := postcode.Validate(str); err != nil {
		return r.err
	}
	return nil
}

// When sets the condition that determines if the validation should be performed.
func (r ZipCodeRule) When(condition bool) ZipCodeRule {
	r.condition = condition
	return r
}

// Error sets the error message for the rule.
func (r ZipCodeRule) Error(message string) ZipCodeRule {
	r.err = r.err.SetMessage(message)
	return r
}

// ErrorObject sets the error struct for the rule.
func (r ZipCodeRule) ErrorObject(err validation.Error) ZipCodeRule {
	r.err = err
	return r
}

// isUKPostcodeFormat checks if the string looks like a UK postcode
func isUKPostcodeFormat(s string) bool {
	// Simple check for UK postcode format (contains letters and numbers)
	ukPattern := regexp.MustCompile(`^[A-Za-z][A-Za-z0-9].*\s+\d[A-Za-z]{2}$`)
	return ukPattern.MatchString(s)
}
