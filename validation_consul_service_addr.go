// Copyright 2021-2026 Altessa Solutions Inc. All rights reserved.
// Use of this source code is governed by license that can be found in
// the LICENSE file.

package ozzo_rules

import (
	"net"
	"regexp"

	"github.com/go-ozzo/ozzo-validation/v4"
)

// ErrConsulServiceAddressInvalid is the error that returns when a value is not a valid consul service address.
var ErrConsulServiceAddressInvalid = validation.NewError("validation_consul_service_address",
	"invalid service address, address must be valid IP address or DNS name or template")

// ConsulServiceAddress is a validation rule that checks if a value is a valid consul service address.
func ConsulServiceAddress() ConsulServiceAddressRule {
	return ConsulServiceAddressRule{condition: true, template: false, err: ErrConsulServiceAddressInvalid}
}

func ConsulServiceAddressOrTemplate() ConsulServiceAddressRule {
	return ConsulServiceAddressRule{condition: true, template: true, err: ErrConsulServiceAddressInvalid}
}

// ConsulServiceAddressRule is a rule that checks if a value is a valid consul service address.
type ConsulServiceAddressRule struct {
	err       validation.Error
	condition bool
	template  bool
}

var rxTemplate = regexp.MustCompile(`\{\{[^}]*\}\}`)

// Validate checks if the given value is valid or not.
func (r ConsulServiceAddressRule) Validate(v any) error {
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

	if r.template && rxTemplate.MatchString(str) {
		// Special case: reject templates that are only dots between braces like {{.}}.{{.}}
		if regexp.MustCompile(`^\{\{\.\}\}(\.\{\{\.\}\})*$`).MatchString(str) {
			return r.err
		}

		// Replace all templates with a valid placeholder to check the rest of the string
		simplified := rxTemplate.ReplaceAllString(str, "placeholder")

		// If the result is empty or just dots, it's invalid
		if simplified == "" || regexp.MustCompile(`^\.+$`).MatchString(simplified) {
			return r.err
		}

		// Check if simplified string is a valid address
		if ip := net.ParseIP(simplified); ip != nil && !ip.IsUnspecified() {
			return nil
		}

		if !isInvalidIPFormat(simplified) && dnsRx.MatchString(simplified) {
			return nil
		}

		// If contains templates but the structure is invalid
		return r.err
	}

	// First check if it's a valid IP address
	if ip := net.ParseIP(str); ip != nil {
		if ip.IsUnspecified() {
			return r.err
		}
		return nil
	}

	// Check if it looks like an invalid IP (contains only dots and numbers)
	if isInvalidIPFormat(str) {
		return r.err
	}

	// Finally check if it's a valid DNS name
	if !dnsRx.MatchString(str) {
		return r.err
	}

	return nil
}

// When sets the condition that determines if the validation should be performed.
func (r ConsulServiceAddressRule) When(condition bool) ConsulServiceAddressRule {
	r.condition = condition
	return r
}

// Error sets the error message for the rule.
func (r ConsulServiceAddressRule) Error(message string) ConsulServiceAddressRule {
	r.err = r.err.SetMessage(message)
	return r
}

// ErrorObject sets the error struct for the rule.
func (r ConsulServiceAddressRule) ErrorObject(err validation.Error) ConsulServiceAddressRule {
	r.err = err
	return r
}

var dnsRx = regexp.MustCompile(
	`^([a-zA-Z0-9_]([a-zA-Z0-9_\-]{0,61}[a-zA-Z0-9_])?)(\.([a-zA-Z0-9_]([a-zA-Z0-9_\-]{0,61}[a-zA-Z0-9_])?))*\.?$`,
)

// isInvalidIPFormat checks if string looks like an IP but is invalid
func isInvalidIPFormat(s string) bool {
	// Check if string contains only numbers, dots and possibly colons (for IPv6)
	ipLikeRx := regexp.MustCompile(`^[\d.:]+$`)
	if !ipLikeRx.MatchString(s) {
		return false
	}

	// If it looks like an IP but net.ParseIP returned nil, it's invalid
	return true
}
