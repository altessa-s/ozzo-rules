// Copyright 2021-2026 Altessa Solutions Inc. All rights reserved.
// Use of this source code is governed by license that can be found in
// the LICENSE file.

package ozzo_rules

import (
	"regexp"

	"github.com/go-ozzo/ozzo-validation/v4"
)

// ErrBirthdateInvalid is the error that returns when a value is not a valid birthdate
var ErrBirthdateInvalid = validation.NewError("validation_birthdate", "invalid birthdate")

var dateRx = regexp.MustCompile(`^\d{4}-[01]\d-[0-3]\d$`)

// Birthdate is a validation rule that checks if a value is a valid birthdate.
func Birthdate() BirthdateRule {
	return BirthdateRule{condition: true, err: ErrBirthdateInvalid}
}

// BirthdateRule is a rule that checks if a value is a valid birthdate
type BirthdateRule struct {
	err       validation.Error
	condition bool
}

// Validate checks if the given value is valid or not.
func (r BirthdateRule) Validate(v any) error {
	if !r.condition {
		return nil
	}

	value, isNil := validation.Indirect(v)
	if isNil {
		return r.err
	}

	str, err := validation.EnsureString(value)
	if err != nil || !dateRx.MatchString(str) {
		return r.err
	}
	return nil
}

// When sets the condition that determines if the validation should be performed.
func (r BirthdateRule) When(condition bool) BirthdateRule {
	r.condition = condition
	return r
}

// Error sets the error message for the rule.
func (r BirthdateRule) Error(message string) BirthdateRule {
	r.err = r.err.SetMessage(message)
	return r
}

// ErrorObject sets the error struct for the rule.
func (r BirthdateRule) ErrorObject(err validation.Error) BirthdateRule {
	r.err = err
	return r
}
