// Copyright 2021-2026 Altessa Solutions Inc. All rights reserved.
// Use of this source code is governed by license that can be found in
// the LICENSE file.

package ozzo_rules

import (
	"testing"

	"github.com/go-ozzo/ozzo-validation/v4"
)

func TestListOffsetRule(t *testing.T) {
	tests := []struct {
		name      string
		value     any
		expectErr bool
	}{
		// Valid offsets
		{"valid offset 0", int64(0), false}, // Note: offset allows 0, unlike limit
		{"valid offset 1", int64(1), false},
		{"valid offset 10", int64(10), false},
		{"valid offset 50", int64(50), false},
		{"valid offset 100", int64(100), false},
		{"valid offset 250", int64(250), false},
		{"valid offset 500", int64(500), false},
		{"valid offset 999", int64(999), false},
		{"valid offset 1000", int64(1000), false},

		// Invalid offsets - out of range
		{"invalid offset -1", int64(-1), true},
		{"invalid offset -10", int64(-10), true},
		{"invalid offset -100", int64(-100), true},
		{"invalid offset 1001", int64(1001), true},
		{"invalid offset 2000", int64(2000), true},
		{"invalid offset 10000", int64(10000), true},
		{"invalid offset max int64", int64(9223372036854775807), true},

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
		{"invalid map", map[string]int{"offset": 10}, true},

		// Edge cases
		{"nil value", nil, true},
		{"pointer to valid", listOffsetInt64Ptr(10), false},
		{"pointer to invalid", listOffsetInt64Ptr(-1), true},
		{"pointer to nil", (*int64)(nil), true},

		// Boundary values
		{"boundary min valid", int64(0), false},
		{"boundary max valid", int64(1000), false},
		{"boundary min invalid", int64(-1), true},
		{"boundary max invalid", int64(1001), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ListOffset().Validate(tt.value)
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

func TestListOffsetRule_When(t *testing.T) {
	err := ListOffset().When(false).Validate(int64(-1))
	if err != nil {
		t.Errorf("expected no error when condition is false, got: %v", err)
	}

	err = ListOffset().When(true).Validate(int64(-1))
	if err == nil {
		t.Error("expected error when condition is true and offset is invalid")
	}
}

func TestListOffsetRule_Error(t *testing.T) {
	customMsg := "custom offset error"
	err := ListOffset().Error(customMsg).Validate(int64(-1))
	if err == nil {
		t.Fatal("expected error but got nil")
	}
	if err.Error() != customMsg {
		t.Errorf("expected error message %q, got %q", customMsg, err.Error())
	}
}

func TestListOffsetRule_ErrorObject(t *testing.T) {
	customErr := validation.NewError("custom_code", "custom message")
	err := ListOffset().ErrorObject(customErr).Validate(int64(-1))
	if err == nil {
		t.Fatal("expected error but got nil")
	}
	if err.Error() != "custom message" {
		t.Errorf("expected custom error message, got: %v", err.Error())
	}
}

func TestListOffsetRule_CommonScenarios(t *testing.T) {
	scenarios := []struct {
		name        string
		description string
		offset      int64
		expectErr   bool
	}{
		// API pagination scenarios
		{
			"API first page",
			"first page with offset 0",
			int64(0),
			false,
		},
		{
			"API second page",
			"second page of 20 items",
			int64(20),
			false,
		},
		{
			"API tenth page",
			"tenth page of 20 items",
			int64(180),
			false,
		},
		{
			"API deep pagination",
			"deep pagination offset",
			int64(500),
			false,
		},
		{
			"API max pagination",
			"maximum allowed offset",
			int64(1000),
			false,
		},
		{
			"API over limit",
			"exceeds maximum offset",
			int64(1001),
			true,
		},

		// Database cursor scenarios
		{
			"DB start",
			"start of result set",
			int64(0),
			false,
		},
		{
			"DB skip 100",
			"skip first 100 records",
			int64(100),
			false,
		},
		{
			"DB skip 250",
			"skip first 250 records",
			int64(250),
			false,
		},

		// UI pagination scenarios
		{
			"UI page 1",
			"first page (offset 0)",
			int64(0),
			false,
		},
		{
			"UI page 5 of 25",
			"page 5 with 25 items per page",
			int64(100),
			false,
		},
		{
			"UI page 10 of 50",
			"page 10 with 50 items per page",
			int64(450),
			false,
		},

		// Search result scenarios
		{
			"Search start",
			"beginning of search results",
			int64(0),
			false,
		},
		{
			"Search page 2",
			"second page of search results",
			int64(10),
			false,
		},
		{
			"Search deep",
			"deep in search results",
			int64(900),
			false,
		},

		// Invalid scenarios
		{
			"Negative offset",
			"going backwards",
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
			err := ListOffset().Validate(sc.offset)
			if sc.expectErr && err == nil {
				t.Errorf("%s: expected error but got nil", sc.description)
			} else if !sc.expectErr && err != nil {
				t.Errorf("%s: expected no error but got: %v", sc.description, err)
			}
		})
	}
}

func TestListOffsetRule_PaginationCalculations(t *testing.T) {
	// Test common pagination offset calculations
	// offset = (page - 1) * limit

	testCases := []struct {
		name      string
		page      int
		limit     int
		offset    int64
		expectErr bool
	}{
		// Page size 10
		{"page 1, limit 10", 1, 10, 0, false},
		{"page 2, limit 10", 2, 10, 10, false},
		{"page 10, limit 10", 10, 10, 90, false},
		{"page 50, limit 10", 50, 10, 490, false},
		{"page 100, limit 10", 100, 10, 990, false},
		{"page 101, limit 10", 101, 10, 1000, false},
		{"page 102, limit 10", 102, 10, 1010, true}, // Over limit

		// Page size 25
		{"page 1, limit 25", 1, 25, 0, false},
		{"page 2, limit 25", 2, 25, 25, false},
		{"page 10, limit 25", 10, 25, 225, false},
		{"page 40, limit 25", 40, 25, 975, false},
		{"page 41, limit 25", 41, 25, 1000, false},
		{"page 42, limit 25", 42, 25, 1025, true}, // Over limit

		// Page size 50
		{"page 1, limit 50", 1, 50, 0, false},
		{"page 2, limit 50", 2, 50, 50, false},
		{"page 10, limit 50", 10, 50, 450, false},
		{"page 20, limit 50", 20, 50, 950, false},
		{"page 21, limit 50", 21, 50, 1000, false},
		{"page 22, limit 50", 22, 50, 1050, true}, // Over limit

		// Page size 100
		{"page 1, limit 100", 1, 100, 0, false},
		{"page 2, limit 100", 2, 100, 100, false},
		{"page 5, limit 100", 5, 100, 400, false},
		{"page 10, limit 100", 10, 100, 900, false},
		{"page 11, limit 100", 11, 100, 1000, false},
		{"page 12, limit 100", 12, 100, 1100, true}, // Over limit
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Calculate offset based on page and limit
			calculatedOffset := int64((tc.page - 1) * tc.limit)
			if calculatedOffset != tc.offset {
				t.Errorf("offset calculation mismatch: expected %d, got %d", tc.offset, calculatedOffset)
			}

			err := ListOffset().Validate(tc.offset)
			if tc.expectErr && err == nil {
				t.Errorf("expected error for offset %d but got nil", tc.offset)
			} else if !tc.expectErr && err != nil {
				t.Errorf("expected no error for offset %d but got: %v", tc.offset, err)
			}
		})
	}
}

func TestListOffsetRule_EdgeValues(t *testing.T) {
	// Test specific edge values
	edges := []struct {
		name      string
		value     int64
		expectErr bool
	}{
		// Around zero
		{"negative large", -1000, true},
		{"negative ten", -10, true},
		{"negative one", -1, true},
		{"zero", 0, false}, // Zero is valid for offset
		{"one", 1, false},
		{"two", 2, false},

		// Common offset values
		{"ten", 10, false},
		{"twenty", 20, false},
		{"fifty", 50, false},
		{"hundred", 100, false},
		{"two hundred", 200, false},

		// Around the maximum
		{"nine ninety eight", 998, false},
		{"nine ninety nine", 999, false},
		{"thousand", 1000, false},
		{"thousand one", 1001, true},
		{"thousand two", 1002, true},

		// Skip patterns
		{"skip 5", 5, false},
		{"skip 15", 15, false},
		{"skip 30", 30, false},
		{"skip 60", 60, false},
		{"skip 120", 120, false},
		{"skip 240", 240, false},
		{"skip 480", 480, false},
		{"skip 960", 960, false},
	}

	for _, edge := range edges {
		t.Run(edge.name, func(t *testing.T) {
			err := ListOffset().Validate(edge.value)
			if edge.expectErr && err == nil {
				t.Errorf("expected error for value %d but got nil", edge.value)
			} else if !edge.expectErr && err != nil {
				t.Errorf("expected no error for value %d but got: %v", edge.value, err)
			}
		})
	}
}

func TestListOffsetRule_TypeHandling(t *testing.T) {
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
			err := ListOffset().Validate(tt.value)
			if tt.expectErr && err == nil {
				t.Errorf("expected error for type %T but got nil", tt.value)
			} else if !tt.expectErr && err != nil {
				t.Errorf("expected no error for type %T but got: %v", tt.value, err)
			}
		})
	}
}

func BenchmarkListOffsetValidation(b *testing.B) {
	values := []any{
		int64(0),
		int64(10),
		int64(100),
		int64(1000),
		int64(-1),
		int64(1001),
		10,   // wrong type
		"10", // wrong type
		nil,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, v := range values {
			_ = ListOffset().Validate(v)
		}
	}
}

// Helper function
func listOffsetInt64Ptr(i int64) *int64 {
	return &i
}
