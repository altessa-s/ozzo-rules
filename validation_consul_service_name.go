// Copyright 2021-2026 Altessa Solutions Inc. All rights reserved.
// Use of this source code is governed by license that can be found in
// the LICENSE file.

package ozzo_rules

import (
	"regexp"

	"github.com/go-ozzo/ozzo-validation/v4"
)

// ErrConsulServiceNameInvalid is the error that returns when a value is not a valid consul service name.
var ErrConsulServiceNameInvalid = validation.NewError("validation_consul_service_name",
	"invalid service name, service name must be valid per RFC 1123 and can contain only alphanumeric characters or dashes")

// ConsulServiceName is a validation rule that checks if a value is a valid consul service name.
func ConsulServiceName() ConsulServiceNameRule {
	return ConsulServiceNameRule{condition: true, allowNil: false, err: ErrConsulServiceNameInvalid}
}

func ConsulServiceNameOrNil() ConsulServiceNameRule {
	return ConsulServiceNameRule{condition: true, allowNil: true, err: ErrConsulServiceNameInvalid}
}

// ConsulServiceNameRule is a rule that checks if a value is a valid consul service name.
type ConsulServiceNameRule struct {
	err       validation.Error
	condition bool
	allowNil  bool
}

// Validate checks if the given value is valid or not.
func (r ConsulServiceNameRule) Validate(v any) error {
	if !r.condition {
		return nil
	}

	value, isNil := validation.Indirect(v)
	if isNil {
		if r.allowNil {
			return nil
		}
		return r.err
	}

	str, err := validation.EnsureString(value)
	if err == nil {
		if str == "" && r.allowNil {
			return nil
		}
	}

	if err != nil || !serviceNameRx.MatchString(str) {
		return r.err
	}

	// Additional check for double hyphens
	if regexp.MustCompile(`--`).MatchString(str) {
		return r.err
	}

	return nil
}

// When sets the condition that determines if the validation should be performed.
func (r ConsulServiceNameRule) When(condition bool) ConsulServiceNameRule {
	r.condition = condition
	return r
}

// Error sets the error message for the rule.
func (r ConsulServiceNameRule) Error(message string) ConsulServiceNameRule {
	r.err = r.err.SetMessage(message)
	return r
}

// ErrorObject sets the error struct for the rule.
func (r ConsulServiceNameRule) ErrorObject(err validation.Error) ConsulServiceNameRule {
	r.err = err
	return r
}

var serviceNameRx = regexp.MustCompile(
	`^[a-zA-Z0-9]([a-zA-Z0-9-]*[a-zA-Z0-9])?(\.[a-zA-Z0-9]([a-zA-Z0-9-]*[a-zA-Z0-9])?)*$`,
)
