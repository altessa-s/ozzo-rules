// Copyright 2021-2026 Altessa Solutions Inc. All rights reserved.
// Use of this source code is governed by license that can be found in
// the LICENSE file.

package ozzo_rules

import (
	"regexp"

	"github.com/go-ozzo/ozzo-validation/v4"
)

// ErrUsernameInvalid is the error that returns when a value is not a valid username.
var ErrUsernameInvalid = validation.NewError("validation_username",
	"can only contain latin letters, numbers, dashes and underscores.")

// Username is a validation rule that checks if a value is a valid username
func Username() UsernameRule {
	return UsernameRule{condition: true, err: ErrUsernameInvalid}
}

// UsernameRule is a rule that checks if a value is a valid username
type UsernameRule struct {
	err       validation.Error
	condition bool
}

// Validate checks if the given value is valid or not.
func (r UsernameRule) Validate(v any) error {
	if !r.condition {
		return nil
	}

	value, isNil := validation.Indirect(v)
	if isNil {
		return r.err
	}

	str, err := validation.EnsureString(value)
	if err != nil || !usernameRx.MatchString(str) {
		return r.err
	}
	return nil
}

// When sets the condition that determines if the validation should be performed.
func (r UsernameRule) When(condition bool) UsernameRule {
	r.condition = condition
	return r
}

// Error sets the error message for the rule.
func (r UsernameRule) Error(message string) UsernameRule {
	r.err = r.err.SetMessage(message)
	return r
}

// ErrorObject sets the error struct for the rule.
func (r UsernameRule) ErrorObject(err validation.Error) UsernameRule {
	r.err = err
	return r
}

var usernameRx = regexp.MustCompile(`^[a-zA-Z0-9_\-]{1,56}$`)
