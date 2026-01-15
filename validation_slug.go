// Copyright 2021-2026 Altessa Solutions Inc. All rights reserved.
// Use of this source code is governed by license that can be found in
// the LICENSE file.

package ozzo_rules

import (
	"regexp"

	"github.com/go-ozzo/ozzo-validation/v4"
)

// ErrSlugInvalid is the error that returns when a value is not a valid slug.
var ErrSlugInvalid = validation.NewError("validation_slug", "must be between 5 and 100 "+
	"characters long and contain only lowercase letters, numbers, and underscores")

// Slug is a validation rule that checks if a value is a valid slug.
func Slug() SlugRule {
	return SlugRule{condition: true, err: ErrSlugInvalid}
}

// SlugRule is a rule that checks if a value is a valid slug.
type SlugRule struct {
	err       validation.Error
	condition bool
}

// Validate checks if the given value is valid or not.
func (r SlugRule) Validate(v any) error {
	if !r.condition {
		return nil
	}

	value, isNil := validation.Indirect(v)
	if isNil {
		return r.err
	}

	str, err := validation.EnsureString(value)
	if err != nil || !slugRx.MatchString(str) {
		return r.err
	}
	return nil
}

// When sets the condition that determines if the validation should be performed.
func (r SlugRule) When(condition bool) SlugRule {
	r.condition = condition
	return r
}

// Error sets the error message for the rule.
func (r SlugRule) Error(message string) SlugRule {
	r.err = r.err.SetMessage(message)
	return r
}

// ErrorObject sets the error struct for the rule.
func (r SlugRule) ErrorObject(err validation.Error) SlugRule {
	r.err = err
	return r
}

var slugRx = regexp.MustCompile("^[a-z0-9_]{5,100}$")
