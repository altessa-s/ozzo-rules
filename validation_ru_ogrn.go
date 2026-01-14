// Copyright 2021-2026 Altessa Solutions Inc. All rights reserved.
// Use of this source code is governed by license that can be found in
// the LICENSE file.

package ozzo_rules

import (
	"regexp"
	"strconv"

	"github.com/go-ozzo/ozzo-validation/v4"
)

// ErrRuOGRNInvalid is the error that returns when a value is not a valid OGRN/OGRNIP number.
var ErrRuOGRNInvalid = validation.NewError("validation_ru_ogrn_invalid", "invalid OGRN/OGRNIP number")

// ErrRuOGRNEntrepreneurInvalid is the error that returns when a value is not a valid OGRNIP number.
var ErrRuOGRNEntrepreneurInvalid = validation.NewError("validation_ru_ogrn_invalid", "invalid OGRNIP number")

// ErrRuOGRNLegalInvalid is the error that returns when a value is not a valid OGRN number.
var ErrRuOGRNLegalInvalid = validation.NewError("validation_ru_ogrn_invalid", "invalid OGRN number")

// RuOGRN is a validation rule that checks if a value is a valid OGRN/OGRNIP number.
func RuOGRN() RuOGRNRule {
	return RuOGRNRule{condition: true, err: ErrRuOGRNInvalid}
}

// RuOGRNEntrepreneur is a validation rule that checks if a value is a valid OGRNIP number.
func RuOGRNEntrepreneur() RuOGRNRule {
	return RuOGRNRule{condition: true, ogrnType: ruOgrnTypeEntrepreneur, err: ErrRuOGRNEntrepreneurInvalid}
}

// RuOGRNLegal is a validation rule that checks if a value is a valid OGRN number.
func RuOGRNLegal() RuOGRNRule {
	return RuOGRNRule{condition: true, ogrnType: ruOgrnTypeLegal, err: ErrRuOGRNLegalInvalid}
}

type ruOgrnType int

const (
	ruOgrnTypeAny          ruOgrnType = 0
	ruOgrnTypeEntrepreneur ruOgrnType = 1
	ruOgrnTypeLegal        ruOgrnType = 2
)

// RuOGRNRule is a rule that checks if a value is a valid OGRN number.
type RuOGRNRule struct {
	err       validation.Error
	condition bool
	ogrnType  ruOgrnType
}

// Validate checks if the given value is valid or not.
func (r RuOGRNRule) Validate(data any) error {
	if !r.condition {
		return nil
	}

	data, isNil := validation.Indirect(data)
	if isNil {
		return r.err
	}

	var ogrn string
	switch t := data.(type) {
	case int64:
		ogrn = strconv.FormatInt(t, 10)
	case string:
		ogrn = t
	default:
		return r.err
	}

	var rx *regexp.Regexp
	switch r.ogrnType {
	case ruOgrnTypeEntrepreneur:
		rx = ruOgrnEntrepreneurRx
	case ruOgrnTypeLegal:
		rx = ruOgrnLegalRx
	default:
		rx = ruOgrnRx
	}

	if ok := rx.MatchString(ogrn); !ok {
		return r.err
	}

	length := len(ogrn)
	var divisor int64
	if length == 13 { //nolint:mnd
		divisor = 11
	} else {
		divisor = 13 // 15 chars case
	}

	// Calculate modulus using streaming method to avoid potential overflow on 32-bit systems
	// although 15 digits fits in int64, using a custom loop is safer/cleaner than Atoi for large numbers.
	// We only need the digits up to length-1.
	var remainder int64
	for i := 0; i < length-1; i++ {
		digit := int64(ogrn[i] - '0')
		remainder = (remainder*10 + digit) % divisor
	}

	// Legal (13 chars):  (N1...N12) % 11. If result == 10, check digit is 0.
	// Entrepreneur (15 chars): (N1...N14) % 13. Check digit is result % 10.

	if length == 15 { //nolint:mnd
		remainder %= 10
	} else if remainder == 10 { //nolint:mnd
		remainder = 0
	}

	// Last digit is the check digit
	lastDigit := int64(ogrn[length-1] - '0')

	if remainder == lastDigit {
		return nil
	}

	return r.err
}

// When sets the condition that determines if the validation should be performed.
func (r RuOGRNRule) When(condition bool) RuOGRNRule {
	r.condition = condition
	return r
}

// Error sets the error message for the rule.
func (r RuOGRNRule) Error(message string) RuOGRNRule {
	r.err = r.err.SetMessage(message)
	return r
}

// ErrorObject sets the error struct for the rule.
func (r RuOGRNRule) ErrorObject(err validation.Error) RuOGRNRule {
	r.err = err
	return r
}

var ruOgrnRx = regexp.MustCompile(`^([\d]{13}|[\d]{15})$`)
var ruOgrnLegalRx = regexp.MustCompile(`^\d{13}$`)
var ruOgrnEntrepreneurRx = regexp.MustCompile(`^\d{15}$`)
