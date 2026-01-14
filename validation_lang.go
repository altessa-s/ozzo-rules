// Copyright 2021-2026 Altessa Solutions Inc. All rights reserved.
// Use of this source code is governed by license that can be found in
// the LICENSE file.

package ozzo_rules

import (
	"github.com/asaskevich/govalidator"
	"github.com/go-ozzo/ozzo-validation/v4"
)

// ErrLangCode2Invalid is the error that returns when a value is not a valid two-letter language code.
var ErrLangCode2Invalid = validation.NewError("validation_lang_code2", "must be a valid two-letter language code")

// LangCode2 is a validation rule that checks if a value is a valid ISO 639-1 language code.
func LangCode2() LangCode2Rule {
	return LangCode2Rule{condition: true, err: ErrLangCode2Invalid}
}

// LangCode2Rule is a rule that checks if a value is a valid ISO 639-1 language code.
type LangCode2Rule struct {
	err       validation.Error
	condition bool
}

// Validate checks if the given value is valid or not.
func (r LangCode2Rule) Validate(v any) error {
	if !r.condition {
		return nil
	}

	value, isNil := validation.Indirect(v)
	if isNil {
		return r.err
	}

	str, err := validation.EnsureString(value)
	if err != nil {
		return r.err
	}

	if str == "" || !govalidator.IsISO693Alpha2(str) {
		return r.err
	}

	return nil
}

// When sets the condition that determines if the validation should be performed.
func (r LangCode2Rule) When(condition bool) LangCode2Rule {
	r.condition = condition
	return r
}

// Error sets the error message for the rule.
func (r LangCode2Rule) Error(message string) LangCode2Rule {
	r.err = r.err.SetMessage(message)
	return r
}

// ErrorObject sets the error struct for the rule.
func (r LangCode2Rule) ErrorObject(err validation.Error) LangCode2Rule {
	r.err = err
	return r
}
