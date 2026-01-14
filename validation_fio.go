// Copyright 2021-2026 Altessa Solutions Inc. All rights reserved.
// Use of this source code is governed by license that can be found in
// the LICENSE file.

package ozzo_rules

import (
	"regexp"

	"github.com/go-ozzo/ozzo-validation/v4"
)

// ErrFirstNameInvalid is the error that returns when a value is not a valid first name.
var ErrFirstNameInvalid = validation.NewError("validation_first_name_invalid",
	"First name must start with a capital letter and be 2-30 characters long, optionally with a hyphen.")

// ErrMiddleNameInvalid is the error that returns when a value is not a valid middle name.
var ErrMiddleNameInvalid = validation.NewError("validation_middle_name_invalid",
	"Middle name must start with a capital letter and be 5-30 characters long.")

// ErrLastNameInvalid is the error that returns when a value is not a valid last name.
var ErrLastNameInvalid = validation.NewError("validation_last_name_invalid",
	"Last name must start with a capital letter and be 2-50 characters long, optionally with a hyphen.")

// RussianFirstName is a validation rule that checks if a value is a valid first name in Russian.
func RussianFirstName() RussianNameRule {
	return RussianNameRule{err: ErrFirstNameInvalid, condition: true, rx: firstNameRx}
}

// RussianMiddleName is a validation rule that checks if a value is a valid middle name in Russian.
func RussianMiddleName() RussianNameRule {
	return RussianNameRule{err: ErrMiddleNameInvalid, condition: true, rx: middleNameRx}
}

// RussianLastName is a validation rule that checks if a value is a valid last name in Russian.
func RussianLastName() RussianNameRule {
	return RussianNameRule{err: ErrLastNameInvalid, condition: true, rx: lastNameRegex}
}

// RussianNameRule is a rule that checks if a value is a valid name.
type RussianNameRule struct {
	err       validation.Error
	condition bool
	rx        *regexp.Regexp
}

// Validate checks if the given value is valid or not.
func (r RussianNameRule) Validate(v any) error {
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

	if !r.rx.MatchString(str) {
		return r.err
	}

	return nil
}

// When sets the condition that determines if the validation should be performed.
func (r RussianNameRule) When(condition bool) RussianNameRule {
	r.condition = condition
	return r
}

// Error sets the error message for the rule.
func (r RussianNameRule) Error(message string) RussianNameRule {
	r.err = r.err.SetMessage(message)
	return r
}

// ErrorObject sets the error struct for the rule.
func (r RussianNameRule) ErrorObject(err validation.Error) RussianNameRule {
	r.err = err
	return r
}

var firstNameRx = regexp.MustCompile(`^[А-ЯЁ][а-яё]+(-[А-ЯЁ][а-яё]+)*$`)
var middleNameRx = regexp.MustCompile(`^[А-ЯЁ][а-яё]+(вич|вна|ич|ична|оглы|кызы)$`)
var lastNameRegex = regexp.MustCompile(`^[А-ЯЁ][а-яё]+(-[А-ЯЁ][а-яё]+)*$`)
