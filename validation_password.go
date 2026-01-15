// Copyright 2021-2026 Altessa Solutions Inc. All rights reserved.
// Use of this source code is governed by license that can be found in
// the LICENSE file.

package ozzo_rules

import (
	"unicode"

	"github.com/go-ozzo/ozzo-validation/v4"
)

// ErrPasswordInvalid is the error that returns when a value is not a valid password.
var ErrPasswordInvalid = validation.NewError("validation_password", "must be at least 8 "+
	"characters long, contain only ASCII symbols, at least one uppercase letter, "+
	"at least one number and at least one punctuation symbol")

// Password is a validation rule that checks if a value is a valid password.
// A value is a valid password if
// - length of at least 8 characters
// - contain at least one uppercase
// - contain at least one number
// - contain at least one punctuation symbol
func Password() PasswordRule {
	return PasswordRule{condition: true, err: ErrPasswordInvalid}
}

// PasswordRule is a rule that checks if a value is a valid password.
type PasswordRule struct {
	err       validation.Error
	condition bool
}

// Validate checks if the given value is valid or not.
func (r PasswordRule) Validate(v any) error {
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

	if str == "" || len(str) < 8 {
		return r.err
	}

	var (
		number bool
		upper  bool
		punct  bool
	)

	for _, char := range str {
		if char > unicode.MaxASCII { // MLM-168
			return r.err
		}
		switch {
		case unicode.IsNumber(char):
			number = true
		case unicode.IsUpper(char):
			upper = true
		case unicode.IsPunct(char) || char == ' ' || unicode.IsSymbol(char):
			punct = true
		}
	}

	if number && upper && punct {
		return nil
	}

	return r.err
}

// When sets the condition that determines if the validation should be performed.
func (r PasswordRule) When(condition bool) PasswordRule {
	r.condition = condition
	return r
}

// Error sets the error message for the rule.
func (r PasswordRule) Error(message string) PasswordRule {
	r.err = r.err.SetMessage(message)
	return r
}

// ErrorObject sets the error struct for the rule.
func (r PasswordRule) ErrorObject(err validation.Error) PasswordRule {
	r.err = err
	return r
}
