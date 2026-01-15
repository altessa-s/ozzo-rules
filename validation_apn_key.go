// Copyright 2021-2026 Altessa Solutions Inc. All rights reserved.
// Use of this source code is governed by license that can be found in
// the LICENSE file.

package ozzo_rules

import (
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/pem"

	"github.com/go-ozzo/ozzo-validation/v4"
)

// ErrApnKeyInvalid is the error that returns when a value is not a valid apn key.
var ErrApnKeyInvalid = validation.NewError("validation_apn_key",
	"must be a data in p8 format")

// ApnKey is a validation rule that checks if a value is a valid apn key.
func ApnKey() ApnKeyRule {
	return ApnKeyRule{condition: true, err: ErrApnKeyInvalid}
}

// ApnKeyRule is a rule that checks if a value is a valid apn key.
type ApnKeyRule struct {
	err       validation.Error
	condition bool
}

// Validate checks if the given value is valid or not.
func (r ApnKeyRule) Validate(v any) error {
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

	block, _ := pem.Decode([]byte(str))
	if block == nil {
		return r.err
	}

	key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err == nil && key != nil {
		switch key.(type) { //nolint:gocritic
		case *ecdsa.PrivateKey:
			return nil
		}
	}
	return r.err
}

// When sets the condition that determines if the validation should be performed.
func (r ApnKeyRule) When(condition bool) ApnKeyRule {
	r.condition = condition
	return r
}

// Error sets the error message for the rule.
func (r ApnKeyRule) Error(message string) ApnKeyRule {
	r.err = r.err.SetMessage(message)
	return r
}

// ErrorObject sets the error struct for the rule.
func (r ApnKeyRule) ErrorObject(err validation.Error) ApnKeyRule {
	r.err = err
	return r
}
