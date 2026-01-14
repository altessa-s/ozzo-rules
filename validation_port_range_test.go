// Copyright 2021-2026 Altessa Solutions Inc. All rights reserved.
// Use of this source code is governed by license that can be found in
// the LICENSE file.

package ozzo_rules

import (
	"testing"

	"github.com/go-ozzo/ozzo-validation/v4"
)

func TestPortRangeRule(t *testing.T) {
	tests := []struct {
		name      string
		value     any
		min       int
		max       int
		expectErr bool
	}{
		// Valid ports - default range (0-65535)
		{"valid port 80", 80, 0, 0, false},
		{"valid port 443", 443, 0, 0, false},
		{"valid port 8080", 8080, 0, 0, false},
		{"valid port 0", 0, 0, 0, false},
		{"valid port 65535", 65535, 0, 0, false},
		{"valid port 22", 22, 0, 0, false},
		{"valid port 3306", 3306, 0, 0, false},

		// Valid ports - custom min only
		{"valid above min 1024", 8080, 1024, 0, false},
		{"valid exactly min 1024", 1024, 1024, 0, false},
		{"valid high port", 50000, 1024, 0, false},

		// Valid ports - custom max only
		{"valid below max 1024", 80, 0, 1024, false},
		{"valid exactly max 1024", 1024, 0, 1024, false},
		{"valid port 0 with max", 0, 0, 1024, false},

		// Valid ports - custom range
		{"valid in range 1024-49151", 8080, 1024, 49151, false},
		{"valid at min of range", 1024, 1024, 49151, false},
		{"valid at max of range", 49151, 1024, 49151, false},
		{"valid mid range", 25000, 1024, 49151, false},

		// Invalid ports - out of default range
		{"negative port", -1, 0, 0, true},
		{"port above 65535", 65536, 0, 0, true},
		{"large negative", -9999, 0, 0, true},
		{"very large port", 99999, 0, 0, true},

		// Invalid ports - below custom min
		{"below min 1024", 80, 1024, 0, true},
		{"just below min", 1023, 1024, 0, true},
		{"zero below min", 0, 1024, 0, true},

		// Invalid ports - above custom max
		{"above max 1024", 1025, 0, 1024, true},
		{"just above max", 1025, 0, 1024, true},
		{"way above max", 9999, 0, 1024, true},

		// Invalid ports - outside custom range
		{"below custom range", 1023, 1024, 49151, true},
		{"above custom range", 49152, 1024, 49151, true},
		{"negative with range", -1, 1024, 49151, true},

		// Invalid types
		{"string port", "8080", 0, 0, true},
		{"float port", 8080.5, 0, 0, true},
		{"int64 port", int64(8080), 0, 0, true},
		{"uint port", uint(8080), 0, 0, true},
		{"bool port", true, 0, 0, true},
		{"slice port", []int{8080}, 0, 0, true},

		// Edge cases
		{"nil value", nil, 0, 0, true},
		{"pointer to valid", portRangeIntPtr(8080), 0, 0, false},
		{"pointer to invalid", portRangeIntPtr(-1), 0, 0, true},
		{"pointer to nil", (*int)(nil), 0, 0, true},

		// Special range tests
		{"privileged port", 80, 0, 1023, false},
		{"non-privileged port", 8080, 1024, 65535, false},
		{"ephemeral port", 50000, 49152, 65535, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := PortWithRange(tt.min, tt.max).Validate(tt.value)
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

func TestPortRangeRule_When(t *testing.T) {
	err := PortWithRange(1024, 65535).When(false).Validate(80)
	if err != nil {
		t.Errorf("expected no error when condition is false, got: %v", err)
	}

	err = PortWithRange(1024, 65535).When(true).Validate(80)
	if err == nil {
		t.Error("expected error when condition is true and port is out of range")
	}
}

func TestPortRangeRule_Error(t *testing.T) {
	customMsg := "custom port error"
	err := PortWithRange(1024, 65535).Error(customMsg).Validate(80)
	if err == nil {
		t.Fatal("expected error but got nil")
	}
	if err.Error() != customMsg {
		t.Errorf("expected error message %q, got %q", customMsg, err.Error())
	}
}

func TestPortRangeRule_ErrorObject(t *testing.T) {
	customErr := validation.NewError("custom_code", "custom message")
	err := PortWithRange(1024, 65535).ErrorObject(customErr).Validate(80)
	if err == nil {
		t.Fatal("expected error but got nil")
	}
	if err.Error() != "custom message" {
		t.Errorf("expected custom error message, got: %v", err.Error())
	}
}

func TestPortRangeRule_CommonRanges(t *testing.T) {
	// Test common port range scenarios
	scenarios := []struct {
		name        string
		description string
		port        int
		min         int
		max         int
		expectErr   bool
	}{
		// System/Well-known ports (0-1023)
		{
			"HTTP on privileged",
			"HTTP on port 80",
			80,
			0,
			1023,
			false,
		},
		{
			"HTTPS on privileged",
			"HTTPS on port 443",
			443,
			0,
			1023,
			false,
		},
		{
			"SSH on privileged",
			"SSH on port 22",
			22,
			0,
			1023,
			false,
		},
		{
			"custom on privileged",
			"custom service trying privileged port",
			8080,
			0,
			1023,
			true,
		},

		// Registered ports (1024-49151)
		{
			"MySQL registered",
			"MySQL on port 3306",
			3306,
			1024,
			49151,
			false,
		},
		{
			"PostgreSQL registered",
			"PostgreSQL on port 5432",
			5432,
			1024,
			49151,
			false,
		},
		{
			"development port",
			"common dev port 8080",
			8080,
			1024,
			49151,
			false,
		},
		{
			"privileged in registered",
			"trying to use port 80 in registered range",
			80,
			1024,
			49151,
			true,
		},

		// Dynamic/Ephemeral ports (49152-65535)
		{
			"ephemeral valid",
			"valid ephemeral port",
			50000,
			49152,
			65535,
			false,
		},
		{
			"ephemeral min",
			"minimum ephemeral port",
			49152,
			49152,
			65535,
			false,
		},
		{
			"ephemeral max",
			"maximum ephemeral port",
			65535,
			49152,
			65535,
			false,
		},
		{
			"non-ephemeral",
			"registered port in ephemeral range",
			8080,
			49152,
			65535,
			true,
		},

		// Custom application ranges
		{
			"microservices range",
			"gRPC service port",
			50051,
			50000,
			60000,
			false,
		},
		{
			"game server range",
			"game server port",
			27015,
			27000,
			27030,
			false,
		},
		{
			"monitoring range",
			"Prometheus port",
			9090,
			9000,
			9999,
			false,
		},

		// Edge cases
		{
			"zero with min 1",
			"port 0 when min is 1",
			0,
			1,
			65535,
			true,
		},
		{
			"max port with lower max",
			"port 65535 when max is lower",
			65535,
			0,
			65534,
			true,
		},
	}

	for _, sc := range scenarios {
		t.Run(sc.name, func(t *testing.T) {
			err := PortWithRange(sc.min, sc.max).Validate(sc.port)
			if sc.expectErr && err == nil {
				t.Errorf("%s: expected error but got nil", sc.description)
			} else if !sc.expectErr && err != nil {
				t.Errorf("%s: expected no error but got: %v", sc.description, err)
			}
		})
	}
}

func TestPortRangeRule_BoundaryValues(t *testing.T) {
	// Test boundary conditions
	tests := []struct {
		name      string
		min       int
		max       int
		testCases []struct {
			port      int
			expectErr bool
		}
	}{
		{
			"default boundaries",
			0, 0,
			[]struct {
				port      int
				expectErr bool
			}{
				{-1, true},
				{0, false},
				{1, false},
				{65534, false},
				{65535, false},
				{65536, true},
			},
		},
		{
			"custom min only",
			1024, 0,
			[]struct {
				port      int
				expectErr bool
			}{
				{0, true},
				{1023, true},
				{1024, false},
				{1025, false},
				{65535, false},
				{65536, true},
			},
		},
		{
			"custom max only",
			0, 1024,
			[]struct {
				port      int
				expectErr bool
			}{
				{-1, true},
				{0, false},
				{1023, false},
				{1024, false},
				{1025, true},
				{65535, true},
			},
		},
		{
			"narrow range",
			8080, 8090,
			[]struct {
				port      int
				expectErr bool
			}{
				{8079, true},
				{8080, false},
				{8085, false},
				{8090, false},
				{8091, true},
			},
		},
		{
			"single port range",
			8080, 8080,
			[]struct {
				port      int
				expectErr bool
			}{
				{8079, true},
				{8080, false},
				{8081, true},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for _, tc := range tt.testCases {
				err := PortWithRange(tt.min, tt.max).Validate(tc.port)
				if tc.expectErr && err == nil {
					t.Errorf("port %d: expected error but got nil", tc.port)
				} else if !tc.expectErr && err != nil {
					t.Errorf("port %d: expected no error but got: %v", tc.port, err)
				}
			}
		})
	}
}

func TestPortRangeRule_ErrorMessage(t *testing.T) {
	// Test that error messages contain the correct min/max values
	tests := []struct {
		name        string
		min         int
		max         int
		port        int
		expectInMsg []string
	}{
		{
			"shows min and max",
			1024,
			49151,
			80,
			[]string{"1024", "49151"},
		},
		{
			"shows defaults when 0",
			0,
			0,
			-1,
			[]string{"0", "0"}, // The error shows the configured values, not the applied defaults
		},
		{
			"shows custom min with default max",
			1024,
			0,
			80,
			[]string{"1024", "0"}, // The error shows the configured values, not the applied defaults
		},
		{
			"shows default min with custom max",
			0,
			1024,
			9999,
			[]string{"0", "1024"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := PortWithRange(tt.min, tt.max).Validate(tt.port)
			if err == nil {
				t.Fatal("expected error but got nil")
			}
			errMsg := err.Error()
			for _, expected := range tt.expectInMsg {
				if !portRangeContains(errMsg, expected) {
					t.Errorf("expected error message to contain %q, got: %s", expected, errMsg)
				}
			}
		})
	}
}

func BenchmarkPortRangeValidation(b *testing.B) {
	ports := []int{
		-1,
		0,
		80,
		443,
		1024,
		8080,
		49151,
		65535,
		65536,
	}

	b.Run("DefaultRange", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			for _, port := range ports {
				_ = PortWithRange(0, 0).Validate(port)
			}
		}
	})

	b.Run("CustomRange", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			for _, port := range ports {
				_ = PortWithRange(1024, 49151).Validate(port)
			}
		}
	})
}

// Helper functions
func portRangeIntPtr(i int) *int {
	return &i
}

func portRangeContains(s, substr string) bool {
	if len(s) < len(substr) {
		return false
	}
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
