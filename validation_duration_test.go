// Copyright 2021-2026 Altessa Solutions Inc. All rights reserved.
// Use of this source code is governed by license that can be found in
// the LICENSE file.

package ozzo_rules

import (
	"testing"
	"time"

	"github.com/go-ozzo/ozzo-validation/v4"
)

func TestDurationRule(t *testing.T) {
	tests := []struct {
		name      string
		value     any
		expectErr bool
	}{
		// Valid durations - string format
		{"valid seconds string", "10s", false},
		{"valid minutes string", "5m", false},
		{"valid hours string", "2h", false},
		{"valid milliseconds string", "500ms", false},
		{"valid complex duration", "1h30m45s", false},
		{"valid numeric string as seconds", "60", false},
		{"valid negative duration", "-5s", false},

		// Valid durations - time.Duration type
		{"valid time.Duration", 5 * time.Second, false},
		{"valid negative time.Duration", -10 * time.Minute, false},

		// Invalid durations
		{"zero duration string", "0", true}, // Duration() requires non-zero
		{"zero duration type", time.Duration(0), true},
		{"empty string", "", true},
		{"invalid format", "10x", true},
		{"invalid text", "ten seconds", true},
		{"invalid mixed", "10 seconds", true},

		// Edge cases
		{"nil value", nil, true},
		{"non-duration type", 123, true}, // int is not handled
		{"pointer to valid duration", durationStrPtr("10s"), false},
		{"pointer to invalid duration", durationStrPtr("invalid"), true},
		{"pointer to nil", (*string)(nil), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Duration().Validate(tt.value)
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

func TestDurationOrZeroRule(t *testing.T) {
	tests := []struct {
		name      string
		value     any
		expectErr bool
	}{
		// Valid durations including zero
		{"valid seconds string", "10s", false},
		{"valid minutes string", "5m", false},
		{"zero duration string", "0", false}, // DurationOrZero() allows zero
		{"zero duration type", time.Duration(0), false},
		{"zero seconds string", "0s", false},

		// Invalid durations
		{"empty string", "", true},
		{"invalid format", "abc", true},
		{"nil value", nil, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := DurationOrZero().Validate(tt.value)
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

func TestDurationRule_When(t *testing.T) {
	err := Duration().When(false).Validate("invalid")
	if err != nil {
		t.Errorf("expected no error when condition is false, got: %v", err)
	}

	err = Duration().When(true).Validate("invalid")
	if err == nil {
		t.Error("expected error when condition is true and duration is invalid")
	}
}

func TestDurationRule_Error(t *testing.T) {
	customMsg := "custom duration error"
	err := Duration().Error(customMsg).Validate("invalid")
	if err == nil {
		t.Fatal("expected error but got nil")
	}
	if err.Error() != customMsg {
		t.Errorf("expected error message %q, got %q", customMsg, err.Error())
	}
}

func TestDurationRule_ErrorObject(t *testing.T) {
	customErr := validation.NewError("custom_code", "custom message")
	err := Duration().ErrorObject(customErr).Validate("invalid")
	if err == nil {
		t.Fatal("expected error but got nil")
	}
	if err.Error() != "custom message" {
		t.Errorf("expected custom error message, got: %v", err.Error())
	}
}

func TestDurationRule_ZeroError(t *testing.T) {
	// Test the specific "must be greater 0" error
	err := Duration().Validate("0")
	if err == nil {
		t.Fatal("expected error for zero duration")
	}
	if err.Error() != "must be greater 0" {
		t.Errorf("expected 'must be greater 0' error, got: %v", err.Error())
	}
}

func TestParseDurationString(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected time.Duration
		hasError bool
	}{
		// Valid formats
		{"seconds suffix", "10s", 10 * time.Second, false},
		{"minutes suffix", "5m", 5 * time.Minute, false},
		{"hours suffix", "2h", 2 * time.Hour, false},
		{"milliseconds suffix", "500ms", 500 * time.Millisecond, false},
		{"complex duration", "1h30m", 90 * time.Minute, false},
		{"numeric only", "60", 60 * time.Second, false},
		{"zero numeric", "0", 0, false},
		{"negative with suffix", "-5s", -5 * time.Second, false},
		{"negative numeric", "-60", -60 * time.Second, false},

		// Invalid formats
		{"invalid suffix", "10x", 0, true},
		{"invalid numeric", "abc", 0, true},
		{"empty string", "", 0, false}, // parseDurationString returns 0 for empty
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := parseDurationString(tt.input)
			if tt.hasError {
				if err == nil {
					t.Error("expected error but got nil")
				}
			} else {
				if err != nil {
					t.Errorf("expected no error but got: %v", err)
				}
				if result != tt.expected {
					t.Errorf("expected %v, got %v", tt.expected, result)
				}
			}
		})
	}
}

func TestDurationRule_VariousFormats(t *testing.T) {
	// Test various duration formats
	validFormats := []struct {
		input    any
		nonZero  bool // should pass Duration() (non-zero required)
		zeroable bool // should pass DurationOrZero()
	}{
		// String formats
		{"1s", true, true},
		{"0s", false, true},
		{"10", true, true},
		{"0", false, true},
		{"-5s", true, true},
		{"1.5h", true, true},
		{"90m", true, true},
		{"1h30m45s", true, true},

		// time.Duration values
		{1 * time.Second, true, true},
		{0 * time.Second, false, true},
		{-5 * time.Minute, true, true},
		{time.Hour + 30*time.Minute, true, true},
	}

	for _, tc := range validFormats {
		name := formatValue(tc.input)

		t.Run("Duration_"+name, func(t *testing.T) {
			err := Duration().Validate(tc.input)
			if tc.nonZero && err != nil {
				t.Errorf("expected valid non-zero duration, got error: %v", err)
			} else if !tc.nonZero && err == nil {
				t.Error("expected error for zero duration")
			}
		})

		t.Run("DurationOrZero_"+name, func(t *testing.T) {
			err := DurationOrZero().Validate(tc.input)
			if tc.zeroable && err != nil {
				t.Errorf("expected valid duration (zero allowed), got error: %v", err)
			} else if !tc.zeroable && err == nil {
				t.Error("expected error")
			}
		})
	}
}

func BenchmarkDurationValidation(b *testing.B) {
	durations := []any{
		"10s",
		"5m",
		"1h30m",
		"60",
		5 * time.Second,
		"invalid",
	}

	b.Run("Duration", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			for _, d := range durations {
				_ = Duration().Validate(d)
			}
		}
	})

	b.Run("DurationOrZero", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			for _, d := range durations {
				_ = DurationOrZero().Validate(d)
			}
		}
	})
}

// Helper functions
func durationStrPtr(s string) *string {
	return &s
}

func formatValue(v any) string {
	switch val := v.(type) {
	case string:
		return val
	case time.Duration:
		return val.String()
	default:
		return "unknown"
	}
}
