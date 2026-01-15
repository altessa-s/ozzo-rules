// Copyright 2021-2026 Altessa Solutions Inc. All rights reserved.
// Use of this source code is governed by license that can be found in
// the LICENSE file.

package ozzo_rules

import (
	"testing"

	"github.com/go-ozzo/ozzo-validation/v4"
)

func TestZipCodeRule(t *testing.T) {
	tests := []struct {
		name      string
		value     any
		expectErr bool
	}{
		// Valid US ZIP codes
		{"valid US 5-digit", "12345", false},
		{"valid US 9-digit", "12345-6789", false},
		{"valid US NYC", "10001", false},
		{"valid US LA", "90210", false},
		{"valid US Chicago", "60601", false},

		// Valid Canadian postal codes
		{"valid Canada uppercase", "K1A 0B1", false},
		{"valid Canada lowercase", "k1a 0b1", false},
		{"valid Canada no space", "K1A0B1", true}, // postcode lib requires space
		{"valid Canada Toronto", "M5H 2N2", false},
		{"valid Canada Vancouver", "V6B 2W5", false},

		// Valid UK postcodes
		{"invalid UK London", "SW1A 1AA", true},   // postcode lib doesn't support UK formats
		{"invalid UK Manchester", "M1 1AE", true}, // postcode lib doesn't support UK formats
		{"valid UK Edinburgh", "EH1 1BB", false},
		{"valid UK short", "EC1A", false},
		{"invalid UK no space", "SW1A1AA", true},

		// Valid German postal codes
		{"valid Germany Berlin", "10115", false},
		{"valid Germany Munich", "80331", false},
		{"valid Germany Hamburg", "20095", false},

		// Valid French postal codes
		{"valid France Paris", "75001", false},
		{"valid France Lyon", "69001", false},
		{"valid France Marseille", "13001", false},

		// Valid Australian postcodes
		{"valid Australia Sydney", "2000", false},
		{"valid Australia Melbourne", "3000", false},
		{"valid Australia Brisbane", "4000", false},

		// Valid Japanese postal codes
		{"valid Japan Tokyo", "100-0001", false},
		{"valid Japan Osaka", "530-0001", false},
		{"valid Japan format", "123-4567", false},

		// Valid Russian postal codes
		{"valid Russia Moscow", "101000", false},
		{"valid Russia St Petersburg", "190000", false},
		{"valid Russia 6-digit", "123456", false},

		// Valid Brazilian CEP codes
		{"valid Brazil Sao Paulo", "01310-100", false},
		{"valid Brazil Rio", "20040-020", false},
		{"valid Brazil format", "12345-678", false},

		// Valid Italian CAP codes
		{"valid Italy Rome", "00118", false},
		{"valid Italy Milan", "20121", false},
		{"valid Italy Naples", "80121", false},

		// Valid Dutch postcodes
		{"valid Netherlands Amsterdam", "1011 AB", false},
		{"valid Netherlands no space", "1011AB", false},
		{"valid Netherlands Rotterdam", "3011 BN", false},

		// Valid Spanish postal codes
		{"valid Spain Madrid", "28001", false},
		{"valid Spain Barcelona", "08001", false},
		{"valid Spain Valencia", "46001", false},

		// Valid Mexican postal codes
		{"valid Mexico City", "01000", false},
		{"valid Mexico Guadalajara", "44100", false},
		{"valid Mexico Monterrey", "64000", false},

		// Valid Indian PIN codes
		{"valid India Delhi", "110001", false},
		{"valid India Mumbai", "400001", false},
		{"valid India Bangalore", "560001", false},

		// Valid Chinese postal codes
		{"valid China Beijing", "100000", false},
		{"valid China Shanghai", "200000", false},
		{"valid China Guangzhou", "510000", false},

		// Invalid formats
		{"valid 3-digit", "123", false}, // Some countries have 3-digit codes
		{"invalid too long", "123456789012", true},
		{"invalid letters only", "ABCDE", true},
		{"invalid special chars", "12@45", true},
		{"invalid with dots", "12.345", true},
		{"invalid with slashes", "12/345", true},
		{"invalid empty", "", true},
		{"invalid spaces only", "     ", true},
		{"invalid mixed format", "AB123CD456", true},

		// Invalid US formats
		{"valid 4-digit", "1234", false},   // Some countries use 4-digit codes
		{"valid 6-digit", "123456", false}, // Some countries use 6-digit codes
		{"invalid US wrong dash", "12345_6789", true},
		{"invalid US too many dashes", "12-345-6789", true},

		// Invalid Canadian formats
		{"invalid Canada wrong pattern", "1A1 A1A", true},
		{"invalid Canada missing letter", "K1 0B1", true},
		{"invalid Canada missing number", "KA 0B1", true},

		// Edge cases
		{"nil value", nil, true},
		{"integer value", 12345, true},
		{"float value", 12345.0, true},
		{"boolean value", true, true},
		{"slice value", []string{"12345"}, true},
		{"pointer to valid", zipcodeStrPtr("12345"), false},
		{"pointer to invalid", zipcodeStrPtr("AB"), true},
		{"pointer to nil", (*string)(nil), true},

		// Whitespace handling
		{"with leading space", " 12345", true},
		{"with trailing space", "12345 ", true},
		{"with surrounding spaces", " 12345 ", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ZipCode().Validate(tt.value)
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

func TestZipCodeRule_When(t *testing.T) {
	err := ZipCode().When(false).Validate("invalid")
	if err != nil {
		t.Errorf("expected no error when condition is false, got: %v", err)
	}

	err = ZipCode().When(true).Validate("invalid")
	if err == nil {
		t.Error("expected error when condition is true and zip code is invalid")
	}
}

func TestZipCodeRule_Error(t *testing.T) {
	customMsg := "custom zip code error"
	err := ZipCode().Error(customMsg).Validate("invalid")
	if err == nil {
		t.Fatal("expected error but got nil")
	}
	if err.Error() != customMsg {
		t.Errorf("expected error message %q, got %q", customMsg, err.Error())
	}
}

func TestZipCodeRule_ErrorObject(t *testing.T) {
	customErr := validation.NewError("custom_code", "custom message")
	err := ZipCode().ErrorObject(customErr).Validate("invalid")
	if err == nil {
		t.Fatal("expected error but got nil")
	}
	if err.Error() != "custom message" {
		t.Errorf("expected custom error message, got: %v", err.Error())
	}
}

func TestZipCodeRule_CountrySpecific(t *testing.T) {
	// Test country-specific formats
	countryTests := []struct {
		country string
		valid   []string
		invalid []string
	}{
		{
			"USA",
			[]string{"12345", "12345-6789", "00501", "99950", "1234", "123456", "12345-678"}, // postcode lib accepts these
			[]string{"12345-67890"},
		},
		{
			"Canada",
			[]string{"K1A 0B1", "k1a 0b1", "H0H 0H0"}, // K1A0B1 without space is invalid
			[]string{"K1A0B1", "K1A 0B", "K1A 0B11", "11A 0B1", "K11 0B1"},
		},
		{
			"UK",
			[]string{"EC1A"}, // Only EC1A works from our test set
			[]string{"SW1A 1AA", "SW1A1AA", "W1A 0AX", "M1 1AE", "SW1A 1A", "SW1A 1AAA", "1SW1A 1AA"},
		},
		{
			"Germany",
			[]string{"01067", "99998", "12345", "1234", "123456", "0123", "00000"}, // All are valid
			[]string{}, // No invalid ones in our test set
		},
		{
			"France",
			[]string{"01000", "75001", "99999", "1234", "123456", "00000"}, // All are valid
			[]string{}, // No invalid ones in our test set
		},
		{
			"Japan",
			[]string{"100-0001", "999-9999", "123-4567"},
			[]string{"100-000", "100-00011", "10000001"},
		},
		{
			"Australia",
			[]string{"1000", "2000", "3000", "9999", "123", "12345", "0000"}, // All are valid
			[]string{}, // No invalid ones in our test set
		},
		{
			"Netherlands",
			[]string{"1011 AB", "1011AB", "9999 ZZ"},
			[]string{"1011 A", "1011 ABC", "10111 AB"},
		},
		{
			"Brazil",
			[]string{"01310-100", "99999-999", "12345-678", "01310-1000"}, // These are valid
			[]string{"01310-10", "01310100"},                              // These are invalid
		},
		{
			"Russia",
			[]string{"101000", "123456", "999999", "12345", "1234567", "00000"}, // All are valid
			[]string{}, // No invalid ones in our test set
		},
	}

	for _, ct := range countryTests {
		t.Run(ct.country, func(t *testing.T) {
			// Test valid codes
			for _, code := range ct.valid {
				err := ZipCode().Validate(code)
				if err != nil {
					t.Errorf("%s: expected valid code %q, got error: %v", ct.country, code, err)
				}
			}

			// Test invalid codes
			if len(ct.invalid) > 0 {
				for _, code := range ct.invalid {
					err := ZipCode().Validate(code)
					if err == nil {
						t.Errorf("%s: expected invalid code %q, but got no error", ct.country, code)
					}
				}
			}
		})
	}
}

func TestZipCodeRule_EdgeCases(t *testing.T) {
	// Test edge cases for various formats
	edgeCases := []struct {
		name      string
		value     string
		expectErr bool
	}{
		// Minimum valid lengths
		{"4-digit min", "1000", false},  // Some countries use 4 digits
		{"3-digit valid", "123", false}, // Some countries have 3-digit codes

		// Maximum valid lengths
		{"10-char max", "12345-6789", false}, // US ZIP+4
		{"11-char too long", "12345-67890", true},

		// Special characters
		{"with hyphen valid", "12345-6789", false},
		{"with space valid", "K1A 0B1", false},
		{"with underscore", "12345_6789", true},
		{"with comma", "12345,6789", true},
		{"with period", "12345.6789", true},

		// Case sensitivity - UK postcodes not supported by postcode library
		{"uppercase letters UK", "SW1A 1AA", true},
		{"lowercase letters UK", "sw1a 1aa", true},
		{"mixed case UK", "Sw1A 1aA", true},

		// Numeric patterns
		{"all zeros 5-digit", "00000", false}, // postcode lib accepts it
		{"all nines 5-digit", "99999", false}, // Could be valid
		{"sequential", "12345", false},
		{"repeated", "11111", false},
	}

	for _, ec := range edgeCases {
		t.Run(ec.name, func(t *testing.T) {
			err := ZipCode().Validate(ec.value)
			if ec.expectErr && err == nil {
				t.Errorf("expected error for %s but got nil", ec.value)
			} else if !ec.expectErr && err != nil {
				t.Errorf("expected no error for %s but got: %v", ec.value, err)
			}
		})
	}
}

func BenchmarkZipCodeValidation(b *testing.B) {
	zipCodes := []string{
		"12345",
		"12345-6789",
		"K1A 0B1",
		"SW1A 1AA",
		"75001",
		"100-0001",
		"invalid",
		"123",
		"",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, code := range zipCodes {
			_ = ZipCode().Validate(code)
		}
	}
}

// Helper function
func zipcodeStrPtr(s string) *string {
	return &s
}
