// Copyright 2021-2026 Altessa Solutions Inc. All rights reserved.
// Use of this source code is governed by license that can be found in
// the LICENSE file.

package ozzo_rules

import (
	"strings"
	"testing"

	"github.com/go-ozzo/ozzo-validation/v4"
)

func TestOneOfRule_String(t *testing.T) {
	tests := []struct {
		name      string
		value     any
		list      []string
		expectErr bool
	}{
		// Valid string values
		{"valid string first", "apple", []string{"apple", "banana", "orange"}, false},
		{"valid string middle", "banana", []string{"apple", "banana", "orange"}, false},
		{"valid string last", "orange", []string{"apple", "banana", "orange"}, false},
		{"valid single option", "only", []string{"only"}, false},
		{"valid from many", "five", []string{"one", "two", "three", "four", "five", "six"}, false},

		// Invalid string values
		{"invalid not in list", "grape", []string{"apple", "banana", "orange"}, true},
		{"invalid empty string", "", []string{"apple", "banana", "orange"}, true},
		{"invalid case sensitive", "Apple", []string{"apple", "banana", "orange"}, true},
		{"invalid with space", "apple ", []string{"apple", "banana", "orange"}, true},
		{"invalid substring", "app", []string{"apple", "banana", "orange"}, true},

		// Edge cases
		{"nil value", nil, []string{"apple", "banana"}, true},
		{"empty list value present", "test", []string{}, true},
		{"empty string in list valid", "", []string{"", "non-empty"}, false},
		{"empty string in list invalid", "test", []string{"", "non-empty"}, true},

		// Duplicates in list (should be deduplicated)
		{"duplicates in list", "apple", []string{"apple", "apple", "banana"}, false},

		// Special characters
		{"special chars valid", "user@example.com", []string{"user@example.com", "admin@example.com"}, false},
		{"special chars invalid", "test@example.com", []string{"user@example.com", "admin@example.com"}, true},

		// Pointer handling
		{"pointer to valid", oneofStrPtr("apple"), []string{"apple", "banana"}, false},
		{"pointer to invalid", oneofStrPtr("grape"), []string{"apple", "banana"}, true},
		{"pointer to nil", (*string)(nil), []string{"apple", "banana"}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rule := OneOf(tt.list...)
			err := rule.Validate(tt.value)
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

func TestOneOfRule_Int(t *testing.T) {
	tests := []struct {
		name      string
		value     any
		list      []int
		expectErr bool
	}{
		// Valid int values
		{"valid int first", 1, []int{1, 2, 3}, false},
		{"valid int middle", 2, []int{1, 2, 3}, false},
		{"valid int last", 3, []int{1, 2, 3}, false},
		{"valid negative", -5, []int{-10, -5, 0, 5, 10}, false},
		{"valid zero", 0, []int{-1, 0, 1}, false},
		{"valid large", 1000, []int{10, 100, 1000}, false},

		// Invalid int values
		{"invalid not in list", 4, []int{1, 2, 3}, true},
		{"invalid between values", 15, []int{10, 20, 30}, true},
		{"invalid negative", -1, []int{1, 2, 3}, true},

		// Type mismatches
		{"invalid string type", "1", []int{1, 2, 3}, true},
		{"invalid float type", 1.0, []int{1, 2, 3}, true},
		{"invalid int64 type", int64(1), []int{1, 2, 3}, true},

		// Edge cases
		{"nil value", nil, []int{1, 2, 3}, true},
		{"empty list", 1, []int{}, true},

		// Duplicates in list
		{"duplicates in list", 5, []int{5, 5, 10, 10}, false},

		// Pointer handling
		{"pointer to valid", oneofIntPtr(2), []int{1, 2, 3}, false},
		{"pointer to invalid", oneofIntPtr(4), []int{1, 2, 3}, true},
		{"pointer to nil", (*int)(nil), []int{1, 2, 3}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rule := OneOf(tt.list...)
			err := rule.Validate(tt.value)
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

func TestOneOfRule_Bool(t *testing.T) {
	tests := []struct {
		name      string
		value     any
		list      []bool
		expectErr bool
	}{
		// Valid bool values
		{"valid true only", true, []bool{true}, false},
		{"valid false only", false, []bool{false}, false},
		{"valid both true", true, []bool{true, false}, false},
		{"valid both false", false, []bool{true, false}, false},

		// Invalid cases
		{"invalid true not allowed", true, []bool{false}, true},
		{"invalid false not allowed", false, []bool{true}, true},
		{"invalid empty list", true, []bool{}, true},

		// Type mismatches
		{"invalid string true", "true", []bool{true}, true},
		{"invalid int 1", 1, []bool{true}, true},
		{"invalid int 0", 0, []bool{false}, true},

		// Edge cases
		{"nil value", nil, []bool{true, false}, true},

		// Duplicates (meaningless for bool but should work)
		{"duplicates true", true, []bool{true, true, true}, false},

		// Pointer handling
		{"pointer to valid true", oneofBoolPtr(true), []bool{true, false}, false},
		{"pointer to valid false", oneofBoolPtr(false), []bool{true, false}, false},
		{"pointer to nil", (*bool)(nil), []bool{true, false}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rule := OneOf(tt.list...)
			err := rule.Validate(tt.value)
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

func TestOneOfRule_When(t *testing.T) {
	rule := OneOf("a", "b", "c")

	err := rule.When(false).Validate("invalid")
	if err != nil {
		t.Errorf("expected no error when condition is false, got: %v", err)
	}

	err = rule.When(true).Validate("invalid")
	if err == nil {
		t.Error("expected error when condition is true and value is invalid")
	}
}

func TestOneOfRule_Error(t *testing.T) {
	customMsg := "custom oneof error"
	rule := OneOf("a", "b", "c")

	err := rule.Error(customMsg).Validate("d")
	if err == nil {
		t.Fatal("expected error but got nil")
	}
	if err.Error() != customMsg {
		t.Errorf("expected error message %q, got %q", customMsg, err.Error())
	}
}

func TestOneOfRule_ErrorObject(t *testing.T) {
	customErr := validation.NewError("custom_code", "custom message")
	rule := OneOf("a", "b", "c")

	err := rule.ErrorObject(customErr).Validate("d")
	if err == nil {
		t.Fatal("expected error but got nil")
	}
	if err.Error() != "custom message" {
		t.Errorf("expected custom error message, got: %v", err.Error())
	}
}

func TestOneOfRule_CommonScenarios(t *testing.T) {
	// Test string scenarios
	t.Run("Status Values", func(t *testing.T) {
		statuses := []string{"pending", "active", "inactive", "deleted"}
		rule := OneOf(statuses...)

		validTests := []string{"pending", "active", "inactive", "deleted"}
		for _, status := range validTests {
			err := rule.Validate(status)
			if err != nil {
				t.Errorf("expected status %q to be valid, got error: %v", status, err)
			}
		}

		invalidTests := []string{"draft", "archived", "ACTIVE", ""}
		for _, status := range invalidTests {
			err := rule.Validate(status)
			if err == nil {
				t.Errorf("expected status %q to be invalid, but got no error", status)
			}
		}
	})

	// Test int scenarios
	t.Run("HTTP Status Codes", func(t *testing.T) {
		httpCodes := []int{200, 201, 204, 400, 401, 403, 404, 500}
		rule := OneOf(httpCodes...)

		validTests := []int{200, 404, 500}
		for _, code := range validTests {
			err := rule.Validate(code)
			if err != nil {
				t.Errorf("expected HTTP code %d to be valid, got error: %v", code, err)
			}
		}

		invalidTests := []int{202, 301, 418, 503}
		for _, code := range invalidTests {
			err := rule.Validate(code)
			if err == nil {
				t.Errorf("expected HTTP code %d to be invalid, but got no error", code)
			}
		}
	})

	// Test enum-like scenarios
	t.Run("User Roles", func(t *testing.T) {
		roles := []string{"admin", "editor", "viewer", "guest"}
		rule := OneOf(roles...)

		scenarios := []struct {
			role      string
			expectErr bool
		}{
			{"admin", false},
			{"editor", false},
			{"viewer", false},
			{"guest", false},
			{"superadmin", true},
			{"moderator", true},
			{"Admin", true}, // case sensitive
			{"", true},
		}

		for _, sc := range scenarios {
			err := rule.Validate(sc.role)
			if sc.expectErr && err == nil {
				t.Errorf("expected role %q to be invalid, but got no error", sc.role)
			} else if !sc.expectErr && err != nil {
				t.Errorf("expected role %q to be valid, got error: %v", sc.role, err)
			}
		}
	})

	// Test priority levels
	t.Run("Priority Levels", func(t *testing.T) {
		priorities := []int{1, 2, 3, 4, 5}
		rule := OneOf(priorities...)

		scenarios := []struct {
			priority  int
			expectErr bool
		}{
			{1, false}, // highest
			{3, false}, // medium
			{5, false}, // lowest
			{0, true},  // too low
			{6, true},  // too high
			{10, true}, // way out of range
		}

		for _, sc := range scenarios {
			err := rule.Validate(sc.priority)
			if sc.expectErr && err == nil {
				t.Errorf("expected priority %d to be invalid, but got no error", sc.priority)
			} else if !sc.expectErr && err != nil {
				t.Errorf("expected priority %d to be valid, got error: %v", sc.priority, err)
			}
		}
	})
}

func TestOneOfRule_ErrorMessage(t *testing.T) {
	// Test that error message includes the allowed values
	t.Run("String List", func(t *testing.T) {
		rule := OneOf("apple", "banana", "orange")
		err := rule.Validate("grape")
		if err == nil {
			t.Fatal("expected error but got nil")
		}

		errMsg := err.Error()
		// Should contain all allowed values in the error message
		expectedParts := []string{"apple", "banana", "orange"}
		for _, part := range expectedParts {
			if !oneofContains(errMsg, part) {
				t.Errorf("expected error message to contain %q, got: %s", part, errMsg)
			}
		}
	})

	t.Run("Int List", func(t *testing.T) {
		rule := OneOf(1, 2, 3)
		err := rule.Validate(4)
		if err == nil {
			t.Fatal("expected error but got nil")
		}

		errMsg := err.Error()
		// Should contain all allowed values in the error message
		expectedParts := []string{"1", "2", "3"}
		for _, part := range expectedParts {
			if !oneofContains(errMsg, part) {
				t.Errorf("expected error message to contain %q, got: %s", part, errMsg)
			}
		}
	})
}

func TestOneOfRule_Deduplication(t *testing.T) {
	// Test that duplicate values in the list are handled correctly
	t.Run("String Duplicates", func(t *testing.T) {
		rule := OneOf("a", "b", "a", "c", "b", "a")

		// Should still validate correctly
		err := rule.Validate("a")
		if err != nil {
			t.Errorf("expected 'a' to be valid despite duplicates, got error: %v", err)
		}

		err = rule.Validate("d")
		if err == nil {
			t.Error("expected 'd' to be invalid, but got no error")
		}
	})

	t.Run("Int Duplicates", func(t *testing.T) {
		rule := OneOf(1, 2, 1, 3, 2, 1)

		// Should still validate correctly
		err := rule.Validate(1)
		if err != nil {
			t.Errorf("expected 1 to be valid despite duplicates, got error: %v", err)
		}

		err = rule.Validate(4)
		if err == nil {
			t.Error("expected 4 to be invalid, but got no error")
		}
	})
}

func BenchmarkOneOfValidation(b *testing.B) {
	// Benchmark with different list sizes
	b.Run("Small List", func(b *testing.B) {
		rule := OneOf("a", "b", "c")
		values := []any{"a", "b", "c", "d", nil}

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			for _, v := range values {
				_ = rule.Validate(v)
			}
		}
	})

	b.Run("Medium List", func(b *testing.B) {
		options := make([]string, 50)
		for i := 0; i < 50; i++ {
			options[i] = string(rune('a'+i%26)) + string(rune('0'+i/26))
		}
		rule := OneOf(options...)

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = rule.Validate("a0")
			_ = rule.Validate("z1")
			_ = rule.Validate("invalid")
		}
	})

	b.Run("Large List", func(b *testing.B) {
		options := make([]int, 1000)
		for i := 0; i < 1000; i++ {
			options[i] = i
		}
		rule := OneOf(options...)

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = rule.Validate(1)
			_ = rule.Validate(500)
			_ = rule.Validate(999)
			_ = rule.Validate(1001)
		}
	})
}

// Helper functions
func oneofStrPtr(s string) *string {
	return &s
}

func oneofIntPtr(i int) *int {
	return &i
}

func oneofBoolPtr(b bool) *bool {
	return &b
}

func oneofContains(s, substr string) bool {
	return len(s) >= len(substr) && strings.Contains(s, substr)
}
