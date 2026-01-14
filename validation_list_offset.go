// Copyright 2021-2026 Altessa Solutions Inc. All rights reserved.
// Use of this source code is governed by license that can be found in
// the LICENSE file.

package ozzo_rules

import (
	"fmt"

	"github.com/go-ozzo/ozzo-validation/v4"
)

const listMaxOffset = 1000

var ErrOffsetInvalid = validation.NewError("validation_list_offset",
	fmt.Sprintf("must be greater than 0 and less than %d", listMaxOffset))

func ListOffset() OffsetRule {
	return OffsetRule{condition: true, err: ErrOffsetInvalid}
}

type OffsetRule struct {
	err       validation.Error
	condition bool
}

func (r OffsetRule) Validate(v any) error {
	if r.condition {
		value, isNil := validation.Indirect(v)
		if isNil {
			return r.err
		}

		if o, ok := value.(int64); ok && o >= 0 && o <= listMaxOffset {
			return nil
		}
		return r.err
	}
	return nil
}

func (r OffsetRule) When(condition bool) OffsetRule {
	r.condition = condition
	return r
}

func (r OffsetRule) Error(message string) OffsetRule {
	r.err = r.err.SetMessage(message)
	return r
}

func (r OffsetRule) ErrorObject(err validation.Error) OffsetRule {
	r.err = err
	return r
}
