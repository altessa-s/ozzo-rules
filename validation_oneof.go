// Copyright 2021-2026 Altessa Solutions Inc. All rights reserved.
// Use of this source code is governed by license that can be found in
// the LICENSE file.

package ozzo_rules

import (
	"fmt"
	"strings"

	"github.com/go-ozzo/ozzo-validation/v4"
)

// OneOf creates a validation rule that checks if a value is one of the allowed values.
func OneOf[T comparable](list ...T) OneOfRule[T] {
	lmap := make(map[T]struct{}, len(list))
	for _, s := range list {
		lmap[s] = struct{}{}
	}

	strs := make([]string, 0, len(lmap))
	for k := range lmap {
		strs = append(strs, fmt.Sprintf("'%v'", k))
	}

	return OneOfRule[T]{
		err: validation.NewError("validation_oneof",
			fmt.Sprintf("must be one of: %s", strings.Join(strs, ", "))),
		condition: true,
		list:      lmap,
	}
}

// OneOfRule is a rule that checks if a value is one of the allowed values.
type OneOfRule[T comparable] struct {
	err       validation.Error
	condition bool
	list      map[T]struct{}
}

// Validate checks if the given value is valid or not.
func (r OneOfRule[T]) Validate(v any) error {
	if !r.condition {
		return nil
	}

	value, isNil := validation.Indirect(v)
	if isNil {
		return r.err
	}

	if val, ok := value.(T); ok {
		if _, ok := r.list[val]; ok {
			return nil
		}
	}

	return r.err
}

// When sets the condition that determines if the validation should be performed.
func (r OneOfRule[T]) When(condition bool) OneOfRule[T] {
	r.condition = condition
	return r
}

// Error sets the error message for the rule.
func (r OneOfRule[T]) Error(message string) OneOfRule[T] {
	r.err = r.err.SetMessage(message)
	return r
}

// ErrorObject sets the error struct for the rule.
func (r OneOfRule[T]) ErrorObject(err validation.Error) OneOfRule[T] {
	r.err = err
	return r
}
