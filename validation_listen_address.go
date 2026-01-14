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

// ErrListenAddressInvalid is the error that returns when a value is not a valid listen address.
var ErrListenAddressInvalid = validation.NewError("validation_listen_address", "invalid listen address")

// ListenAddress is a validation rule that checks if a value is a valid listen address.
func ListenAddress() ListenAddressRule {
	return ListenAddressRule{condition: true, err: ErrListenAddressInvalid}
}

type ListenAddressRule struct {
	err       validation.Error
	condition bool
}

// Validate checks if the given value is valid or not.
func (r ListenAddressRule) Validate(v any) error {
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

	splits := addrRx.Split(str, -1)
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

	host := addrComps[0]

	if host == "" || strings.HasPrefix(host, "_") || strings.HasPrefix(host, ".") {
		return r.err
	}

	if !govalidator.IsHost(host) {
		return r.err
	}

	return nil
}

// When sets the condition that determines if the validation should be performed.
func (r ListenAddressRule) When(condition bool) ListenAddressRule {
	r.condition = condition
	return r
}

// Error sets the error message for the rule.
func (r ListenAddressRule) Error(message string) ListenAddressRule {
	r.err = r.err.SetMessage(message)
	return r
}

// ErrorObject sets the error struct for the rule.
func (r ListenAddressRule) ErrorObject(err validation.Error) ListenAddressRule {
	r.err = err
	return r
}

var addrRx = regexp.MustCompile("^[a-zA-Z]+://")
