// Copyright 2021-2026 Altessa Solutions Inc. All rights reserved.
// Use of this source code is governed by license that can be found in
// the LICENSE file.

package ozzo_rules

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/asaskevich/govalidator"
	"github.com/go-ozzo/ozzo-validation/v4"
)

// ErrServerAddressInvalid is the error that returns when a value is not a valid server address.
var ErrServerAddressInvalid = validation.NewError("validation_server_address", "invalid server address")

// ServerAddress is a validation rule that checks if a value is a valid server address.
func ServerAddress() ServerAddressRule {
	return ServerAddressRule{condition: true, err: ErrServerAddressInvalid}
}

// ServerAddressRule is a rule that checks if a value is a valid server address.
type ServerAddressRule struct {
	err       validation.Error
	condition bool
}

// Validate checks if the given value is valid or not.
func (r ServerAddressRule) Validate(v any) error {
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

	re := regexp.MustCompile("^[a-zA-Z]+://")
	splits := re.Split(str, -1)
	if len(splits) > 1 {
		str = splits[1]
	}

	comps := strings.Split(str, ":")
	if len(comps) != 2 { //nolint:mnd
		return r.err
	}

	p, e := strconv.Atoi(comps[1])
	if e != nil {
		return r.err
	}

	if p <= 0 || p >= 65535 { //nolint:mnd
		return r.err
	}

	if !govalidator.IsHost(comps[0]) {
		return r.err
	}

	return nil
}

// When sets the condition that determines if the validation should be performed.
func (r ServerAddressRule) When(condition bool) ServerAddressRule {
	r.condition = condition
	return r
}

// Error sets the error message for the rule.
func (r ServerAddressRule) Error(message string) ServerAddressRule {
	r.err = r.err.SetMessage(message)
	return r
}

// ErrorObject sets the error struct for the rule.
func (r ServerAddressRule) ErrorObject(err validation.Error) ServerAddressRule {
	r.err = err
	return r
}
