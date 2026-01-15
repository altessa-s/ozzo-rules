// Copyright 2021-2026 Altessa Solutions Inc. All rights reserved.
// Use of this source code is governed by license that can be found in
// the LICENSE file.

package ozzo_rules

import (
	"github.com/go-ozzo/ozzo-validation/v4"
)

// errPortOutOfRange is the error that returns when a value is not a valid port number or port out of range.
var errPortOutOfRange = validation.NewError("validation_port_out_of_range",
	"the port must be between {{.min}} and {{.max}}")

// PortWithRange is a validation rule that checks if a value is a valid port.
func PortWithRange(minPort, maxPort int) PortRangeRule {
	return PortRangeRule{
		condition: true,
		err:       errPortOutOfRange.SetParams(map[string]any{"min": minPort, "max": maxPort}),
		min:       minPort,
		max:       maxPort,
	}
}

// PortRangeRule is a rule that checks if a value is a valid port.
type PortRangeRule struct {
	err       validation.Error
	condition bool
	min       int
	max       int
}

// Validate checks if the given value is valid or not.
func (r PortRangeRule) Validate(v any) error {
	if !r.condition {
		return nil
	}

	value, isNil := validation.Indirect(v)
	if isNil {
		return r.err
	}

	p, ok := value.(int)
	if !ok {
		return r.err
	}

	minPort := 0
	if r.min > 0 {
		minPort = r.min
	}

	maxPort := 65535
	if r.max > 0 && r.max < 65535 {
		maxPort = r.max
	}

	if p < minPort || p > maxPort {
		return r.err
	}

	return nil
}

// When sets the condition that determines if the validation should be performed.
func (r PortRangeRule) When(condition bool) PortRangeRule {
	r.condition = condition
	return r
}

// Error sets the error message for the rule.
func (r PortRangeRule) Error(message string) PortRangeRule {
	r.err = r.err.SetMessage(message)
	return r
}

// ErrorObject sets the error struct for the rule.
func (r PortRangeRule) ErrorObject(err validation.Error) PortRangeRule {
	r.err = err
	return r
}
