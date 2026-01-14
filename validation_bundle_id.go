// Copyright 2021-2026 Altessa Solutions Inc. All rights reserved.
// Use of this source code is governed by license that can be found in
// the LICENSE file.

package ozzo_rules

import (
	"regexp"

	"github.com/go-ozzo/ozzo-validation/v4"
)

// ErrBundleIDInvalid is the error that returns when a value is not a valid bundle id.
var ErrBundleIDInvalid = validation.NewError("validation_bundle_id",
	"can only contain alphanumeric characters, hyphens, and periods. It can't end on a period.")

// BundleID is a validation rule that checks if a value is a valid bundle id.
func BundleID() BundleIDRule {
	return BundleIDRule{condition: true, err: ErrBundleIDInvalid}
}

// BundleIDRule is a rule that checks if a value is a valid bundle id.
type BundleIDRule struct {
	err       validation.Error
	condition bool
}

// Validate checks if the given value is valid or not.
func (r BundleIDRule) Validate(v any) error {
	if !r.condition {
		return nil
	}

	value, isNil := validation.Indirect(v)
	if isNil {
		return r.err
	}

	str, err := validation.EnsureString(value)
	if err != nil || !bundleIDRx.MatchString(str) {
		return r.err
	}
	return nil
}

// When sets the condition that determines if the validation should be performed.
func (r BundleIDRule) When(condition bool) BundleIDRule {
	r.condition = condition
	return r
}

// Error sets the error message for the rule.
func (r BundleIDRule) Error(message string) BundleIDRule {
	r.err = r.err.SetMessage(message)
	return r
}

// ErrorObject sets the error struct for the rule.
func (r BundleIDRule) ErrorObject(err validation.Error) BundleIDRule {
	r.err = err
	return r
}

var bundleIDRx = regexp.MustCompile(`^[a-z0-9]+(\.[a-z0-9]+)+$`)
