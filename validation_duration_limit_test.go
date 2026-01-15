// Copyright 2021-2026 Altessa Solutions Inc. All rights reserved.
// Use of this source code is governed by license that can be found in
// the LICENSE file.

package ozzo_rules

import (
	"strings"
	"testing"
	"time"

	"github.com/go-ozzo/ozzo-validation/v4"
)

func TestDurationLimitRule(t *testing.T) {
	// Test with time.Duration min/max
	t.Run("time.Duration limits", func(t *testing.T) {
		tests := []struct {
			name      string
			value     any
			min       time.Duration
			max       time.Duration
			expectErr bool
		}{
			// Valid durations
			{"valid in range", 5 * time.Second, 1 * time.Second, 10 * time.Second, false},
			{"valid at min", 1 * time.Second, 1 * time.Second, 10 * time.Second, false},
			{"valid at max", 10 * time.Second, 1 * time.Second, 10 * time.Second, false},
			{"valid string seconds", "5s", 1 * time.Second, 10 * time.Second, false},
			{"valid string minutes", "2m", 1 * time.Minute, 5 * time.Minute, false},
			{"valid string hours", "1h", 30 * time.Minute, 2 * time.Hour, false},
			{"valid int64 seconds", int64(5), 1 * time.Second, 10 * time.Second, false},

			// Invalid durations - out of range
			{"invalid below min", 500 * time.Millisecond, 1 * time.Second, 10 * time.Second, true},
			{"invalid above max", 11 * time.Second, 1 * time.Second, 10 * time.Second, true},
			{"invalid string below", "500ms", 1 * time.Second, 10 * time.Second, true},
			{"invalid string above", "11s", 1 * time.Second, 10 * time.Second, true},
			{"invalid int64 below", int64(0), 1 * time.Second, 10 * time.Second, true},
			{"invalid int64 above", int64(11), 1 * time.Second, 10 * time.Second, true},

			// Invalid types
			{"invalid float", 5.5, 1 * time.Second, 10 * time.Second, true},
			{"invalid bool", true, 1 * time.Second, 10 * time.Second, true},
			{"invalid slice", []int{5}, 1 * time.Second, 10 * time.Second, true},

			// Edge cases
			{"nil value", nil, 1 * time.Second, 10 * time.Second, true},
			{"invalid string format", "invalid", 1 * time.Second, 10 * time.Second, true},
			{"empty string", "", 1 * time.Second, 10 * time.Second, true},

			// Zero duration handling
			{"zero duration with positive min", 0 * time.Second, 1 * time.Second, 10 * time.Second, true},
			{"zero duration allowed", 0 * time.Second, 0 * time.Second, 10 * time.Second, false},

			// Negative duration handling (parseDuration returns -1 for errors)
			{"negative duration", -5 * time.Second, 0 * time.Second, 10 * time.Second, true},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				rule := DurationWithLimit(tt.min, tt.max)
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
	})

	// Test with string min/max
	t.Run("string limits", func(t *testing.T) {
		tests := []struct {
			name      string
			value     any
			min       string
			max       string
			expectErr bool
		}{
			// Valid durations
			{"valid in range", "5s", "1s", "10s", false},
			{"valid at min", "1m", "1m", "5m", false},
			{"valid at max", "2h", "1h", "2h", false},
			{"valid duration value", 5 * time.Second, "1s", "10s", false},
			{"valid int64", int64(30), "10s", "1m", false},

			// Invalid durations
			{"invalid below min", "500ms", "1s", "10s", true},
			{"invalid above max", "2m", "10s", "1m", true},

			// Complex duration strings
			{"valid complex", "1h30m", "1h", "2h", false},
			{"valid microseconds", "100µs", "50µs", "200µs", false},
			{"valid nanoseconds", "500ns", "100ns", "1µs", false},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				rule := DurationWithLimit(tt.min, tt.max)
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
	})

	// Test with int64 min/max (seconds)
	t.Run("int64 limits", func(t *testing.T) {
		tests := []struct {
			name      string
			value     any
			min       int64
			max       int64
			expectErr bool
		}{
			// Valid durations
			{"valid in range", int64(5), int64(1), int64(10), false},
			{"valid duration", 5 * time.Second, int64(1), int64(10), false},
			{"valid string", "5s", int64(1), int64(10), false},
			{"valid at boundaries", int64(10), int64(10), int64(10), false},

			// Invalid durations
			{"invalid below", int64(0), int64(1), int64(10), true},
			{"invalid above", int64(11), int64(1), int64(10), true},
			{"invalid duration below", 500 * time.Millisecond, int64(1), int64(10), true},
			{"invalid string above", "11s", int64(1), int64(10), true},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				rule := DurationWithLimit(tt.min, tt.max)
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
	})
}

func TestDurationLimitRule_When(t *testing.T) {
	rule := DurationWithLimit(1*time.Second, 10*time.Second)

	err := rule.When(false).Validate("invalid")
	if err != nil {
		t.Errorf("expected no error when condition is false, got: %v", err)
	}

	err = rule.When(true).Validate("invalid")
	if err == nil {
		t.Error("expected error when condition is true and duration is invalid")
	}
}

func TestDurationLimitRule_Error(t *testing.T) {
	customMsg := "custom duration error"
	rule := DurationWithLimit(1*time.Second, 10*time.Second)

	err := rule.Error(customMsg).Validate("invalid")
	if err == nil {
		t.Fatal("expected error but got nil")
	}
	if err.Error() != customMsg {
		t.Errorf("expected error message %q, got %q", customMsg, err.Error())
	}
}

func TestDurationLimitRule_ErrorObject(t *testing.T) {
	customErr := validation.NewError("custom_code", "custom message")
	rule := DurationWithLimit(1*time.Second, 10*time.Second)

	err := rule.ErrorObject(customErr).Validate("invalid")
	if err == nil {
		t.Fatal("expected error but got nil")
	}
	if err.Error() != "custom message" {
		t.Errorf("expected custom error message, got: %v", err.Error())
	}
}

func TestDurationLimitRule_CommonScenarios(t *testing.T) {
	scenarios := []struct {
		name        string
		description string
		value       any
		min         any
		max         any
		expectErr   bool
	}{
		// API timeout scenarios
		{
			"API request timeout",
			"typical HTTP request timeout",
			"30s",
			"1s",
			"5m",
			false,
		},
		{
			"API timeout too short",
			"timeout below minimum",
			"500ms",
			"1s",
			"5m",
			true,
		},
		{
			"API timeout too long",
			"timeout above maximum",
			"10m",
			"1s",
			"5m",
			true,
		},

		// Cache TTL scenarios
		{
			"Cache TTL valid",
			"typical cache TTL",
			1 * time.Hour,
			5 * time.Minute,
			24 * time.Hour,
			false,
		},
		{
			"Cache TTL too short",
			"TTL below minimum",
			2 * time.Minute,
			5 * time.Minute,
			24 * time.Hour,
			true,
		},

		// Retry delay scenarios
		{
			"Retry delay valid",
			"exponential backoff delay",
			int64(5),
			int64(1),
			int64(60),
			false,
		},
		{
			"Retry delay zero",
			"no delay between retries",
			int64(0),
			int64(1),
			int64(60),
			true,
		},

		// Session timeout scenarios
		{
			"Session timeout valid",
			"user session duration",
			"30m",
			"5m",
			"2h",
			false,
		},
		{
			"Session timeout infinite",
			"very long session",
			"24h",
			"5m",
			"2h",
			true,
		},

		// Rate limit window scenarios
		{
			"Rate limit window",
			"API rate limit window",
			"1m",
			"10s",
			"5m",
			false,
		},
		{
			"Rate limit too narrow",
			"window too small",
			"5s",
			"10s",
			"5m",
			true,
		},
	}

	for _, sc := range scenarios {
		t.Run(sc.name, func(t *testing.T) {
			var err error
			switch min := sc.min.(type) {
			case string:
				rule := DurationWithLimit(min, sc.max.(string))
				err = rule.Validate(sc.value)
			case time.Duration:
				rule := DurationWithLimit(min, sc.max.(time.Duration))
				err = rule.Validate(sc.value)
			case int64:
				rule := DurationWithLimit(min, sc.max.(int64))
				err = rule.Validate(sc.value)
			}

			if sc.expectErr && err == nil {
				t.Errorf("%s: expected error but got nil", sc.description)
			} else if !sc.expectErr && err != nil {
				t.Errorf("%s: expected no error but got: %v", sc.description, err)
			}
		})
	}
}

func TestDurationLimitRule_EdgeCases(t *testing.T) {
	// Test edge cases for duration parsing
	tests := []struct {
		name      string
		value     any
		min       time.Duration
		max       time.Duration
		expectErr bool
	}{
		// Pointer handling
		{"pointer to duration", durationLimitDurationPtr(5 * time.Second), 1 * time.Second, 10 * time.Second, false},
		{"pointer to string", durationLimitStrPtr("5s"), 1 * time.Second, 10 * time.Second, false},
		{"pointer to int64", durationLimitInt64Ptr(5), 1 * time.Second, 10 * time.Second, false},
		{"pointer to nil duration", (*time.Duration)(nil), 1 * time.Second, 10 * time.Second, true},
		{"pointer to nil string", (*string)(nil), 1 * time.Second, 10 * time.Second, true},
		{"pointer to nil int64", (*int64)(nil), 1 * time.Second, 10 * time.Second, true},

		// Boundary conditions
		{"exactly zero with zero min", 0 * time.Second, 0 * time.Second, 10 * time.Second, false},
		{"negative min ignored", 5 * time.Second, -1 * time.Second, 10 * time.Second, false},
		{"negative max ignored", 5 * time.Second, 1 * time.Second, -1 * time.Second, false},

		// String format variations
		{"fractional seconds", "1.5s", 1 * time.Second, 2 * time.Second, false},
		{"multiple units", "1m30s", 1 * time.Minute, 2 * time.Minute, false},
		{"space in string", "1 s", 1 * time.Second, 10 * time.Second, true}, // Invalid format
		{"uppercase unit", "5S", 1 * time.Second, 10 * time.Second, true},   // Invalid format
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rule := DurationWithLimit(tt.min, tt.max)
			err := rule.Validate(tt.value)
			if tt.expectErr && err == nil {
				t.Errorf("expected error but got nil")
			} else if !tt.expectErr && err != nil {
				t.Errorf("expected no error but got: %v", err)
			}
		})
	}
}

func TestDurationLimitRule_ErrorMessage(t *testing.T) {
	// Test that error messages contain the correct min/max values
	tests := []struct {
		name     string
		min      any
		max      any
		value    any
		contains []string
	}{
		{
			"shows duration format",
			1 * time.Second,
			10 * time.Second,
			"20s",
			[]string{"1s", "10s"},
		},
		{
			"shows string format",
			"5m",
			"1h",
			"2h",
			[]string{"5m", "1h"},
		},
		{
			"shows int64 as seconds",
			int64(60),
			int64(300),
			int64(400),
			[]string{"1m", "5m"}, // 60s = 1m, 300s = 5m
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var err error
			switch min := tt.min.(type) {
			case string:
				rule := DurationWithLimit(min, tt.max.(string))
				err = rule.Validate(tt.value)
			case time.Duration:
				rule := DurationWithLimit(min, tt.max.(time.Duration))
				err = rule.Validate(tt.value)
			case int64:
				rule := DurationWithLimit(min, tt.max.(int64))
				err = rule.Validate(tt.value)
			}

			if err == nil {
				t.Fatal("expected error but got nil")
			}

			errMsg := err.Error()
			for _, expected := range tt.contains {
				if !durationLimitContains(errMsg, expected) {
					t.Errorf("expected error message to contain %q, got: %s", expected, errMsg)
				}
			}
		})
	}
}

func BenchmarkDurationLimitValidation(b *testing.B) {
	rule := DurationWithLimit(1*time.Second, 1*time.Minute)
	values := []any{
		5 * time.Second,
		"30s",
		int64(45),
		"invalid",
		100 * time.Second,
		nil,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, v := range values {
			_ = rule.Validate(v)
		}
	}
}

// Helper functions
func durationLimitDurationPtr(d time.Duration) *time.Duration {
	return &d
}

func durationLimitStrPtr(s string) *string {
	return &s
}

func durationLimitInt64Ptr(i int64) *int64 {
	return &i
}

func durationLimitContains(s, substr string) bool {
	return len(s) >= len(substr) && strings.Contains(s, substr)
}
