// Copyright 2021-2026 Altessa Solutions Inc. All rights reserved.
// Use of this source code is governed by license that can be found in
// the LICENSE file.

package ozzo_rules

import (
	"errors"
	"time"

	"github.com/go-ozzo/ozzo-validation/v4"
)

var (
	ErrDurationOutOfRange = validation.NewError("validation_duration_out_of_range",
		"the duration must be between {{.min}} and {{.max}}")
)

type DurationLimit interface {
	string | time.Duration | int64
}

// DurationWithLimit is a validation rule that checks if a value is not empty and a valid duration.
func DurationWithLimit[T DurationLimit](minDur, maxDur T) DurationLimitRule {
	return DurationLimitRule{
		condition: true,
		err:       buildRuleError(minDur, maxDur),
		min:       minDur,
		max:       maxDur,
	}
}

// DurationLimitRule is a rule that checks if a value is a valid duration.
type DurationLimitRule struct {
	err       validation.Error
	condition bool
	min       any
	max       any
}

// Validate checks if the given value is valid or not.
func (r DurationLimitRule) Validate(v any) error {
	if !r.condition {
		return nil
	}

	value, isNil := validation.Indirect(v)
	if isNil {
		return r.err
	}

	minDuration, _ := parseDuration(r.min) //nolint:errcheck
	maxDuration, _ := parseDuration(r.max) //nolint:errcheck

	dur, err := parseDuration(value)
	if err != nil {
		return r.err
	}

	if minDuration >= 0 && dur < minDuration {
		return r.err
	} else if maxDuration >= 0 && dur > maxDuration {
		return r.err
	}

	return nil
}

// When sets the condition that determines if the validation should be performed.
func (r DurationLimitRule) When(condition bool) DurationLimitRule {
	r.condition = condition
	return r
}

// Error sets the error message for the rule.
func (r DurationLimitRule) Error(message string) DurationLimitRule {
	r.err = r.err.SetMessage(message)
	return r
}

// ErrorObject sets the error struct for the rule.
func (r DurationLimitRule) ErrorObject(err validation.Error) DurationLimitRule {
	r.err = err
	return r
}

func parseDuration(duration any) (time.Duration, error) {
	switch t := duration.(type) {
	case string:
		return parseDurationString(t)
	case time.Duration:
		return t, nil
	case int64:
		return time.Duration(t) * time.Second, nil
	}

	return -1, errors.New("invalid type") //nolint:goerr113
}

func buildRuleError(m, mx any) validation.Error {
	minDur, err := parseDuration(m)
	if err != nil {
		panic("invalid min value: " + err.Error())
	}

	maxDur, err := parseDuration(mx)
	if err != nil {
		panic("invalid min value: " + err.Error())
	}

	return ErrDurationOutOfRange.
		SetParams(map[string]any{"min": minDur.String(), "max": maxDur.String()})
}
