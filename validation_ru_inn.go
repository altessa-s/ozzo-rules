// Copyright 2021-2026 Altessa Solutions Inc. All rights reserved.
// Use of this source code is governed by license that can be found in
// the LICENSE file.

package ozzo_rules

import (
	"regexp"
	"strconv"

	"github.com/go-ozzo/ozzo-validation/v4"
)

// ErrRuINNInvalid is the error that returns when a value is not a valid INN number.
var ErrRuINNInvalid = validation.NewError("validation_ru_inn_invalid", "invalid INN number, INN must be 10 or "+
	"12 characters long and contain only digits")

var ErrRuINNPersonalInvalid = validation.NewError("validation_ru_inn_personal_invalid", "invalid personal INN number, "+
	"personal INN must be 12 characters long and contain only digits")

var ErrRuINNLegalInvalid = validation.NewError("validation_ru_inn_legal_invalid", "invalid legal INN number, "+
	"legal INN must be 10 characters long and contain only digits")

// RuINN is a validation rule that checks if a value is a valid INN number.
func RuINN() RuINNRule {
	return RuINNRule{condition: true, err: ErrRuINNInvalid}
}

// RuINNPersonal is a validation rule that checks if a value is a valid personal INN number.
func RuINNPersonal() RuINNRule {
	return RuINNRule{condition: true, innType: ruInnTypePersonal, err: ErrRuINNPersonalInvalid}
}

// RuINNLegal is a validation rule that checks if a value is a valid legal INN number.
func RuINNLegal() RuINNRule {
	return RuINNRule{condition: true, innType: ruInnTypeLegal, err: ErrRuINNLegalInvalid}
}

type ruInnType int

const (
	ruInnTypeAny      ruInnType = 0
	ruInnTypePersonal ruInnType = 1
	ruInnTypeLegal    ruInnType = 2
)

// RuINNRule is a rule that checks if a value is a valid INN number.
type RuINNRule struct {
	err       validation.Error
	condition bool
	innType   ruInnType
}

// Validate checks if the given value is valid or not.
func (r RuINNRule) Validate(v any) error {
	if !r.condition {
		return nil
	}

	value, isNil := validation.Indirect(v)
	if isNil {
		return r.err
	}

	str := ""

	switch t := value.(type) {
	case int64:
		str = strconv.FormatInt(t, 10)
	case string:
		str = t
	default:
		return r.err
	}

	var rx *regexp.Regexp
	switch r.innType {
	case ruInnTypePersonal:
		rx = ruInnPersonalRx
	case ruInnTypeLegal:
		rx = ruInnLegalRx
	default:
		rx = ruInnRx
	}

	if ok := rx.MatchString(str); !ok {
		return r.err
	}

	if r.innType == ruInnTypeAny {
		r.innType = ruInnTypePersonal
		if len(str) == 10 { //nolint:mnd
			r.innType = ruInnTypeLegal
		}
	}

	checkSum := func(factors []int) int {
		sum := 0
		for i, f := range factors {
			sum += f * int(str[i]-'0')
		}
		return sum % 11 % 10 //nolint:mnd
	}

	switch r.innType {
	case ruInnTypeLegal:
		n10 := checkSum([]int{2, 4, 10, 3, 5, 9, 4, 6, 8})
		if n10 == int(str[9]-'0') {
			return nil
		}
	case ruInnTypePersonal:
		n11 := checkSum([]int{7, 2, 4, 10, 3, 5, 9, 4, 6, 8})
		n12 := checkSum([]int{3, 7, 2, 4, 10, 3, 5, 9, 4, 6, 8})
		if n11 == int(str[10]-'0') && n12 == int(str[11]-'0') {
			return nil
		}
	default:
		return r.err
	}

	return r.err
}

// When sets the condition that determines if the validation should be performed.
func (r RuINNRule) When(condition bool) RuINNRule {
	r.condition = condition
	return r
}

// Error sets the error message for the rule.
func (r RuINNRule) Error(message string) RuINNRule {
	r.err = r.err.SetMessage(message)
	return r
}

// ErrorObject sets the error struct for the rule.
func (r RuINNRule) ErrorObject(err validation.Error) RuINNRule {
	r.err = err
	return r
}

var ruInnRx = regexp.MustCompile(`^([\d]{10}|[\d]{12})$`)
var ruInnPersonalRx = regexp.MustCompile(`^[\d]{12}$`)
var ruInnLegalRx = regexp.MustCompile(`^[\d]{10}$`)
