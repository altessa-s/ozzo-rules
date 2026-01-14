// Copyright 2021-2026 Altessa Solutions Inc. All rights reserved.
// Use of this source code is governed by license that can be found in
// the LICENSE file.

package ozzo_rules

import "github.com/go-ozzo/ozzo-validation/v4"

var (
	// ErrRequired is the error that returns when a value is required.
	ErrRequired = validation.NewError("validation_required", "cannot be blank")
)

// OneofRequired creates a new RequiredRule that checks if at least one of the given values is not empty.
func OneofRequired(val ...any) RequiredRule {
	return RequiredRule{condition: true, vals: val, err: ErrRequired}
}

// AllRequired creates a new RequiredRule that checks if all the given values are not empty.
func AllRequired(val ...any) RequiredRule {
	return RequiredRule{condition: true, vals: val, err: ErrRequired, allRequired: true}
}

type RequiredRule struct {
	condition   bool
	err         validation.Error
	vals        []any
	allRequired bool
}

func (r RequiredRule) valsIsiEmpty() bool {
	if len(r.vals) == 0 {
		return false
	}

	for _, val := range r.vals {
		if !validation.IsEmpty(val) {
			return false
		}
	}

	return true
}

// Validate checks if the given value is valid or not.
func (r RequiredRule) Validate(value any) error {
	if r.condition {
		value, _ = validation.Indirect(value)
		if !r.allRequired {
			if validation.IsEmpty(value) && r.valsIsiEmpty() {
				return r.err
			}
			return nil
		}

		if validation.IsEmpty(value) && !r.valsIsiEmpty() {
			return r.err
		}
	}
	return nil
}

// When sets the condition that determines if the validation should be performed.
func (r RequiredRule) When(condition bool) RequiredRule {
	r.condition = condition
	return r
}

// Error sets the error message for the rule.
func (r RequiredRule) Error(message string) RequiredRule {
	r.err = r.err.SetMessage(message)
	return r
}

// ErrorObject sets the error struct for the rule.
func (r RequiredRule) ErrorObject(err validation.Error) RequiredRule {
	r.err = err
	return r
}
