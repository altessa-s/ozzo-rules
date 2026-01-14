// Copyright 2021-2026 Altessa Solutions Inc. All rights reserved.
// Use of this source code is governed by license that can be found in
// the LICENSE file.

package ozzo_rules

import (
	"testing"
	"time"

	"github.com/go-ozzo/ozzo-validation/v4"
)

func TestTimestampRule(t *testing.T) {
	now := time.Now().Unix()
	yesterday := now - 86400
	tomorrow := now + 86400

	tests := []struct {
		name      string
		value     any
		min       int64
		max       int64
		expectErr bool
	}{
		// Valid timestamps - no bounds
		{"valid current timestamp", now, 0, 0, false},
		{"valid zero timestamp", int64(0), 0, 0, false},
		{"valid future timestamp", tomorrow, 0, 0, false},
		{"valid past timestamp", yesterday, 0, 0, false},
		{"valid large timestamp", int64(9999999999), 0, 0, false},

		// Valid timestamps - with min bound
		{"valid above min", now, yesterday, 0, false},
		{"valid exactly min", yesterday, yesterday, 0, false},

		// Valid timestamps - with max bound
		{"valid below max", now, 0, tomorrow, false},
		{"valid exactly max", tomorrow, 0, tomorrow, false},

		// Valid timestamps - with both bounds
		{"valid within range", now, yesterday, tomorrow, false},
		{"valid at min in range", yesterday, yesterday, tomorrow, false},
		{"valid at max in range", tomorrow, yesterday, tomorrow, false},

		// Invalid timestamps - negative
		{"negative timestamp", int64(-1), 0, 0, true},
		{"large negative", int64(-9999999999), 0, 0, true},

		// Invalid timestamps - out of bounds
		{"below min", yesterday, now, 0, true},
		{"above max", tomorrow, 0, now, true},
		{"outside range low", yesterday - 1, yesterday, tomorrow, true},
		{"outside range high", tomorrow + 1, yesterday, tomorrow, true},

		// Invalid types
		{"string timestamp", "1234567890", 0, 0, true},
		{"float timestamp", 1234567890.5, 0, 0, true},
		{"int32 timestamp", int32(1234567890), 0, 0, true},
		{"uint64 timestamp", uint64(1234567890), 0, 0, true},

		// Edge cases
		{"nil value", nil, 0, 0, true},
		{"pointer to valid", timestampInt64Ptr(now), 0, 0, false},
		{"pointer to invalid", timestampInt64Ptr(-1), 0, 0, true},
		{"pointer to nil", (*int64)(nil), 0, 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Timestamp(tt.min, tt.max).Validate(tt.value)
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

func TestTimestampRule_When(t *testing.T) {
	err := Timestamp(0, 0).When(false).Validate(-1)
	if err != nil {
		t.Errorf("expected no error when condition is false, got: %v", err)
	}

	err = Timestamp(0, 0).When(true).Validate(-1)
	if err == nil {
		t.Error("expected error when condition is true and timestamp is invalid")
	}
}

func TestTimestampRule_Error(t *testing.T) {
	customMsg := "custom timestamp error"
	err := Timestamp(0, 0).Error(customMsg).Validate(-1)
	if err == nil {
		t.Fatal("expected error but got nil")
	}
	if err.Error() != customMsg {
		t.Errorf("expected error message %q, got %q", customMsg, err.Error())
	}
}

func TestTimestampRule_ErrorObject(t *testing.T) {
	customErr := validation.NewError("custom_code", "custom message")
	err := Timestamp(0, 0).ErrorObject(customErr).Validate(-1)
	if err == nil {
		t.Fatal("expected error but got nil")
	}
	if err.Error() != "custom message" {
		t.Errorf("expected custom error message, got: %v", err.Error())
	}
}

func TestTimestampRule_CommonScenarios(t *testing.T) {
	// Common timestamp scenarios
	now := time.Now()

	scenarios := []struct {
		name        string
		description string
		value       int64
		min         int64
		max         int64
		expectErr   bool
	}{
		// Birth date validation (between 1900 and now)
		{
			"valid birth date",
			"birth date in 1990",
			time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC).Unix(),
			time.Date(1900, 1, 1, 0, 0, 0, 0, time.UTC).Unix(),
			now.Unix(),
			false,
		},
		{
			"future birth date",
			"birth date next year",
			now.AddDate(1, 0, 0).Unix(),
			time.Date(1900, 1, 1, 0, 0, 0, 0, time.UTC).Unix(),
			now.Unix(),
			true,
		},

		// Event scheduling (future dates only)
		{
			"valid future event",
			"event next month",
			now.AddDate(0, 1, 0).Unix(),
			now.Unix(),
			0,
			false,
		},
		{
			"past event",
			"event last month",
			now.AddDate(0, -1, 0).Unix(),
			now.Unix(),
			0,
			true,
		},

		// Token expiry (within next 24 hours)
		{
			"valid token expiry",
			"expires in 1 hour",
			now.Add(1 * time.Hour).Unix(),
			now.Unix(),
			now.Add(24 * time.Hour).Unix(),
			false,
		},
		{
			"expired token",
			"expired 1 hour ago",
			now.Add(-1 * time.Hour).Unix(),
			now.Unix(),
			now.Add(24 * time.Hour).Unix(),
			true,
		},

		// Historical data (past 30 days)
		{
			"valid historical",
			"15 days ago",
			now.AddDate(0, 0, -15).Unix(),
			now.AddDate(0, 0, -30).Unix(),
			now.Unix(),
			false,
		},
		{
			"too old historical",
			"45 days ago",
			now.AddDate(0, 0, -45).Unix(),
			now.AddDate(0, 0, -30).Unix(),
			now.Unix(),
			true,
		},
	}

	for _, sc := range scenarios {
		t.Run(sc.name, func(t *testing.T) {
			err := Timestamp(sc.min, sc.max).Validate(sc.value)
			if sc.expectErr && err == nil {
				t.Errorf("%s: expected error but got nil", sc.description)
			} else if !sc.expectErr && err != nil {
				t.Errorf("%s: expected no error but got: %v", sc.description, err)
			}
		})
	}
}

func TestTimestampRule_EdgeValues(t *testing.T) {
	// Test edge values
	tests := []struct {
		name      string
		value     int64
		expectErr bool
	}{
		// Unix epoch boundaries
		{"unix epoch start", 0, false},
		{"before unix epoch", -1, true},
		{"one second after epoch", 1, false},

		// Year 2038 problem (32-bit systems)
		{"year 2038 boundary", 2147483647, false}, // Max 32-bit signed int
		{"after year 2038", 2147483648, false},    // Works with 64-bit

		// Far future
		{"year 3000", 32503680000, false},
		{"year 9999", 253402300799, false},

		// Maximum int64 (effectively no limit)
		{"near max int64", 9223372036854775807, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Timestamp(0, 0).Validate(tt.value)
			if tt.expectErr && err == nil {
				t.Error("expected error but got nil")
			} else if !tt.expectErr && err != nil {
				t.Errorf("expected no error but got: %v", err)
			}
		})
	}
}

func BenchmarkTimestampValidation(b *testing.B) {
	now := time.Now().Unix()
	timestamps := []any{
		now,
		int64(0),
		int64(-1),
		"not a timestamp",
		now + 86400,
		now - 86400,
	}

	b.Run("NoBounds", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			for _, ts := range timestamps {
				_ = Timestamp(0, 0).Validate(ts)
			}
		}
	})

	b.Run("WithBounds", func(b *testing.B) {
		min := now - 86400
		max := now + 86400
		for i := 0; i < b.N; i++ {
			for _, ts := range timestamps {
				_ = Timestamp(min, max).Validate(ts)
			}
		}
	})
}

// Helper function
func timestampInt64Ptr(i int64) *int64 {
	return &i
}
