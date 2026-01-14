// Copyright 2021-2026 Altessa Solutions Inc. All rights reserved.
// Use of this source code is governed by license that can be found in
// the LICENSE file.

package ozzo_rules

import (
	"strconv"

	"github.com/go-ozzo/ozzo-validation/v4"
)

// ErrRsPIBInvalid is the error that returns when a value is not a valid PIB number.
var ErrRsPIBInvalid = validation.NewError("validation_rs_pib_invalid", "invalid PIB number")

// RsPIB is a validation rule that checks if a value is a valid PIB number.
func RsPIB() RsPIBRule {
	return RsPIBRule{condition: true, err: ErrRsPIBInvalid}
}

// RsPIBRule is a rule that checks if a value is a valid PIB number.
// PIB number is a 9-digit number that is used in Serbia to identify legal entities.
type RsPIBRule struct {
	err       validation.Error
	condition bool
}

// Validate checks if the given value is valid or not.
func (r RsPIBRule) Validate(data any) error {
	if !r.condition {
		return nil
	}

	data, isNil := validation.Indirect(data)
	if isNil {
		return r.err
	}

	var ok bool
	switch t := data.(type) {
	case string:
		ok = isValidMod1110(t)
	case int64:
		ok = isValidMod1110(t)
	default:
		return r.err
	}

	if !ok {
		return r.err
	}

	return nil
}

// When sets the condition that determines if the validation should be performed.
func (r RsPIBRule) When(condition bool) RsPIBRule {
	r.condition = condition
	return r
}

// Error sets the error message for the rule.
func (r RsPIBRule) Error(message string) RsPIBRule {
	r.err = r.err.SetMessage(message)
	return r
}

// ErrorObject sets the error struct for the rule.
func (r RsPIBRule) ErrorObject(err validation.Error) RsPIBRule {
	r.err = err
	return r
}

// nolint: mnd
func mod1110[T interface{ string | int64 }](data T) (bool, error) {
	var num string

	switch t := any(data).(type) {
	case int64:
		num = strconv.FormatInt(t, 10)
	default:
		var ok bool
		num, ok = t.(string)
		if !ok {
			return false, errInvalidNumber
		}
	}

	if len(num) != 9 {
		return false, errInvalidNumber
	}

	product := int64(10)
	sum := int64(0)

	numRunes := []rune(num)

	for i := range 8 {
		n, err := strconv.ParseInt(string(numRunes[i]), 10, 64)
		if err != nil {
			return false, errInvalidNumber
		}

		if sum = (n + product) % 10; sum == 0 {
			sum = 10
		}

		product = (2 * sum) % 11
	}

	lastDigit, err := strconv.ParseInt(string(numRunes[8]), 10, 64)
	if err != nil {
		return false, errInvalidNumber
	}

	checkDigit := (product + lastDigit) % 10

	return checkDigit == int64(1), nil
}

// IsValidMod11_10 checks if the ISO 7064, MOD 11-10 check digit is valid.
func isValidMod1110[T interface{ string | int64 }](data T) bool {
	ok, err := mod1110(data)
	if err != nil {
		return false
	}
	return ok
}
