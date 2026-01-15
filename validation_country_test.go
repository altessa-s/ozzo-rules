// Copyright 2021-2026 Altessa Solutions Inc. All rights reserved.
// Use of this source code is governed by license that can be found in
// the LICENSE file.

package ozzo_rules

import (
	"strings"
	"testing"

	"github.com/go-ozzo/ozzo-validation/v4"
)

func TestCountryCode2(t *testing.T) {
	tests := []struct {
		name      string
		value     any
		expectErr bool
	}{
		// Valid country codes - Major countries
		{"valid US", "US", false},
		{"valid UK", "GB", false},
		{"valid Canada", "CA", false},
		{"valid Germany", "DE", false},
		{"valid France", "FR", false},
		{"valid Japan", "JP", false},
		{"valid China", "CN", false},
		{"valid Russia", "RU", false},
		{"valid Brazil", "BR", false},
		{"valid India", "IN", false},
		{"valid Australia", "AU", false},
		{"valid Mexico", "MX", false},

		// Valid country codes - lowercase (should work as it's converted to uppercase)
		{"valid us lowercase", "us", false},
		{"valid gb lowercase", "gb", false},
		{"valid de lowercase", "de", false},
		{"valid mixed case", "Us", false},
		{"valid mixed case 2", "gB", false},

		// Valid country codes - Europe
		{"valid Italy", "IT", false},
		{"valid Spain", "ES", false},
		{"valid Poland", "PL", false},
		{"valid Netherlands", "NL", false},
		{"valid Belgium", "BE", false},
		{"valid Switzerland", "CH", false},
		{"valid Austria", "AT", false},
		{"valid Sweden", "SE", false},
		{"valid Norway", "NO", false},
		{"valid Denmark", "DK", false},
		{"valid Finland", "FI", false},
		{"valid Greece", "GR", false},
		{"valid Portugal", "PT", false},
		{"valid Czech Republic", "CZ", false},
		{"valid Romania", "RO", false},
		{"valid Hungary", "HU", false},
		{"valid Ireland", "IE", false},

		// Valid country codes - Asia
		{"valid South Korea", "KR", false},
		{"valid Indonesia", "ID", false},
		{"valid Thailand", "TH", false},
		{"valid Singapore", "SG", false},
		{"valid Malaysia", "MY", false},
		{"valid Philippines", "PH", false},
		{"valid Vietnam", "VN", false},
		{"valid Pakistan", "PK", false},
		{"valid Bangladesh", "BD", false},
		{"valid Turkey", "TR", false},
		{"valid Saudi Arabia", "SA", false},
		{"valid UAE", "AE", false},
		{"valid Israel", "IL", false},

		// Valid country codes - Americas
		{"valid Argentina", "AR", false},
		{"valid Chile", "CL", false},
		{"valid Colombia", "CO", false},
		{"valid Peru", "PE", false},
		{"valid Venezuela", "VE", false},
		{"valid Ecuador", "EC", false},
		{"valid Uruguay", "UY", false},
		{"valid Paraguay", "PY", false},
		{"valid Bolivia", "BO", false},
		{"valid Costa Rica", "CR", false},
		{"valid Panama", "PA", false},
		{"valid Cuba", "CU", false},
		{"valid Dominican Republic", "DO", false},

		// Valid country codes - Africa
		{"valid South Africa", "ZA", false},
		{"valid Egypt", "EG", false},
		{"valid Nigeria", "NG", false},
		{"valid Kenya", "KE", false},
		{"valid Morocco", "MA", false},
		{"valid Algeria", "DZ", false},
		{"valid Tunisia", "TN", false},
		{"valid Ghana", "GH", false},
		{"valid Ethiopia", "ET", false},
		{"valid Tanzania", "TZ", false},

		// Valid country codes - Oceania
		{"valid New Zealand", "NZ", false},
		{"valid Fiji", "FJ", false},
		{"valid Papua New Guinea", "PG", false},

		// Valid country codes - Small nations
		{"valid Luxembourg", "LU", false},
		{"valid Malta", "MT", false},
		{"valid Cyprus", "CY", false},
		{"valid Iceland", "IS", false},
		{"valid Andorra", "AD", false},
		{"valid Monaco", "MC", false},
		{"valid Liechtenstein", "LI", false},
		{"valid San Marino", "SM", false},
		{"valid Vatican", "VA", false},

		// Invalid country codes
		{"invalid UK code", "UK", true}, // Common mistake - should be GB
		{"invalid single letter", "A", true},
		{"invalid three letters", "USA", true},
		{"invalid four letters", "USAA", true},
		{"invalid numeric", "12", true},
		{"invalid alphanumeric", "U1", true},
		{"invalid special chars", "U$", true},
		{"invalid with space", "U S", true},
		{"invalid with dash", "U-S", true},
		{"invalid empty", "", true},
		{"invalid spaces", "  ", true},
		{"invalid non-existent", "XX", true},
		{"invalid ZZ", "ZZ", true},
		{"invalid AA", "AA", true},

		// Edge cases
		{"nil value", nil, true},
		{"integer value", 12, true},
		{"float value", 12.34, true},
		{"boolean value", true, true},
		{"slice value", []string{"US"}, true},

		// Special territories and regions
		{"valid Hong Kong", "HK", false},
		{"valid Macau", "MO", false},
		{"valid Taiwan", "TW", false},
		{"valid Puerto Rico", "PR", false},
		{"valid Greenland", "GL", false},
		{"valid Palestine", "PS", false},
		{"invalid Kosovo", "XK", true}, // User-assigned code, not in ISO standard

		// Former countries (should be invalid)
		{"invalid Yugoslavia", "YU", true},
		{"invalid Soviet Union", "SU", true},
		{"invalid Czechoslovakia", "CS", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validation.Validate(tt.value, CountryCode2())
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

func TestCountryCode2_CommonMistakes(t *testing.T) {
	// Test common mistakes people make with country codes
	mistakes := []struct {
		name        string
		input       string
		description string
	}{
		{"UK instead of GB", "UK", "United Kingdom code is GB, not UK"},
		{"USA instead of US", "USA", "Should use 2-letter code US"},
		{"GER instead of DE", "GER", "Germany code is DE, not GER"},
		{"HOL instead of NL", "HOL", "Netherlands code is NL, not HOL"},
		{"ENG instead of GB", "ENG", "England is part of GB"},
		{"SCO instead of GB", "SCO", "Scotland is part of GB"},
		{"EUR instead of EU", "EUR", "EU is not a country code"},
		{"ASIA", "ASIA", "Continent name, not country code"},
		{"1A", "1A", "Starts with number"},
		{"A1", "A1", "Contains number"},
	}

	for _, m := range mistakes {
		t.Run(m.name, func(t *testing.T) {
			err := validation.Validate(m.input, CountryCode2())
			if err == nil {
				t.Errorf("%s: expected error but got nil", m.description)
			}
		})
	}
}

func TestCountryCode2_RegionalGroups(t *testing.T) {
	// Test country codes by regional groups
	regions := map[string][]string{
		"EU Members": {
			"AT", "BE", "BG", "HR", "CY", "CZ", "DK", "EE", "FI", "FR",
			"DE", "GR", "HU", "IE", "IT", "LV", "LT", "LU", "MT", "NL",
			"PL", "PT", "RO", "SK", "SI", "ES", "SE",
		},
		"G7": {
			"CA", "FR", "DE", "IT", "JP", "GB", "US",
		},
		"BRICS": {
			"BR", "RU", "IN", "CN", "ZA",
		},
		"ASEAN": {
			"BN", "KH", "ID", "LA", "MY", "MM", "PH", "SG", "TH", "VN",
		},
		"Nordic": {
			"DK", "FI", "IS", "NO", "SE",
		},
		"Caribbean": {
			"AG", "BS", "BB", "CU", "DM", "DO", "GD", "HT", "JM", "KN",
			"LC", "VC", "TT",
		},
		"Middle East": {
			"AE", "BH", "EG", "IR", "IQ", "IL", "JO", "KW", "LB", "OM",
			"QA", "SA", "SY", "TR", "YE",
		},
	}

	for region, codes := range regions {
		t.Run(region, func(t *testing.T) {
			for _, code := range codes {
				err := validation.Validate(code, CountryCode2())
				if err != nil {
					t.Errorf("%s: country code %s should be valid, got error: %v", region, code, err)
				}
			}
		})
	}
}

func TestCountryCode2_CaseSensitivity(t *testing.T) {
	// Test that validation is case-insensitive
	testCodes := []string{"US", "GB", "DE", "FR", "JP", "CN", "BR", "IN"}

	for _, code := range testCodes {
		variations := []string{
			code,                  // UPPERCASE
			strings.ToLower(code), // lowercase
			string(code[0]) + strings.ToLower(string(code[1])), // Mixedcase
			strings.ToLower(string(code[0])) + string(code[1]), // mIXEDCASE
		}

		for _, variant := range variations {
			t.Run(variant, func(t *testing.T) {
				err := validation.Validate(variant, CountryCode2())
				if err != nil {
					t.Errorf("expected %s to be valid (case-insensitive), got error: %v", variant, err)
				}
			})
		}
	}
}

func TestCountryCode2_SpecialCases(t *testing.T) {
	// Test special and disputed territories
	specialCases := []struct {
		name      string
		code      string
		expectErr bool
		note      string
	}{
		// Valid special territories
		{"Antarctica", "AQ", false, "International territory"},
		{"Bouvet Island", "BV", false, "Norwegian dependency"},
		{"Christmas Island", "CX", false, "Australian territory"},
		{"Heard Island", "HM", false, "Australian territory"},
		{"French Southern Territories", "TF", false, "French territory"},

		// User-assigned codes (should be valid in ISO but may vary)
		{"Kosovo", "XK", true, "User-assigned code not in ISO standard"},

		// Reserved codes (should be invalid)
		{"Reserved AA", "AA", true, "Reserved code"},
		{"Reserved QM-QZ", "QO", true, "Reserved range"},
		{"Reserved XA-XZ", "XZ", true, "Reserved range"},
		{"Reserved ZZ", "ZZ", true, "Reserved code"},
	}

	for _, sc := range specialCases {
		t.Run(sc.name, func(t *testing.T) {
			err := validation.Validate(sc.code, CountryCode2())
			if sc.expectErr && err == nil {
				t.Errorf("%s (%s): expected error but got nil - %s", sc.name, sc.code, sc.note)
			} else if !sc.expectErr && err != nil {
				t.Errorf("%s (%s): expected no error but got: %v - %s", sc.name, sc.code, err, sc.note)
			}
		})
	}
}

func TestCountryCode2Rule_When(t *testing.T) {
	err := CountryCode2().When(false).Validate("XX")
	if err != nil {
		t.Errorf("expected no error when condition is false, got: %v", err)
	}

	err = CountryCode2().When(true).Validate("XX")
	if err == nil {
		t.Error("expected error when condition is true and country code is invalid")
	}
}

func TestCountryCode2Rule_Error(t *testing.T) {
	customMsg := "custom country code error"
	err := CountryCode2().Error(customMsg).Validate("XX")
	if err == nil {
		t.Fatal("expected error but got nil")
	}
	if err.Error() != customMsg {
		t.Errorf("expected error message %q, got %q", customMsg, err.Error())
	}
}

func TestCountryCode2Rule_ErrorObject(t *testing.T) {
	customErr := validation.NewError("custom_code", "custom message")
	err := CountryCode2().ErrorObject(customErr).Validate("XX")
	if err == nil {
		t.Fatal("expected error but got nil")
	}
	if err.Error() != "custom message" {
		t.Errorf("expected custom error message, got: %v", err.Error())
	}
}

func BenchmarkCountryCode2Validation(b *testing.B) {
	codes := []any{
		"US",
		"GB",
		"de",
		"FR",
		"INVALID",
		"USA",
		"12",
		"",
		nil,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, code := range codes {
			_ = validation.Validate(code, CountryCode2())
		}
	}
}
