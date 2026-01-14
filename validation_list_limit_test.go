// Copyright 2021-2026 Altessa Solutions Inc. All rights reserved.
// Use of this source code is governed by license that can be found in
// the LICENSE file.

package ozzo_rules

import (
	"testing"

	"github.com/go-ozzo/ozzo-validation/v4"
)

func TestListLimitRule(t *testing.T) {
	tests := []struct {
		name      string
		value     any
		expectErr bool
	}{
		// Valid limits
		{"valid limit 1", int64(1), false},
		{"valid limit 10", int64(10), false},
		{"valid limit 50", int64(50), false},
		{"valid limit 100", int64(100), false},
		{"valid limit 250", int64(250), false},
		{"valid limit 500", int64(500), false},
		{"valid limit 999", int64(999), false},
		{"valid limit 1000", int64(1000), false},

		// Invalid limits - out of range
		{"invalid limit 0", int64(0), true},
		{"invalid limit -1", int64(-1), true},
		{"invalid limit -10", int64(-10), true},
		{"invalid limit -100", int64(-100), true},
		{"invalid limit 1001", int64(1001), true},
		{"invalid limit 2000", int64(2000), true},
		{"invalid limit 10000", int64(10000), true},
		{"invalid limit max int64", int64(9223372036854775807), true},

		// Invalid types (only int64 is accepted)
		{"invalid int", 10, true},
		{"invalid int32", int32(10), true},
		{"invalid uint", uint(10), true},
		{"invalid uint64", uint64(10), true},
		{"invalid float32", float32(10.0), true},
		{"invalid float64", float64(10.0), true},
		{"invalid string", "10", true},
		{"invalid string text", "ten", true},
		{"invalid bool", true, true},
		{"invalid slice", []int{10}, true},
		{"invalid map", map[string]int{"limit": 10}, true},

		// Edge cases
		{"nil value", nil, true},
		{"pointer to valid", listLimitInt64Ptr(10), false},
		{"pointer to invalid", listLimitInt64Ptr(0), true},
		{"pointer to nil", (*int64)(nil), true},

		// Boundary values
		{"boundary min valid", int64(1), false},
		{"boundary max valid", int64(1000), false},
		{"boundary min invalid", int64(0), true},
		{"boundary max invalid", int64(1001), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ListLimit().Validate(tt.value)
			if tt.expectErr {
				if err == nil {
					t.Errorf("expected error but got nil")
				}
			} else {
				if err != nil {
					t.Errorf("expected no error but got: %v", err)
				}
			}
		})
	}
}

func TestListLimitRule_When(t *testing.T) {
	err := ListLimit().When(false).Validate(int64(0))
	if err != nil {
		t.Errorf("expected no error when condition is false, got: %v", err)
	}

	err = ListLimit().When(true).Validate(int64(0))
	if err == nil {
		t.Error("expected error when condition is true and limit is invalid")
	}
}

func TestListLimitRule_Error(t *testing.T) {
	customMsg := "custom limit error"
	err := ListLimit().Error(customMsg).Validate(int64(0))
	if err == nil {
		t.Fatal("expected error but got nil")
	}
	if err.Error() != customMsg {
		t.Errorf("expected error message %q, got %q", customMsg, err.Error())
	}
}

func TestListLimitRule_ErrorObject(t *testing.T) {
	customErr := validation.NewError("custom_code", "custom message")
	err := ListLimit().ErrorObject(customErr).Validate(int64(0))
	if err == nil {
		t.Fatal("expected error but got nil")
	}
	if err.Error() != "custom message" {
		t.Errorf("expected custom error message, got: %v", err.Error())
	}
}

func TestListLimitRule_CommonScenarios(t *testing.T) {
	scenarios := []struct {
		name        string
		description string
		limit       int64
		expectErr   bool
	}{
		// API pagination scenarios
		{
			"API default page size",
			"common default page size for APIs",
			int64(20),
			false,
		},
		{
			"API medium page size",
			"medium page size for data lists",
			int64(50),
			false,
		},
		{
			"API large page size",
			"large page size for bulk operations",
			int64(100),
			false,
		},
		{
			"API max allowed",
			"maximum allowed for performance",
			int64(1000),
			false,
		},
		{
			"API over limit",
			"exceeds maximum allowed",
			int64(1001),
			true,
		},

		// Database query scenarios
		{
			"DB small batch",
			"small batch for updates",
			int64(10),
			false,
		},
		{
			"DB medium batch",
			"medium batch for processing",
			int64(250),
			false,
		},
		{
			"DB large batch",
			"large batch for migrations",
			int64(500),
			false,
		},

		// UI display scenarios
		{
			"UI dropdown items",
			"items in a dropdown",
			int64(25),
			false,
		},
		{
			"UI table rows",
			"rows in a data table",
			int64(50),
			false,
		},
		{
			"UI grid items",
			"items in a grid view",
			int64(48),
			false,
		},

		// Export scenarios
		{
			"Export small",
			"small export batch",
			int64(100),
			false,
		},
		{
			"Export medium",
			"medium export batch",
			int64(500),
			false,
		},
		{
			"Export max",
			"maximum export batch",
			int64(1000),
			false,
		},

		// Invalid scenarios
		{
			"Zero items",
			"no items requested",
			int64(0),
			true,
		},
		{
			"Negative items",
			"negative number of items",
			int64(-1),
			true,
		},
		{
			"Way over limit",
			"significantly over limit",
			int64(10000),
			true,
		},
	}

	for _, sc := range scenarios {
		t.Run(sc.name, func(t *testing.T) {
			err := ListLimit().Validate(sc.limit)
			if sc.expectErr && err == nil {
				t.Errorf("%s: expected error but got nil", sc.description)
			} else if !sc.expectErr && err != nil {
				t.Errorf("%s: expected no error but got: %v", sc.description, err)
			}
		})
	}
}

func TestListLimitRule_EdgeValues(t *testing.T) {
	// Test specific edge values
	edges := []struct {
		name      string
		value     int64
		expectErr bool
	}{
		// Around zero
		{"negative large", -1000, true},
		{"negative one", -1, true},
		{"zero", 0, true},
		{"one", 1, false},
		{"two", 2, false},

		// Common pagination values
		{"ten", 10, false},
		{"twenty", 20, false},
		{"twenty-five", 25, false},
		{"fifty", 50, false},
		{"hundred", 100, false},

		// Around the maximum
		{"nine ninety eight", 998, false},
		{"nine ninety nine", 999, false},
		{"thousand", 1000, false},
		{"thousand one", 1001, true},
		{"thousand two", 1002, true},

		// Powers of 2 (common in computing)
		{"power of 2: 16", 16, false},
		{"power of 2: 32", 32, false},
		{"power of 2: 64", 64, false},
		{"power of 2: 128", 128, false},
		{"power of 2: 256", 256, false},
		{"power of 2: 512", 512, false},
		{"power of 2: 1024", 1024, true}, // Over limit

		// Common UI pagination
		{"UI: 15", 15, false},
		{"UI: 30", 30, false},
		{"UI: 60", 60, false},
		{"UI: 90", 90, false},
		{"UI: 120", 120, false},
	}

	for _, edge := range edges {
		t.Run(edge.name, func(t *testing.T) {
			err := ListLimit().Validate(edge.value)
			if edge.expectErr && err == nil {
				t.Errorf("expected error for value %d but got nil", edge.value)
			} else if !edge.expectErr && err != nil {
				t.Errorf("expected no error for value %d but got: %v", edge.value, err)
			}
		})
	}
}

func TestListLimitRule_TypeHandling(t *testing.T) {
	// Test how different numeric types are handled
	// Only int64 should be valid
	value := 10

	tests := []struct {
		name      string
		value     any
		expectErr bool
	}{
		{"int", int(value), true},
		{"int8", int8(value), true},
		{"int16", int16(value), true},
		{"int32", int32(value), true},
		{"int64", int64(value), false},
		{"uint", uint(value), true},
		{"uint8", uint8(value), true},
		{"uint16", uint16(value), true},
		{"uint32", uint32(value), true},
		{"uint64", uint64(value), true},
		{"float32", float32(value), true},
		{"float64", float64(value), true},
		{"string int", "10", true},
		{"string float", "10.0", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ListLimit().Validate(tt.value)
			if tt.expectErr && err == nil {
				t.Errorf("expected error for type %T but got nil", tt.value)
			} else if !tt.expectErr && err != nil {
				t.Errorf("expected no error for type %T but got: %v", tt.value, err)
			}
		})
	}
}

func BenchmarkListLimitValidation(b *testing.B) {
	values := []any{
		int64(1),
		int64(10),
		int64(100),
		int64(1000),
		int64(0),
		int64(1001),
		int64(-1),
		10,   // wrong type
		"10", // wrong type
		nil,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, v := range values {
			_ = ListLimit().Validate(v)
		}
	}
}

// Helper function
func listLimitInt64Ptr(i int64) *int64 {
	return &i
}
