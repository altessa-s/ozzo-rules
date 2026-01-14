// Copyright 2021-2026 Altessa Solutions Inc. All rights reserved.
// Use of this source code is governed by license that can be found in
// the LICENSE file.

package ozzo_rules

import (
	"testing"

	"github.com/go-ozzo/ozzo-validation/v4"
)

type rsMbrTest struct {
	Value any
}

func (mt *rsMbrTest) Validate() error {
	return validation.ValidateStruct(mt,
		validation.Field(&mt.Value, RsMBR()),
	)
}

func TestRsMbrRule_Validate(t *testing.T) {
	rsMbrTestCases := []struct {
		value     *rsMbrTest
		expectErr bool
		name      string
	}{
		{&rsMbrTest{Value: "17232037"}, false, "Valid RsMBR 1"},
		{&rsMbrTest{Value: "07899688"}, false, "Valid RsMBR 2"},
		{&rsMbrTest{Value: "29021783"}, false, "Valid RsMBR 3"},
		{&rsMbrTest{Value: "17150286"}, false, "Valid RsMBR 4"},
		{&rsMbrTest{Value: int64(17645676)}, false, "Valid RsMBR 5"},
		{&rsMbrTest{Value: int64(17337068)}, false, "Valid RsMBR 6"},
		{&rsMbrTest{Value: int64(7019904)}, false, "Valid RsMBR 7"},

		{&rsMbrTest{Value: "12345678"}, true, "Invalid RsMBR 1"},
		{&rsMbrTest{Value: "07899689"}, true, "Invalid RsMBR 2"},
		{&rsMbrTest{Value: "17150206"}, true, "Invalid RsMBR 3"},
		{&rsMbrTest{Value: "0789968915786587412598741577"}, true, "Invalid RsMBR 4"},
		{&rsMbrTest{Value: "0789bx9689"}, true, "Invalid RsMBR 5"},
		{&rsMbrTest{Value: "abvgdstryc"}, true, "Invalid RsMBR 6"},
		{&rsMbrTest{Value: "123"}, true, "Invalid RsMBR 7"},
		{&rsMbrTest{Value: "0"}, true, "Invalid RsMBR 8"},
		{&rsMbrTest{Value: ""}, true, "Invalid RsMBR 9"},
		{&rsMbrTest{Value: int64(12345678)}, true, "Invalid RsMBR 10"},
		{&rsMbrTest{Value: int64(789928)}, true, "Invalid RsMBR 11"},
		{&rsMbrTest{Value: int64(1232158778952121478)}, true, "Invalid RsMBR 12"},
		{&rsMbrTest{Value: int64(123)}, true, "Invalid RsMBR 13"},
		{&rsMbrTest{Value: int64(0)}, true, "Invalid RsMBR 14"},
	}

	for _, tt := range rsMbrTestCases {
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
