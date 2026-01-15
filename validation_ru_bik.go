// Copyright 2021-2026 Altessa Solutions Inc. All rights reserved.
// Use of this source code is governed by license that can be found in
// the LICENSE file.

package ozzo_rules

import (
	"regexp"

	"github.com/go-ozzo/ozzo-validation/v4"
)

// ErrRuBICInvalid is the error that returns when a value is not a valid Russian BIC number.
var ErrRuBICInvalid = validation.NewError("validation_ru_bik", "must be an 11-character code, starting with "+
	"prefix 04 (RU) and followed by 9-digits")

// RuBIC is a validation rule that checks if a value is a valid Russian BIC number.
func RuBIC() RuBICRule {
	return RuBICRule{err: ErrRuBICInvalid, condition: true}
}

// RuBICRule is a rule that checks if a value is a valid Russian BIC number.
type RuBICRule struct {
	err       validation.Error
	condition bool
}

// Validate checks if the given value is valid or not.
func (r RuBICRule) Validate(v any) error {
	if !r.condition {
		return nil
	}

	value, isNil := validation.Indirect(v)
	if isNil {
		return r.err
	}

	str, err := validation.EnsureString(value)
	if err != nil || !ruBicRx.MatchString(str) {
		return r.err
	}

	return nil
}

// When sets the condition that determines if the validation should be performed.
func (r RuBICRule) When(condition bool) RuBICRule {
	r.condition = condition
	return r
}

// Error sets the error message for the rule.
func (r RuBICRule) Error(message string) RuBICRule {
	r.err = r.err.SetMessage(message)
	return r
}

// ErrorObject sets the error struct for the rule.
func (r RuBICRule) ErrorObject(err validation.Error) RuBICRule {
	r.err = err
	return r
}

var ruBicRx = regexp.MustCompile(`^(04)[\d]{7}$`)
