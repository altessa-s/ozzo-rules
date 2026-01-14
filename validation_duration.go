// Copyright 2021-2026 Altessa Solutions Inc. All rights reserved.
// Use of this source code is governed by license that can be found in
// the LICENSE file.

package ozzo_rules

import (
	"strconv"
	"strings"
	"time"

	"github.com/go-ozzo/ozzo-validation/v4"
)

// ErrDurationInvalid is the error that returns when a value is not a valid duration.
var ErrDurationInvalid = validation.NewError("validation_duration", "invalid duration")

// Duration is a validation rule that checks if a value is not empty and a valid duration.
func Duration() DurationRule {
	return DurationRule{condition: true, err: ErrDurationInvalid}
}

// DurationOrZero is a validation rule that checks if a value is valid duration or empty.
func DurationOrZero() DurationRule {
	return DurationRule{condition: true, allowedZero: true, err: ErrDurationInvalid}
}

// DurationRule is a rule that checks if a value is a valid duration.
type DurationRule struct {
	err         validation.Error
	allowedZero bool
	condition   bool
}

// Validate checks if the given value is valid or not.
func (r DurationRule) Validate(v any) error {
	if !r.condition {
		return nil
	}

	value, isNil := validation.Indirect(v)
	if isNil {
		return r.err
	}

	var dur time.Duration = 0

	switch valType := value.(type) {
	case string:
		dString := valType
		if dString == "" {
			return r.err
		}
		var err error
		if dur, err = parseDurationString(dString); err != nil {
			return r.err
		}
	case time.Duration:
		dur = valType
	}

	if dur == 0 && !r.allowedZero {
		return r.err.SetMessage("must be greater 0")
	}

	return nil
}

// When sets the condition that determines if the validation should be performed.
func (r DurationRule) When(condition bool) DurationRule {
	r.condition = condition
	return r
}

// Error sets the error message for the rule.
func (r DurationRule) Error(message string) DurationRule {
	r.err = r.err.SetMessage(message)
	return r
}

// ErrorObject sets the error struct for the rule.
func (r DurationRule) ErrorObject(err validation.Error) DurationRule {
	r.err = err
	return r
}

func parseDurationString(durationString string) (time.Duration, error) {
	if durationString == "" {
		return 0, nil
	}

	if strings.HasSuffix(durationString, "s") || strings.HasSuffix(durationString, "m") ||
		strings.HasSuffix(durationString, "h") || strings.HasSuffix(durationString, "ms") {
		return time.ParseDuration(durationString)
	}

	secs, err := strconv.ParseInt(durationString, 10, 64)
	if err != nil {
		return -1, err
	}
	return time.Duration(secs) * time.Second, nil
}
