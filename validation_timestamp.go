// Copyright 2021-2026 Altessa Solutions Inc. All rights reserved.
// Use of this source code is governed by license that can be found in
// the LICENSE file.

package ozzo_rules

import (
	"github.com/go-ozzo/ozzo-validation/v4"
)

// ErrTimestampInvalid is the error that returns when a value is not a valid timestamp.
var ErrTimestampInvalid = validation.NewError("validation_timestamp", "invalid timestamp")

// Timestamp is a validation rule that checks if a value is a valid timestamp.
func Timestamp(minTms, maxTms int64) *TimestampRule {
	return &TimestampRule{condition: true, min: minTms, max: maxTms, err: ErrTimestampInvalid}
}

// TimestampRule is a rule that checks if a value is a valid timestamp.
type TimestampRule struct {
	err       validation.Error
	condition bool
	min       int64
	max       int64
}

// Validate checks if the given value is valid or not.
func (r TimestampRule) Validate(v any) error {
	if !r.condition {
		return nil
	}

	value, isNil := validation.Indirect(v)
	if isNil {
		return r.err
	}

	timestamp, ok := value.(int64)
	if !ok || timestamp < 0 {
		return r.err
	}

	if r.min > 0 && timestamp < r.min {
		return r.err
	}

	if r.max > 0 && timestamp > r.max {
		return r.err
	}

	return nil
}

// When sets the condition that determines if the validation should be performed.
func (r TimestampRule) When(condition bool) TimestampRule {
	r.condition = condition
	return r
}

// Error sets the error message for the rule.
func (r TimestampRule) Error(message string) TimestampRule {
	r.err = r.err.SetMessage(message)
	return r
}

// ErrorObject sets the error struct for the rule.
func (r TimestampRule) ErrorObject(err validation.Error) TimestampRule {
	r.err = err
	return r
}
