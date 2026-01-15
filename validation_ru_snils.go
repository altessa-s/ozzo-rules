// Copyright 2021-2026 Altessa Solutions Inc. All rights reserved.
// Use of this source code is governed by license that can be found in
// the LICENSE file.

package ozzo_rules

import (
	"regexp"
	"strconv"

	"github.com/go-ozzo/ozzo-validation/v4"
)

// ErrRuSNILSInvalid is the error that returns when a value is not a valid SNILS number.
var ErrRuSNILSInvalid = validation.NewError("validation_ru_snils_invalid", "invalid SNILS number")

// RuSNILS is a validation rule that checks if a value is a valid SNILS number.
func RuSNILS() RuSNILSRule {
	return RuSNILSRule{condition: true, err: ErrRuSNILSInvalid}
}

// RuSNILSRule is a rule that checks if a value is a valid SNILS number.
type RuSNILSRule struct {
	err       validation.Error
	condition bool
}

// Validate checks if the given value is valid or not.
func (r RuSNILSRule) Validate(data any) error {
	data, isNil := validation.Indirect(data)
	if isNil {
		return r.err
	}

	var snils string
	switch t := data.(type) {
	case int64:
		snils = strconv.FormatInt(t, 10)
	case string:
		snils = t
	default:
		return r.err
	}

	if !ruSnilsRx.MatchString(snils) {
		return r.err
	}

	// Manual parsing to extract digits for checksum, avoiding allocs from FindStringSubmatch
	var digits [11]int
	var idx int

	for _, c := range snils {
		if c >= '0' && c <= '9' {
			// We know from regex match that there are exactly 11 digits
			digits[idx] = int(c - '0')
			idx++
		}
	}

	// We rely on regex for length check

	sum := 0
	for i := 0; i < 9; i++ {
		sum += digits[i] * (9 - i) //nolint:mnd
	}

	sum %= 101
	if sum == 100 { //nolint:mnd
		sum = 0
	}

	// Checksum is the last two digits (index 9 and 10) treated as a number
	controlSum := digits[9]*10 + digits[10] //nolint:mnd

	if sum != controlSum {
		return r.err
	}

	return nil
}

// When sets the condition that determines if the validation should be performed.
func (r RuSNILSRule) When(condition bool) RuSNILSRule {
	r.condition = condition
	return r
}

// Error sets the error message for the rule.
func (r RuSNILSRule) Error(message string) RuSNILSRule {
	r.err = r.err.SetMessage(message)
	return r
}

// ErrorObject sets the error struct for the rule.
func (r RuSNILSRule) ErrorObject(err validation.Error) RuSNILSRule {
	r.err = err
	return r
}

var ruSnilsRx = regexp.MustCompile(`^(\d{3})[\s-]?(\d{3})[\s-]?(\d{3})[\s-]?(\d{2})$`)
