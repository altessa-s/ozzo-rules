// Copyright 2021-2026 Altessa Solutions Inc. All rights reserved.
// Use of this source code is governed by license that can be found in
// the LICENSE file.

package ozzo_rules

import (
	"github.com/go-ozzo/ozzo-validation/v4"
)

// ErrRegexInvalid is the error that returns if the regular expression is invalid.
var ErrRegexInvalid = validation.NewError("validation_regex_invalid", "invalid regular expression")

// Regex returns a new rule that checks if a value is a valid regular expression.
func Regex() RegexRule {
	return RegexRule{condition: true, err: ErrRegexInvalid}
}

// RegexRule is a rule that checks if a value is a valid regular expression.
type RegexRule struct {
	err       validation.Error
	condition bool
}

// Validate checks if the given value is valid or not.
func (r RegexRule) Validate(data any) error {
	if !r.condition {
		return nil
	}

	re, isNil := validation.Indirect(data)
	if isNil {
		return r.err
	}

	if str, ok := re.(string); ok {
		if _, err := compileRegex(str); err != nil {
			return r.err
		}
		return nil
	}

	return r.err
}

// When sets the condition that determines if the validation should be performed.
func (r RegexRule) When(condition bool) RegexRule {
	r.condition = condition
	return r
}

// Error sets the error message for the rule.
func (r RegexRule) Error(message string) RegexRule {
	r.err = r.err.SetMessage(message)
	return r
}

// ErrorObject sets the error struct for the rule.
func (r RegexRule) ErrorObject(err validation.Error) RegexRule {
	r.err = err
	return r
}
