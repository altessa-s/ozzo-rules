// Copyright 2021-2026 Altessa Solutions Inc. All rights reserved.
// Use of this source code is governed by license that can be found in
// the LICENSE file.

package ozzo_rules

import (
	"fmt"

	"github.com/go-ozzo/ozzo-validation/v4"
)

const listMaxLimit = 1000

// ErrListLimitInvalid is the error that returns when a value is not a valid list limit.
var ErrListLimitInvalid = validation.NewError("validation_list_limit",
	fmt.Sprintf("must be greater than 0 and less than %d", listMaxLimit))

// ListLimit creates a validation rule that checks if a value is a valid list limit.
func ListLimit() ListLimitRule {
	return ListLimitRule{err: ErrListLimitInvalid, condition: true}
}

// ListLimitRule is a rule that checks if a value is a valid list limit.
type ListLimitRule struct {
	err       validation.Error
	condition bool
}

// Validate checks if the given value is valid or not.
func (r ListLimitRule) Validate(v any) error {
	if !r.condition {
		return nil
	}

	value, isNil := validation.Indirect(v)
	if isNil {
		return r.err
	}

	if l, ok := value.(int64); ok && l > 0 && l <= listMaxLimit {
		return nil
	}

	return r.err
}

// When sets the condition that determines if the validation should be performed.
func (r ListLimitRule) When(condition bool) ListLimitRule {
	r.condition = condition
	return r
}

// Error sets the error message for the rule.
func (r ListLimitRule) Error(message string) ListLimitRule {
	r.err = r.err.SetMessage(message)
	return r
}

// ErrorObject sets the error struct for the rule.
func (r ListLimitRule) ErrorObject(err validation.Error) ListLimitRule {
	r.err = err
	return r
}
