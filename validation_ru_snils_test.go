// Copyright 2021-2026 Altessa Solutions Inc. All rights reserved.
// Use of this source code is governed by license that can be found in
// the LICENSE file.

package ozzo_rules

import (
	"testing"

	"github.com/go-ozzo/ozzo-validation/v4"
)

type ruSnilsTest struct {
	Value any
}

func (pt *ruSnilsTest) Validate() error {
	return validation.ValidateStruct(pt,
		validation.Field(&pt.Value, RuSNILS()),
	)
}

func TestRuSNILSRule_Validate(t *testing.T) {
	validStrSNILS := "783 365 845 51"
	validIntSNILS := int64(78336584551)

	invalidStrSNILS := "12323ASDS"
	invalidIntSNILS := int64(87654678923)

	stringsSNILS := []struct {
		value     *ruSnilsTest
		expectErr bool
		name      string
	}{
		{&ruSnilsTest{Value: "783 365 845 51"}, false, "Valid SNILS 1"},
		{&ruSnilsTest{Value: "81234744887"}, false, "Valid SNILS 2"},
		{&ruSnilsTest{Value: "008 533 980 55"}, false, "Valid SNILS 3"},
		{&ruSnilsTest{Value: "783 365 845 51"}, false, "Valid SNILS 4"},
		{&ruSnilsTest{Value: "78336584551"}, false, "Valid SNILS 5"},
		{&ruSnilsTest{Value: "783-365-845 51"}, false, "Valid SNILS 6"},
		{&ruSnilsTest{Value: "783365-84551"}, false, "Valid SNILS 7"},
		{&ruSnilsTest{Value: "783-365-845-51"}, false, "Valid SNILS 8"},
		{&ruSnilsTest{Value: int64(78336584551)}, false, "Valid SNILS 9"},
		{&ruSnilsTest{Value: &validStrSNILS}, false, "Valid SNILS 10"},
		{&ruSnilsTest{Value: &validIntSNILS}, false, "Valid SNILS 11"},

		{&ruSnilsTest{Value: "78-336584551"}, true, "Invalid SNILS 1"},
		{&ruSnilsTest{Value: "78-3 365 845 51"}, true, "Invalid SNILS 2"},
		{&ruSnilsTest{Value: "7833658455 1"}, true, "Invalid SNILS 3"},
		{&ruSnilsTest{Value: "7833658455-1"}, true, "Invalid SNILS 4"},
		{&ruSnilsTest{Value: "783 365 8455 1"}, true, "Invalid SNILS 5"},
		{&ruSnilsTest{Value: "783 365 845 5i"}, true, "Invalid SNILS 6"},
		{&ruSnilsTest{Value: &invalidStrSNILS}, true, "Invalid SNILS 7"},
		{&ruSnilsTest{Value: &invalidIntSNILS}, true, "Invalid SNILS 8"},
	}

	for _, tt := range stringsSNILS {
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
