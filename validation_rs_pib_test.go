// Copyright 2021-2026 Altessa Solutions Inc. All rights reserved.
// Use of this source code is governed by license that can be found in
// the LICENSE file.

package ozzo_rules

import (
	"testing"

	"github.com/go-ozzo/ozzo-validation/v4"
)

type rsPibTest struct {
	Value any
}

func (pt *rsPibTest) Validate() error {
	return validation.ValidateStruct(pt,
		validation.Field(&pt.Value, RsPIB()),
	)
}

func TestRsPIBRule_Validate(t *testing.T) {
	validStrPIB := "112932542"
	validIntPIB := int64(112932542)

	invalidStrPIB := "112932543"
	invalidIntPIB := int64(112932543)

	rsPibTestCases := []struct {
		value     *rsPibTest
		expectErr bool
		name      string
	}{
		{&rsPibTest{Value: "112932542"}, false, "Valid RsPIB 1"},
		{&rsPibTest{Value: int64(112932542)}, false, "Valid RsPIB 2"},
		{&rsPibTest{Value: &validStrPIB}, false, "Valid RsPIB 3"},
		{&rsPibTest{Value: &validIntPIB}, false, "Valid RsPIB 4"},

		{&rsPibTest{Value: "123456789"}, true, "Invalid RsPIB 1"},
		{&rsPibTest{Value: "112932543"}, true, "Invalid RsPIB 2"},
		{&rsPibTest{Value: int64(123456789)}, true, "Invalid RsPIB 3"},
		{&rsPibTest{Value: int64(112932543)}, true, "Invalid RsPIB 4"},
		{&rsPibTest{Value: &invalidIntPIB}, true, "Invalid RsPIB 5"},
		{&rsPibTest{Value: &invalidStrPIB}, true, "Invalid RsPIB 6"},
	}

	for _, tt := range rsPibTestCases {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.value.Validate()

			if tt.expectErr && err == nil {
				t.Errorf("Expected error, got nil")
			}

			if !tt.expectErr && err != nil {
				t.Errorf("Expected no error, got %v", err)
			}
		})
	}
}
