// Copyright 2021-2026 Altessa Solutions Inc. All rights reserved.
// Use of this source code is governed by license that can be found in
// the LICENSE file.

package ozzo_rules

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/go-ozzo/ozzo-validation/v4"
)

const mbrLength = 8

// ErrRsMBRInvalid is the error that returns when a value is not a valid Serbian Company Registration Number(MBR).
var ErrRsMBRInvalid = validation.NewError("validation_rs_mbr_invalid", "invalid mbr number")

// RsMBR is a validation rule that checks if a value is a valid Serbian Company Registration Number(MBR).
func RsMBR() RsMBRRule {
	return RsMBRRule{condition: true, err: ErrRsMBRInvalid}
}

// RsMBRRule is a rule that checks if a value is a valid Serbian Company Registration Number(MBR).
// MBR is a 8-digit number that is used in Serbia to identify legal entities.
type RsMBRRule struct {
	err       validation.Error
	condition bool
}

// Validate checks if the given value is valid or not.
func (r RsMBRRule) Validate(data any) error {
	if !r.condition {
		return nil
	}

	data, isNil := validation.Indirect(data)
	if isNil {
		return r.err
	}

	var num string

	switch t := data.(type) {
	case int64:
		num = strconv.FormatInt(t, 10)
		// Sometimes the number is 7 digits long and it is valid number, so we need to add a leading zero
		if len(num) == 7 { // nolint: mnd
			num = fmt.Sprintf("0%s", num)
		}
	case string:
		num = t
	default:
		return r.err
	}

	if len(num) != mbrLength {
		return r.err
	}

	ok, err := mod11_7(num)
	if !ok || err != nil {
		return r.err
	}

	return nil
}

// When sets the condition that determines if the validation should be performed.
func (r RsMBRRule) When(condition bool) RsMBRRule {
	r.condition = condition
	return r
}

// Error sets the error message for the rule.
func (r RsMBRRule) Error(message string) RsMBRRule {
	r.err = r.err.SetMessage(message)
	return r
}

// ErrorObject sets the error struct for the rule.
func (r RsMBRRule) ErrorObject(err validation.Error) RsMBRRule {
	r.err = err
	return r
}

var errInvalidNumber = errors.New("invalid number")

// The mod11_7 function calculates the MOD 11-7 check digit for an 8-digit number.
// sum = 2a1 + 3a2 + 4a3 + 5a4 + 6a5 + 7a6 + 2a7 mod 11
// nolint: mnd
func mod11_7(in string) (bool, error) {
	product := int64(2)
	var sum int64

	for i := len(in) - 2; i >= 0; i-- {
		n, err := strconv.ParseInt(string(in[i]), 10, 64)
		if err != nil {
			return false, errInvalidNumber
		}

		sum += n * product
		if product == 7 {
			product = 2
		} else {
			product++
		}
	}

	control := 11 - (sum % 11)

	if control == 10 || control == 11 {
		control = 0
	}

	lastDigit, err := strconv.ParseInt(string(in[len(in)-1]), 10, 64)
	if err != nil {
		return false, errInvalidNumber
	}

	return control == lastDigit, nil
}
