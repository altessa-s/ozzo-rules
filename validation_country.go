// Copyright 2021-2026 Altessa Solutions Inc. All rights reserved.
// Use of this source code is governed by license that can be found in
// the LICENSE file.

package ozzo_rules

import (
	"strings"

	"github.com/asaskevich/govalidator"
	"github.com/go-ozzo/ozzo-validation/v4"
)

// ErrCountryCode2Invalid is the error that returns when a value is not a valid two-letter country code.
var ErrCountryCode2Invalid = validation.NewError("validation_country_code2", "must be a valid two-letter country code")

// CountryCode2 is a validation rule that checks if a value is a valid ISO 3166-1 alpha-2 country code.
func CountryCode2() CountryCode2Rule {
	return CountryCode2Rule{condition: true, err: ErrCountryCode2Invalid}
}

// CountryCode2Rule is a rule that checks if a value is a valid ISO 3166-1 alpha-2 country code.
type CountryCode2Rule struct {
	err       validation.Error
	condition bool
}

// Validate checks if the given value is valid or not.
func (r CountryCode2Rule) Validate(v any) error {
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

	if str == "" || !govalidator.IsISO3166Alpha2(strings.ToUpper(str)) {
		return r.err
	}

	return nil
}

// When sets the condition that determines if the validation should be performed.
func (r CountryCode2Rule) When(condition bool) CountryCode2Rule {
	r.condition = condition
	return r
}

// Error sets the error message for the rule.
func (r CountryCode2Rule) Error(message string) CountryCode2Rule {
	r.err = r.err.SetMessage(message)
	return r
}

// ErrorObject sets the error struct for the rule.
func (r CountryCode2Rule) ErrorObject(err validation.Error) CountryCode2Rule {
	r.err = err
	return r
}
