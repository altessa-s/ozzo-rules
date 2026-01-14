// Copyright 2021-2026 Altessa Solutions Inc. All rights reserved.
// Use of this source code is governed by license that can be found in
// the LICENSE file.

package ozzo_rules

import (
	"regexp"

	"github.com/go-ozzo/ozzo-validation/v4"
)

// ErrApnKeyIdInvalid is the error that returns when a value is not a valid apn key.
var ErrApnKeyIdInvalid = validation.NewError("validation_apn_key_id",
	"must be 10 characters long and may only contain alphanumeric characters")

// ApnKeyId is a validation rule that checks if a value is a valid apn key.
func ApnKeyId() ApnKeyIdRule {
	return ApnKeyIdRule{condition: true, err: ErrApnKeyIdInvalid}
}

// ApnKeyIdRule is a rule that checks if a value is a valid apn key.
type ApnKeyIdRule struct {
	err       validation.Error
	condition bool
}

// Validate checks if the given value is valid or not.
func (r ApnKeyIdRule) Validate(v any) error {
	if !r.condition {
		return nil
	}

	value, isNil := validation.Indirect(v)
	if isNil {
		return r.err
	}

	str, err := validation.EnsureString(value)
	if err != nil || !keyRx.MatchString(str) {
		return r.err
	}
	return nil
}

// When sets the condition that determines if the validation should be performed.
func (r ApnKeyIdRule) When(condition bool) ApnKeyIdRule {
	r.condition = condition
	return r
}

// Error sets the error message for the rule.
func (r ApnKeyIdRule) Error(message string) ApnKeyIdRule {
	r.err = r.err.SetMessage(message)
	return r
}

// ErrorObject sets the error struct for the rule.
func (r ApnKeyIdRule) ErrorObject(err validation.Error) ApnKeyIdRule {
	r.err = err
	return r
}

var keyRx = regexp.MustCompile(`^[A-Z0-9.\-]{10}$`)
