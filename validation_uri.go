// Copyright 2021-2026 Altessa Solutions Inc. All rights reserved.
// Use of this source code is governed by license that can be found in
// the LICENSE file.

package ozzo_rules

import (
	"strconv"
	"strings"

	"github.com/asaskevich/govalidator"
	"github.com/go-ozzo/ozzo-validation/v4"
)

// ErrURIInvalid is the error that returns when a value is not a valid URI address.
var ErrURIInvalid = validation.NewError("validation_uri_address", "invalid uri address")

// URI is a validation rule that checks if a value is a valid URI address.
func URI() URIRule {
	return URIRule{err: ErrURIInvalid, condition: true}
}

// URIRule is a rule that checks if a value is a valid URI address.
type URIRule struct {
	err       validation.Error
	condition bool
}

// Validate checks if the given value is valid or not.
func (r URIRule) Validate(v any) error {
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

	splits := strings.Split(str, "://")

	if splits[0] != "starttls" && splits[0] != "ssl" {
		return r.err
	}

	if len(splits) > 1 {
		str = splits[1]
	}

	addrComps := strings.Split(str, ":")

	if len(addrComps) != 2 { //nolint:mnd
		return r.err
	}

	p, e := strconv.Atoi(addrComps[1])
	if e != nil {
		return r.err
	}

	if p <= 0 || p >= 65535 { //nolint:mnd
		return r.err
	}

	if !govalidator.IsHost(addrComps[0]) {
		return r.err
	}

	return nil
}

// When sets the condition that determines if the validation should be performed.
func (r URIRule) When(condition bool) URIRule {
	r.condition = condition
	return r
}

// Error sets the error message for the rule.
func (r URIRule) Error(message string) URIRule {
	r.err = r.err.SetMessage(message)
	return r
}

// ErrorObject sets the error struct for the rule.
func (r URIRule) ErrorObject(err validation.Error) URIRule {
	r.err = err
	return r
}
