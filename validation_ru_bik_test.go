// Copyright 2021-2026 Altessa Solutions Inc. All rights reserved.
// Use of this source code is governed by license that can be found in
// the LICENSE file.

package ozzo_rules

import (
	"testing"

	"github.com/go-ozzo/ozzo-validation/v4"
)

func TestRuBICRule(t *testing.T) {
	tests := []struct {
		name      string
		value     any
		expectErr bool
	}{
		// Valid RuBIC codes (starting with 04 + 7 digits = 9 total)
		{"valid RuBIC 1", "041234567", false},
		{"valid RuBIC 2", "049999999", false},
		{"valid RuBIC 3", "040000000", false},
		{"valid RuBIC 4", "045678901", false},
		{"valid RuBIC 5", "047654321", false},
		{"valid RuBIC all zeros after 04", "040000000", false},
		{"valid RuBIC all nines after 04", "049999999", false},

		// Real Russian bank BICs (examples)
		{"valid Sberbank Moscow", "044525225", false},
		{"valid VTB Moscow", "044525187", false},
		{"valid Alfa Bank Moscow", "044525593", false},
		{"valid Gazprombank Moscow", "044525823", false},
		{"valid Raiffeisenbank Moscow", "044525700", false},

		// Invalid RuBIC codes - wrong prefix
		{"invalid prefix 00", "001234567", true},
		{"invalid prefix 01", "011234567", true},
		{"invalid prefix 02", "021234567", true},
		{"invalid prefix 03", "031234567", true},
		{"invalid prefix 05", "051234567", true},
		{"invalid prefix 10", "101234567", true},
		{"invalid prefix 44", "441234567", true}, // Even though real BICs start with 044

		// Invalid RuBIC codes - wrong length
		{"invalid too short 8 digits", "04123456", true},
		{"invalid too short 7 digits", "0412345", true},
		{"invalid too short 6 digits", "041234", true},
		{"invalid too long 10 digits", "0412345678", true},
		{"invalid too long 11 digits", "04123456789", true},
		{"invalid empty", "", true},

		// Invalid RuBIC codes - non-numeric
		{"invalid letters", "04ABCDEFG", true},
		{"invalid mixed", "0412345AB", true},
		{"invalid special chars", "041234$67", true},
		{"invalid spaces", "041 234 567", true},
		{"invalid dashes", "041-234-567", true},
		{"invalid dots", "041.234.567", true},

		// Invalid RuBIC codes - format variations
		{"invalid with leading zero", "0041234567", true},
		{"invalid with trailing space", "041234567 ", true},
		{"invalid with leading space", " 041234567", true},

		// Edge cases
		{"nil value", nil, true},
		{"integer value", 41234567, true},
		{"integer with leading zero", 041234567, true}, // This is octal in Go
		{"float value", 41234567.0, true},
		{"boolean value", true, true},
		{"slice value", []string{"041234567"}, true},

		// Pointer handling
		{"pointer to valid", bikStrPtr("041234567"), false},
		{"pointer to invalid", bikStrPtr("001234567"), true},
		{"pointer to nil", (*string)(nil), true},

		// Boundary tests
		{"starts with 040", "040000000", false},
		{"starts with 049", "049999999", false},
		{"almost valid 039999999", "039999999", true},
		{"almost valid 050000000", "050000000", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := RuBIC().Validate(tt.value)
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

func TestRuBICRule_When(t *testing.T) {
	err := RuBIC().When(false).Validate("invalid")
	if err != nil {
		t.Errorf("expected no error when condition is false, got: %v", err)
	}

	err = RuBIC().When(true).Validate("invalid")
	if err == nil {
		t.Error("expected error when condition is true and RuBIC is invalid")
	}
}

func TestRuBICRule_Error(t *testing.T) {
	customMsg := "custom RuBIC error"
	err := RuBIC().Error(customMsg).Validate("invalid")
	if err == nil {
		t.Fatal("expected error but got nil")
	}
	if err.Error() != customMsg {
		t.Errorf("expected error message %q, got %q", customMsg, err.Error())
	}
}

func TestRuBICRule_ErrorObject(t *testing.T) {
	customErr := validation.NewError("custom_code", "custom message")
	err := RuBIC().ErrorObject(customErr).Validate("invalid")
	if err == nil {
		t.Fatal("expected error but got nil")
	}
	if err.Error() != "custom message" {
		t.Errorf("expected custom error message, got: %v", err.Error())
	}
}

func TestRuBICRule_RussianBankCodes(t *testing.T) {
	// Test actual Russian bank RuBIC codes structure
	// Russian BICs have specific meaning in their digits:
	// 04 - Russia country code
	// Next 2 digits - territory code
	// Next 3 digits - bank unit number
	// Last 2 digits - branch number

	testCases := []struct {
		name        string
		bic         string
		description string
		expectErr   bool
	}{
		// Moscow banks (territory code 45)
		{
			"Sberbank Moscow main",
			"044525225",
			"Sberbank main branch in Moscow",
			false,
		},
		{
			"Central Bank Moscow",
			"044525000",
			"Central Bank of Russia in Moscow",
			false,
		},

		// St. Petersburg banks (territory code 40)
		{
			"Bank St. Petersburg",
			"044030790",
			"Bank in St. Petersburg",
			false,
		},

		// Regional banks
		{
			"Novosibirsk bank",
			"045004867",
			"Bank in Novosibirsk region",
			false,
		},
		{
			"Yekaterinburg bank",
			"046577906",
			"Bank in Yekaterinburg",
			false,
		},

		// Invalid formats
		{
			"Old format 6 digits",
			"044525",
			"Old RuBIC format not supported",
			true,
		},
		{
			"International SWIFT",
			"SABRRUMM",
			"SWIFT code not RuBIC",
			true,
		},
		{
			"Wrong country prefix",
			"031234567",
			"Non-Russian country code",
			true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := RuBIC().Validate(tc.bic)
			if tc.expectErr && err == nil {
				t.Errorf("%s: expected error but got nil", tc.description)
			} else if !tc.expectErr && err != nil {
				t.Errorf("%s: expected no error but got: %v", tc.description, err)
			}
		})
	}
}

func TestRuBICRule_EdgeCases(t *testing.T) {
	// Test edge cases and special scenarios
	edgeCases := []struct {
		name      string
		value     string
		expectErr bool
	}{
		// Numeric patterns
		{"sequential digits", "041234567", false},
		{"repeated digits", "041111111", false},
		{"palindrome after 04", "041234321", false},

		// Almost valid patterns
		{"one digit short", "04123456", true},
		{"one digit long", "0412345678", true},
		{"wrong first digit", "141234567", true},
		{"wrong second digit", "051234567", true},

		// Format variations that might occur
		{"with dots", "04.123.4567", true},
		{"with dashes", "04-123-4567", true},
		{"with spaces", "04 123 4567", true},
		{"parentheses", "(04)1234567", true},
		{"plus prefix", "+041234567", true},

		// Case sensitivity (all digits, so not applicable)
		{"lowercase o", "o41234567", true},
		{"uppercase O", "O41234567", true},
	}

	for _, ec := range edgeCases {
		t.Run(ec.name, func(t *testing.T) {
			err := RuBIC().Validate(ec.value)
			if ec.expectErr && err == nil {
				t.Errorf("expected error for %s but got nil", ec.value)
			} else if !ec.expectErr && err != nil {
				t.Errorf("expected no error for %s but got: %v", ec.value, err)
			}
		})
	}
}

func TestRuBICRule_CommonMistakes(t *testing.T) {
	// Test common mistakes when entering RuBIC codes
	mistakes := []struct {
		name        string
		value       string
		description string
	}{
		{
			"SWIFT instead of RuBIC",
			"SABRRUMM",
			"Entered SWIFT code instead of RuBIC",
		},
		{
			"INN instead of RuBIC",
			"7707083893",
			"Entered INN (tax number) instead of RuBIC",
		},
		{
			"OGRN instead of RuBIC",
			"1027700132195",
			"Entered OGRN instead of RuBIC",
		},
		{
			"Account number",
			"40817810099910004312",
			"Entered account number instead of RuBIC",
		},
		{
			"Correspondent account",
			"30101810400000000225",
			"Entered correspondent account instead of RuBIC",
		},
		{
			"Card number",
			"4276123456789012",
			"Entered card number instead of RuBIC",
		},
		{
			"Phone number",
			"79161234567",
			"Entered phone number instead of RuBIC",
		},
	}

	for _, m := range mistakes {
		t.Run(m.name, func(t *testing.T) {
			err := RuBIC().Validate(m.value)
			if err == nil {
				t.Errorf("%s: expected error but got nil", m.description)
			}
		})
	}
}

func BenchmarkRuBICValidation(b *testing.B) {
	values := []string{
		"041234567",
		"044525225",
		"001234567",
		"04123456",
		"0412345678",
		"04ABC4567",
		"",
		"041234567890",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, v := range values {
			_ = RuBIC().Validate(v)
		}
	}
}

// Helper function
func bikStrPtr(s string) *string {
	return &s
}
