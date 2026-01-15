// Copyright 2021-2026 Altessa Solutions Inc. All rights reserved.
// Use of this source code is governed by license that can be found in
// the LICENSE file.

package ozzo_rules

import (
	"regexp"

	"github.com/go-ozzo/ozzo-validation/v4"
)

// ErrApnTeamIdInvalid is the error that returns when a value is not a valid apn team.
var ErrApnTeamIdInvalid = validation.NewError("validation_apn_team_id",
	"must be 10 characters long and may only contain alphanumeric characters")

// ApnTeamId is a validation rule that checks if a value is a valid apn team.
func ApnTeamId() ApnTeamIdRule {
	return ApnTeamIdRule{condition: true, err: ErrApnTeamIdInvalid}
}

// ApnTeamIdRule is a rule that checks if a value is a valid apn team.
type ApnTeamIdRule struct {
	err       validation.Error
	condition bool
}

// Validate checks if the given value is valid or not.
func (r ApnTeamIdRule) Validate(v any) error {
	if !r.condition {
		return nil
	}

	value, isNil := validation.Indirect(v)
	if isNil {
		return r.err
	}

	str, err := validation.EnsureString(value)
	if err != nil || !teamRx.MatchString(str) {
		return r.err
	}
	return nil
}

// When sets the condition that determines if the validation should be performed.
func (r ApnTeamIdRule) When(condition bool) ApnTeamIdRule {
	r.condition = condition
	return r
}

// Error sets the error message for the rule.
func (r ApnTeamIdRule) Error(message string) ApnTeamIdRule {
	r.err = r.err.SetMessage(message)
	return r
}

// ErrorObject sets the error struct for the rule.
func (r ApnTeamIdRule) ErrorObject(err validation.Error) ApnTeamIdRule {
	r.err = err
	return r
}

var teamRx = regexp.MustCompile(`^[A-Z0-9.\-]{10}$`)
