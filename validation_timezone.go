// Copyright 2021-2026 Altessa Solutions Inc. All rights reserved.
// Use of this source code is governed by license that can be found in
// the LICENSE file.

package ozzo_rules

import (
	"time"

	"github.com/go-ozzo/ozzo-validation/v4"
)

// ErrTimezoneInvalid is the error that returns when a value is not a valid timezone.
var ErrTimezoneInvalid = validation.NewError("validation_timezone", "must be valid timezone")

// TimezoneOrNil creates a validation rule that allows nil values.
func TimezoneOrNil() TimezoneRule {
	return TimezoneRule{err: ErrTimezoneInvalid, condition: true, allowNil: true}
}

// Timezone creates a validation rule that checks if a value is a valid timezone.
func Timezone() TimezoneRule {
	return TimezoneRule{err: ErrTimezoneInvalid, condition: true}
}

// TimezoneRule is a rule that checks if a value is a valid timezone.
type TimezoneRule struct {
	err       validation.Error
	condition bool
	allowNil  bool
}

// Validate checks if the given value is valid or not.
func (r TimezoneRule) Validate(v any) error {
	if !r.condition {
		return nil
	}

	value, isNil := validation.Indirect(v)
	if isNil {
		if r.allowNil {
			return nil
		}
		return r.err
	}

	tzone, err := validation.EnsureString(value)
	if err != nil {
		return r.err
	}

	if tzone == "" {
		return r.err
	}

	_, err = time.LoadLocation(tzone)
	if err != nil {
		if isKnownTimezoneAbbreviation(tzone) {
			return nil
		}
		return r.err
	}

	return nil
}

func (r TimezoneRule) When(condition bool) TimezoneRule {
	r.condition = condition
	return r
}

func (r TimezoneRule) Error(message string) TimezoneRule {
	r.err = r.err.SetMessage(message)
	return r
}

func (r TimezoneRule) ErrorObject(err validation.Error) TimezoneRule {
	r.err = err
	return r
}

// isKnownTimezoneAbbreviation checks if the given string is a known timezone abbreviation
func isKnownTimezoneAbbreviation(tz string) bool {
	// Common timezone abbreviations that are valid but not recognized by time.LoadLocation
	knownAbbreviations := map[string]bool{
		"PST": true, // Pacific Standard Time
		"CST": true, // Central Standard Time
		// "EST" is handled by time.LoadLocation, so not included here
		"MST":  true, // Mountain Standard Time
		"PDT":  true, // Pacific Daylight Time
		"CDT":  true, // Central Daylight Time
		"EDT":  true, // Eastern Daylight Time
		"MDT":  true, // Mountain Daylight Time
		"GMT":  true, // Greenwich Mean Time
		"BST":  true, // British Summer Time
		"IST":  true, // Indian Standard Time
		"JST":  true, // Japan Standard Time
		"KST":  true, // Korea Standard Time
		"AEST": true, // Australian Eastern Standard Time
		"AEDT": true, // Australian Eastern Daylight Time
		"AWST": true, // Australian Western Standard Time
		"ACST": true, // Australian Central Standard Time
		"ACDT": true, // Australian Central Daylight Time
		"NZST": true, // New Zealand Standard Time
		"NZDT": true, // New Zealand Daylight Time
	}

	return knownAbbreviations[tz]
}
