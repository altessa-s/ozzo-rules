// Copyright 2021-2026 Altessa Solutions Inc. All rights reserved.
// Use of this source code is governed by license that can be found in
// the LICENSE file.

package ozzo_rules

import (
	"github.com/go-ozzo/ozzo-validation/v4"
	"github.com/nyaruka/phonenumbers"
)

// ErrPhoneInvalid is the error that returns when a value is not a valid phone number.
var ErrPhoneInvalid = validation.NewError("validation_phone", "invalid phone number")

// Phone is a validation rule that checks if a value is a valid phone number.
func Phone(countryCode string) PhoneRule {
	return PhoneRule{condition: true, countryCode: countryCode, err: ErrPhoneInvalid}
}

const DefaultCountryCode = "US"

// PhoneRule is a rule that checks if a value is a valid phone number.
type PhoneRule struct {
	err         validation.Error
	countryCode string
	condition   bool
}

// Validate checks if the given value is valid or not.
func (r PhoneRule) Validate(v any) error {
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

	if len(str) == 0 {
		return r.err
	}

	if str[0] != '+' {
		str = "+" + str
	}

	p, err := phonenumbers.Parse(str, r.countryCode)
	if err != nil {
		_ = r.err.SetMessage(err.Error()) //nolint:errcheck
		return r.err
	}

	if !phonenumbers.IsValidNumber(p) {
		return r.err
	}

	return nil
}

// When sets the condition that determines if the validation should be performed.
func (r PhoneRule) When(condition bool) PhoneRule {
	r.condition = condition
	return r
}

// Error sets the error message for the rule.
func (r PhoneRule) Error(message string) PhoneRule {
	r.err = r.err.SetMessage(message)
	return r
}

// ErrorObject sets the error struct for the rule.
func (r PhoneRule) ErrorObject(err validation.Error) PhoneRule {
	r.err = err
	return r
}
