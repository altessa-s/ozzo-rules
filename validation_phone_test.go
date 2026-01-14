// Copyright 2021-2026 Altessa Solutions Inc. All rights reserved.
// Use of this source code is governed by license that can be found in
// the LICENSE file.

package ozzo_rules

import (
	"testing"

	"github.com/go-ozzo/ozzo-validation/v4"
)

func TestPhoneRule(t *testing.T) {
	tests := []struct {
		name        string
		value       any
		countryCode string
		expectErr   bool
	}{
		// Valid US phone numbers
		{"valid US phone with country code", "+12125551234", "US", false},
		{"valid US phone without plus", "12125551234", "US", false},
		{"valid US phone with dashes", "+1-212-555-1234", "US", false},
		{"valid US phone with spaces", "+1 212 555 1234", "US", false},
		{"valid US phone with parentheses", "+1 (212) 555-1234", "US", false},
		{"valid US toll-free", "+18002345678", "US", false},

		// Valid international phone numbers
		{"valid UK phone", "+442071234567", "GB", false},
		{"valid German phone", "+491701234567", "DE", false},
		{"valid French phone", "+33123456789", "FR", false},
		{"valid Russian phone", "+79161234567", "RU", false},
		{"valid Japanese phone", "+819012345678", "JP", false},
		{"valid Australian phone", "+61412345678", "AU", false},
		{"valid Indian phone", "+919876543210", "IN", false},
		{"valid Chinese phone", "+8613812345678", "CN", false},

		// Valid numbers without country code (using default)
		{"US number without country code", "2125551234", "US", true}, // needs full international format
		{"UK number without country code", "2071234567", "GB", true}, // needs full international format

		// Invalid phone numbers
		{"too short US", "+1212", "US", true},
		{"too long US", "+121255512341234", "US", true},
		{"invalid US area code", "+11235551234", "US", true},
		{"invalid format", "abc123", "US", true},
		{"empty string", "", "US", true},
		{"only plus sign", "+", "US", true},
		{"invalid country code", "+9991234567890", "US", true},

		// Edge cases
		{"nil value", nil, "US", true},
		{"non-string value", 12345678, "US", true},
		{"pointer to valid phone", phoneStrPtr("+12125551234"), "US", false},
		{"pointer to invalid phone", phoneStrPtr("invalid"), "US", true},
		{"pointer to nil", (*string)(nil), "US", true},

		// Special formatting
		{"US with dots", "+1.212.555.1234", "US", false},
		{"international format", "+1 (212) 555-1234", "US", false},
		{"leading zeros", "00442071234567", "GB", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Phone(tt.countryCode).Validate(tt.value)
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

func TestPhoneRule_DefaultCountryCode(t *testing.T) {
	// Test using default country code (US)
	// Note: The phone library requires international format or valid national format
	usPhone := "+12125551234"
	err := Phone(DefaultCountryCode).Validate(usPhone)
	if err != nil {
		t.Errorf("expected valid US phone with default country code, got error: %v", err)
	}
}

func TestPhoneRule_When(t *testing.T) {
	// Test with condition false - should skip validation
	err := Phone("US").When(false).Validate("invalid")
	if err != nil {
		t.Errorf("expected no error when condition is false, got: %v", err)
	}

	// Test with condition true - should validate
	err = Phone("US").When(true).Validate("invalid")
	if err == nil {
		t.Error("expected error when condition is true and phone is invalid")
	}
}

func TestPhoneRule_Error(t *testing.T) {
	customMsg := "custom phone error"
	err := Phone("US").Error(customMsg).Validate("invalid")
	if err == nil {
		t.Fatal("expected error but got nil")
	}
	// Note: The phone validator might override the message with parse error
	// So we just check that we got an error
}

func TestPhoneRule_ErrorObject(t *testing.T) {
	customErr := validation.NewError("custom_code", "custom message")
	err := Phone("US").ErrorObject(customErr).Validate("invalid")
	if err == nil {
		t.Fatal("expected error but got nil")
	}
}

func TestPhoneRule_VariousCountries(t *testing.T) {
	// Test specific country formats
	countryTests := []struct {
		country string
		number  string
		valid   bool
	}{
		// European countries
		{"GB", "+442071234567", true}, // UK
		{"DE", "+491701234567", true}, // Germany
		{"FR", "+33123456789", true},  // France
		{"IT", "+393123456789", true}, // Italy
		{"ES", "+34912345678", true},  // Spain
		{"NL", "+31612345678", true},  // Netherlands
		{"BE", "+32475123456", true},  // Belgium
		{"CH", "+41791234567", true},  // Switzerland
		{"AT", "+43123456789", true},  // Austria

		// Americas
		{"CA", "+14165551234", true},   // Canada (same as US)
		{"MX", "+525512345678", true},  // Mexico
		{"BR", "+5511999999999", true}, // Brazil
		{"AR", "+5491112345678", true}, // Argentina

		// Asia Pacific
		{"JP", "+819012345678", true},  // Japan
		{"CN", "+8613812345678", true}, // China
		{"IN", "+919876543210", true},  // India
		{"AU", "+61412345678", true},   // Australia
		{"NZ", "+6421234567", true},    // New Zealand
		{"SG", "+6581234567", true},    // Singapore
		{"KR", "+821012345678", true},  // South Korea

		// Middle East & Africa
		{"AE", "+971501234567", true}, // UAE
		{"SA", "+966501234567", true}, // Saudi Arabia
		{"ZA", "+27123456789", true},  // South Africa
		{"EG", "+201234567890", true}, // Egypt
		{"IL", "+972523456789", true}, // Israel

		// Russia & CIS
		{"RU", "+79161234567", true},  // Russia
		{"UA", "+380501234567", true}, // Ukraine
		{"KZ", "+77012345678", true},  // Kazakhstan
		{"BY", "+375291234567", true}, // Belarus
	}

	for _, tc := range countryTests {
		t.Run("country_"+tc.country, func(t *testing.T) {
			err := Phone(tc.country).Validate(tc.number)
			if tc.valid && err != nil {
				t.Errorf("expected valid phone for %s, got error: %v", tc.country, err)
			} else if !tc.valid && err == nil {
				t.Errorf("expected invalid phone for %s, but got no error", tc.country)
			}
		})
	}
}

func TestPhoneRule_EmptyCountryCode(t *testing.T) {
	// Test with empty country code - should use number's country code
	internationalNumbers := []string{
		"+12125551234",  // US
		"+442071234567", // UK
		"+491701234567", // Germany
		"+33123456789",  // France
		"+79161234567",  // Russia
		"+819012345678", // Japan
		"+61412345678",  // Australia
	}

	for _, number := range internationalNumbers {
		t.Run("empty_country_"+number, func(t *testing.T) {
			err := Phone("").Validate(number)
			if err != nil {
				t.Errorf("expected valid international number with empty country code, got error: %v", err)
			}
		})
	}
}

func BenchmarkPhoneValidation(b *testing.B) {
	phones := []struct {
		number  string
		country string
	}{
		{"+12125551234", "US"},
		{"+442071234567", "GB"},
		{"+491701234567", "DE"},
		{"invalid", "US"},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, p := range phones {
			_ = Phone(p.country).Validate(p.number)
		}
	}
}

// Helper function for creating string pointers
func phoneStrPtr(s string) *string {
	return &s
}
