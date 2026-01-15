// Copyright 2021-2026 Altessa Solutions Inc. All rights reserved.
// Use of this source code is governed by license that can be found in
// the LICENSE file.

package ozzo_rules

import (
	"testing"

	"github.com/go-ozzo/ozzo-validation/v4"
)

// Test data for APN key validation
const (
	// Valid P8 ECDSA private key (test key, not for production use)
	validP8Key = `-----BEGIN PRIVATE KEY-----
MIGHAgEAMBMGByqGSM49AgEGCCqGSM49AwEHBG0wawIBAQQgIvlsbaXp1SRkysF9
1AQUXZxPvJIk5HdDLmJpIiKU9HehRANCAASEqMfKJvM1T8n9oYqOmJm1tMnPRNH+
hqV9gD8eEiMGMlNBrPqb8eNda1nJZDpCLW1fIhcbdCfk+FBgRvAVET1R
-----END PRIVATE KEY-----`

	// Invalid P8 key (RSA instead of ECDSA)
	invalidP8KeyRSA = `-----BEGIN PRIVATE KEY-----
MIIBVAIBADANBgkqhkiG9w0BAQEFAASCAT4wggE6AgEAAkEAwJZYY6KbMYXZJLxK
RmE+vBHwLZLwO8rJQ5mzQVQmgs4P4RmANryKEhNvbHr8TQwP2LN0u7cXOmj5V7s0
VGkfLwIDAQABAkEAuRoFLCGJA7voR+6vspbpKvDbfGx5VqaWvq0PFpZf84RM8nhO
k2F6yb07h4ztIC5zPjm8WKZdJ3pGaZlFFx8xAQIhAO0qo2YAcnMFxXdFp1PKSpnE
8Yr9yBqkJHU+tz+qg5lhAiEA0ER6IzH0pEWZ7s6TwM0FR13hTAJH8nsU1qTjJVV5
VY8CIDa+2D+/O4jfpUF7bIRZQIEiLz+0QXMWYSsKs3Mzd4RBAiBwvNJ7b7kTLnn/
xrewRx9uBm5BZIWIDzFwnQdbV3tqHwIgXIH2SJmxWV3Nf9PHVgJDMAsJEsHJC5mP
s1qWEC7IdQM=
-----END PRIVATE KEY-----`

	// Invalid format (not PEM)
	invalidP8Format = `not a valid PEM format`
)

// TestApnKeyRule tests the ApnKey validator
func TestApnKeyRule(t *testing.T) {
	tests := []struct {
		name      string
		value     any
		expectErr bool
	}{
		// Valid cases
		{"valid P8 ECDSA key", validP8Key, false},
		{"valid P8 key with extra whitespace", validP8Key + "\n\n", false},

		// Invalid cases
		{"invalid RSA key", invalidP8KeyRSA, true},
		{"invalid not PEM format", invalidP8Format, true},
		{"invalid empty string", "", true},
		{"invalid missing header", "MIGHAgEAMBMGByqGSM49AgEGCCqGSM49AwEH", true},
		{"invalid corrupted PEM", "-----BEGIN PRIVATE KEY-----\ninvalid base64\n-----END PRIVATE KEY-----", true},
		{"invalid wrong PEM type", "-----BEGIN PUBLIC KEY-----\nMIGHAgEAMBMG\n-----END PUBLIC KEY-----", true},

		// Edge cases
		{"nil value", nil, true},
		{"integer value", 12345, true},
		{"boolean value", true, true},
		{"slice value", []string{validP8Key}, true},

		// Pointer handling
		{"pointer to valid", appleStrPtr(validP8Key), false},
		{"pointer to invalid", appleStrPtr(invalidP8Format), true},
		{"pointer to nil", (*string)(nil), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ApnKey().Validate(tt.value)
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

// TestApnKeyIdRule tests the ApnKeyId validator
func TestApnKeyIdRule(t *testing.T) {
	tests := []struct {
		name      string
		value     any
		expectErr bool
	}{
		// Valid key IDs
		{"valid all uppercase letters", "ABCDEFGHIJ", false},
		{"valid all numbers", "1234567890", false},
		{"valid mixed letters and numbers", "ABC123DEF4", false},
		{"valid with dots", "A.B.C.D.E.", false},
		{"valid with hyphens", "A-B-C-D-E-", false},
		{"valid mixed special chars", "A1B2.C3-D4", false},
		{"valid real Apple format", "2Y9R4HXF34", false},

		// Invalid key IDs - wrong length
		{"invalid too short 9 chars", "ABC123DEF", true},
		{"invalid too long 11 chars", "ABC123DEF45", true},
		{"invalid empty", "", true},

		// Invalid key IDs - wrong characters
		{"invalid lowercase letters", "abcdefghij", true},
		{"invalid mixed case", "AbCdEfGhIj", true},
		{"invalid with underscore", "ABC_123_DE", true},
		{"invalid with space", "ABC 123 DE", true},
		{"invalid with slash", "ABC/123/DE", true},
		{"invalid special chars", "ABC@123#DE", true},

		// Edge cases
		{"nil value", nil, true},
		{"integer value", 1234567890, true},
		{"boolean value", true, true},
		{"slice value", []string{"ABCDEFGHIJ"}, true},

		// Pointer handling
		{"pointer to valid", appleStrPtr("2Y9R4HXF34"), false},
		{"pointer to invalid", appleStrPtr("invalid"), true},
		{"pointer to nil", (*string)(nil), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ApnKeyId().Validate(tt.value)
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

// TestApnTeamIdRule tests the ApnTeamId validator
func TestApnTeamIdRule(t *testing.T) {
	tests := []struct {
		name      string
		value     any
		expectErr bool
	}{
		// Valid team IDs (same format as key IDs)
		{"valid all uppercase letters", "ABCDEFGHIJ", false},
		{"valid all numbers", "1234567890", false},
		{"valid mixed letters and numbers", "ABC123DEF4", false},
		{"valid with dots", "A.B.C.D.E.", false},
		{"valid with hyphens", "A-B-C-D-E-", false},
		{"valid mixed special chars", "A1B2.C3-D4", false},
		{"valid real Apple format", "EQHXZ8M8AV", false},

		// Invalid team IDs - wrong length
		{"invalid too short 9 chars", "ABC123DEF", true},
		{"invalid too long 11 chars", "ABC123DEF45", true},
		{"invalid empty", "", true},

		// Invalid team IDs - wrong characters
		{"invalid lowercase letters", "abcdefghij", true},
		{"invalid mixed case", "AbCdEfGhIj", true},
		{"invalid with underscore", "ABC_123_DE", true},
		{"invalid with space", "ABC 123 DE", true},
		{"invalid with slash", "ABC/123/DE", true},
		{"invalid special chars", "ABC@123#DE", true},

		// Edge cases
		{"nil value", nil, true},
		{"integer value", 1234567890, true},
		{"boolean value", true, true},
		{"slice value", []string{"ABCDEFGHIJ"}, true},

		// Pointer handling
		{"pointer to valid", appleStrPtr("EQHXZ8M8AV"), false},
		{"pointer to invalid", appleStrPtr("invalid"), true},
		{"pointer to nil", (*string)(nil), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ApnTeamId().Validate(tt.value)
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

// TestBundleIDRule tests the BundleID validator
func TestBundleIDRule(t *testing.T) {
	tests := []struct {
		name      string
		value     any
		expectErr bool
	}{
		// Valid bundle IDs
		{"valid simple", "com.example", false},
		{"valid three parts", "com.example.app", false},
		{"valid four parts", "com.example.app.ios", false},
		{"valid with numbers", "com.example123.app456", false},
		{"valid all numbers segments", "123.456.789", false},
		{"valid long", "com.company.division.team.project.app", false},

		// Real-world examples
		{"valid Apple app", "com.apple.safari", false},
		{"valid Google app", "com.google.chrome", false},
		{"valid Facebook app", "com.facebook.messenger", false},
		{"valid game", "com.supercell.clashofclans", false},

		// Invalid bundle IDs - structure
		{"invalid single segment", "com", true},
		{"invalid no dots", "comexampleapp", true},
		{"invalid starts with dot", ".com.example", true},
		{"invalid ends with dot", "com.example.", true},
		{"invalid double dots", "com..example", true},
		{"invalid empty segment", "com.example..app", true},

		// Invalid bundle IDs - characters
		{"invalid uppercase", "com.Example.app", true},
		{"invalid uppercase COM", "COM.example.app", true},
		{"invalid hyphens", "com.example-app", true},
		{"invalid underscores", "com.example_app", true},
		{"invalid special chars", "com.example@app", true},
		{"invalid spaces", "com.example app", true},

		// Edge cases
		{"nil value", nil, true},
		{"empty string", "", true},
		{"integer value", 12345, true},
		{"boolean value", true, true},
		{"slice value", []string{"com.example"}, true},

		// Pointer handling
		{"pointer to valid", appleStrPtr("com.example.app"), false},
		{"pointer to invalid", appleStrPtr("invalid"), true},
		{"pointer to nil", (*string)(nil), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := BundleID().Validate(tt.value)
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

// Test When method for all validators
func TestAppleValidators_When(t *testing.T) {
	testCases := []struct {
		name      string
		validator validation.Rule
		value     string
	}{
		{"ApnKey", ApnKey(), "invalid"},
		{"ApnKeyId", ApnKeyId(), "invalid"},
		{"ApnTeamId", ApnTeamId(), "invalid"},
		{"BundleID", BundleID(), "invalid"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Test with When(false)
			var err error
			switch v := tc.validator.(type) {
			case ApnKeyRule:
				err = v.When(false).Validate(tc.value)
			case ApnKeyIdRule:
				err = v.When(false).Validate(tc.value)
			case ApnTeamIdRule:
				err = v.When(false).Validate(tc.value)
			case BundleIDRule:
				err = v.When(false).Validate(tc.value)
			}

			if err != nil {
				t.Errorf("expected no error when condition is false, got: %v", err)
			}
		})
	}
}

// Test Error method for all validators
func TestAppleValidators_Error(t *testing.T) {
	customMsg := "custom error message"

	testCases := []struct {
		name      string
		validator validation.Rule
		value     string
	}{
		{"ApnKey", ApnKey(), "invalid"},
		{"ApnKeyId", ApnKeyId(), "lowercase"},
		{"ApnTeamId", ApnTeamId(), "lowercase"},
		{"BundleID", BundleID(), "invalid"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var err error
			switch v := tc.validator.(type) {
			case ApnKeyRule:
				err = v.Error(customMsg).Validate(tc.value)
			case ApnKeyIdRule:
				err = v.Error(customMsg).Validate(tc.value)
			case ApnTeamIdRule:
				err = v.Error(customMsg).Validate(tc.value)
			case BundleIDRule:
				err = v.Error(customMsg).Validate(tc.value)
			}

			if err == nil {
				t.Fatal("expected error but got nil")
			}
			if err.Error() != customMsg {
				t.Errorf("expected error message %q, got %q", customMsg, err.Error())
			}
		})
	}
}

// Test ErrorObject method for all validators
func TestAppleValidators_ErrorObject(t *testing.T) {
	customErr := validation.NewError("custom_code", "custom message")

	testCases := []struct {
		name      string
		validator validation.Rule
		value     string
	}{
		{"ApnKey", ApnKey(), "invalid"},
		{"ApnKeyId", ApnKeyId(), "lowercase"},
		{"ApnTeamId", ApnTeamId(), "lowercase"},
		{"BundleID", BundleID(), "invalid"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var err error
			switch v := tc.validator.(type) {
			case ApnKeyRule:
				err = v.ErrorObject(customErr).Validate(tc.value)
			case ApnKeyIdRule:
				err = v.ErrorObject(customErr).Validate(tc.value)
			case ApnTeamIdRule:
				err = v.ErrorObject(customErr).Validate(tc.value)
			case BundleIDRule:
				err = v.ErrorObject(customErr).Validate(tc.value)
			}

			if err == nil {
				t.Fatal("expected error but got nil")
			}
			if err.Error() != "custom message" {
				t.Errorf("expected custom error message, got: %v", err.Error())
			}
		})
	}
}

// Benchmark tests
func BenchmarkApnKeyValidation(b *testing.B) {
	values := []string{
		validP8Key,
		invalidP8KeyRSA,
		invalidP8Format,
		"",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, v := range values {
			_ = ApnKey().Validate(v)
		}
	}
}

func BenchmarkApnKeyIdValidation(b *testing.B) {
	values := []string{
		"2Y9R4HXF34",
		"ABCDEFGHIJ",
		"lowercase1",
		"SHORT",
		"",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, v := range values {
			_ = ApnKeyId().Validate(v)
		}
	}
}

func BenchmarkBundleIDValidation(b *testing.B) {
	values := []string{
		"com.example.app",
		"com.apple.safari",
		"invalid",
		"com.Example.App",
		"",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, v := range values {
			_ = BundleID().Validate(v)
		}
	}
}

// Helper function
func appleStrPtr(s string) *string {
	return &s
}
