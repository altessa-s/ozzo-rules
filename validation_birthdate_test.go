// Copyright 2021-2026 Altessa Solutions Inc. All rights reserved.
// Use of this source code is governed by license that can be found in
// the LICENSE file.

package ozzo_rules

import (
	"fmt"
	"testing"
	"time"

	"github.com/go-ozzo/ozzo-validation/v4"
)

func TestBirthdateRule(t *testing.T) {
	tests := []struct {
		name      string
		value     any
		expectErr bool
	}{
		// Valid birthdates - various years
		{"valid birthdate 2000", "2000-01-01", false},
		{"valid birthdate 1990", "1990-12-31", false},
		{"valid birthdate 1980", "1980-06-15", false},
		{"valid birthdate 1970", "1970-03-25", false},
		{"valid birthdate 1950", "1950-07-04", false},
		{"valid birthdate 1900", "1900-01-01", false},
		{"valid birthdate 2023", "2023-12-25", false},

		// Valid birthdates - various months
		{"valid January", "1990-01-15", false},
		{"valid February", "1990-02-28", false},
		{"valid March", "1990-03-31", false},
		{"valid April", "1990-04-30", false},
		{"valid May", "1990-05-01", false},
		{"valid June", "1990-06-30", false},
		{"valid July", "1990-07-31", false},
		{"valid August", "1990-08-31", false},
		{"valid September", "1990-09-30", false},
		{"valid October", "1990-10-31", false},
		{"valid November", "1990-11-30", false},
		{"valid December", "1990-12-31", false},

		// Valid birthdates - edge cases for days
		{"valid first day", "2000-01-01", false},
		{"valid last day 31", "2000-01-31", false},
		{"valid last day 30", "2000-04-30", false},
		{"valid leap year", "2000-02-29", false},
		{"valid non-leap year", "2001-02-28", false},

		// Valid with leading zeros
		{"valid with zeros", "2000-01-01", false},
		{"valid single digit day", "2000-01-09", false},
		{"valid single digit month", "2000-09-01", false},

		// Invalid formats
		{"invalid format slash", "2000/01/01", true},
		{"invalid format dot", "2000.01.01", true},
		{"invalid format US", "01-01-2000", true},
		{"invalid format EU", "01.01.2000", true},
		{"invalid no separators", "20000101", true},
		{"invalid with time", "2000-01-01 00:00:00", true},
		{"invalid with T", "2000-01-01T00:00:00", true},
		{"invalid ISO8601", "2000-01-01T00:00:00Z", true},

		// Invalid dates - month out of range
		{"invalid month 00", "2000-00-01", false}, // Regex allows 00-19
		{"invalid month 13", "2000-13-01", false}, // Regex allows 00-19
		{"invalid month 20", "2000-20-01", true},  // Regex doesn't allow 20+
		{"invalid month 99", "2000-99-01", true},

		// Invalid dates - day out of range
		{"invalid day 00", "2000-01-00", false}, // Regex allows 00-39
		{"invalid day 32", "2000-01-32", false}, // Regex allows 00-39
		{"invalid day 40", "2000-01-40", true},  // Regex doesn't allow 40+
		{"invalid day 99", "2000-01-99", true},

		// Invalid dates - regex validation (doesn't check actual date validity)
		{"invalid Feb 31", "2000-02-31", false},      // Regex allows it
		{"invalid Apr 31", "2000-04-31", false},      // Regex allows it
		{"invalid leap Feb 30", "2000-02-30", false}, // Regex allows it

		// Invalid year formats
		{"invalid year 2 digits", "00-01-01", true},
		{"invalid year 3 digits", "200-01-01", true},
		{"invalid year 5 digits", "20000-01-01", true},
		{"invalid year letters", "abcd-01-01", true},

		// Invalid separators
		{"invalid no dash 1", "2000 01 01", true},
		{"invalid no dash 2", "200001-01", true},
		{"invalid no dash 3", "2000-0101", true},
		{"invalid extra dash", "2000-01-01-", true},
		{"invalid double dash", "2000--01-01", true},

		// Invalid types
		{"empty string", "", true},
		{"nil value", nil, true},
		{"integer value", 20000101, true},
		{"float value", 2000.0101, true},
		{"boolean value", true, true},
		{"slice value", []string{"2000-01-01"}, true},

		// Edge cases
		{"pointer to valid", birthdateStrPtr("2000-01-01"), false},
		{"pointer to invalid", birthdateStrPtr("2000/01/01"), true},
		{"pointer to nil", (*string)(nil), true},

		// Special strings
		{"text instead of date", "birthday", true},
		{"partial date", "2000-01", true},
		{"just year", "2000", true},
		{"random numbers", "1234-56-78", true},
		{"with spaces", " 2000-01-01 ", true},
		{"with prefix", "Date: 2000-01-01", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Birthdate().Validate(tt.value)
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

func TestBirthdateRule_When(t *testing.T) {
	err := Birthdate().When(false).Validate("invalid-date")
	if err != nil {
		t.Errorf("expected no error when condition is false, got: %v", err)
	}

	err = Birthdate().When(true).Validate("invalid-date")
	if err == nil {
		t.Error("expected error when condition is true and date is invalid")
	}
}

func TestBirthdateRule_Error(t *testing.T) {
	customMsg := "custom birthdate error"
	err := Birthdate().Error(customMsg).Validate("invalid")
	if err == nil {
		t.Fatal("expected error but got nil")
	}
	if err.Error() != customMsg {
		t.Errorf("expected error message %q, got %q", customMsg, err.Error())
	}
}

func TestBirthdateRule_ErrorObject(t *testing.T) {
	customErr := validation.NewError("custom_code", "custom message")
	err := Birthdate().ErrorObject(customErr).Validate("invalid")
	if err == nil {
		t.Fatal("expected error but got nil")
	}
	if err.Error() != "custom message" {
		t.Errorf("expected custom error message, got: %v", err.Error())
	}
}

func TestBirthdateRule_CommonScenarios(t *testing.T) {
	// Current year for age calculations
	currentYear := time.Now().Year()

	scenarios := []struct {
		name        string
		description string
		date        string
		expectErr   bool
	}{
		// Age ranges
		{
			"adult age 25",
			"25 year old",
			fmt.Sprintf("%d-06-15", currentYear-25),
			false,
		},
		{
			"senior age 65",
			"65 year old",
			fmt.Sprintf("%d-01-01", currentYear-65),
			false,
		},
		{
			"child age 10",
			"10 year old",
			fmt.Sprintf("%d-12-25", currentYear-10),
			false,
		},
		{
			"newborn this year",
			"born this year",
			fmt.Sprintf("%d-03-15", currentYear),
			false,
		},

		// Historical dates
		{
			"millennium baby",
			"born in 2000",
			"2000-01-01",
			false,
		},
		{
			"90s kid",
			"born in 1995",
			"1995-07-20",
			false,
		},
		{
			"baby boomer",
			"born in 1960",
			"1960-05-10",
			false,
		},
		{
			"very old person",
			"born in 1920",
			"1920-11-30",
			false,
		},

		// Special dates
		{
			"leap day baby",
			"born on leap day",
			"2000-02-29",
			false,
		},
		{
			"new year baby",
			"born on Jan 1",
			"2000-01-01",
			false,
		},
		{
			"christmas baby",
			"born on Christmas",
			"1990-12-25",
			false,
		},
		{
			"valentine baby",
			"born on Valentine's",
			"1985-02-14",
			false,
		},

		// Invalid common mistakes
		{
			"US format",
			"MM/DD/YYYY format",
			"12/25/2000",
			true,
		},
		{
			"EU format",
			"DD.MM.YYYY format",
			"25.12.2000",
			true,
		},
		{
			"wrong separator",
			"using dots",
			"2000.12.25",
			true,
		},
		{
			"with timestamp",
			"includes time",
			"2000-12-25 00:00:00",
			true,
		},
	}

	for _, sc := range scenarios {
		t.Run(sc.name, func(t *testing.T) {
			err := Birthdate().Validate(sc.date)
			if sc.expectErr && err == nil {
				t.Errorf("%s: expected error but got nil", sc.description)
			} else if !sc.expectErr && err != nil {
				t.Errorf("%s: expected no error but got: %v", sc.description, err)
			}
		})
	}
}

func TestBirthdateRule_EdgeDates(t *testing.T) {
	// Test month boundaries
	monthDays := map[string]int{
		"01": 31, "02": 29, "03": 31, "04": 30,
		"05": 31, "06": 30, "07": 31, "08": 31,
		"09": 30, "10": 31, "11": 30, "12": 31,
	}

	for month, maxDay := range monthDays {
		// Test valid last day
		date := fmt.Sprintf("2000-%s-%02d", month, maxDay)
		t.Run(fmt.Sprintf("valid_month_%s_day_%d", month, maxDay), func(t *testing.T) {
			err := Birthdate().Validate(date)
			if err != nil {
				t.Errorf("expected valid date %s, got error: %v", date, err)
			}
		})

		// Test invalid day (note: regex doesn't validate actual calendar)
		if maxDay < 31 {
			date = fmt.Sprintf("2000-%s-31", month)
			t.Run(fmt.Sprintf("month_%s_day_31", month), func(t *testing.T) {
				err := Birthdate().Validate(date)
				// The regex allows day 31 for any month
				if err != nil {
					t.Errorf("regex should allow %s even if calendar-invalid", date)
				}
			})
		}
	}
}

func TestBirthdateRule_YearBoundaries(t *testing.T) {
	tests := []struct {
		year      string
		expectErr bool
	}{
		// Valid 4-digit years
		{"0001", false},
		{"1000", false},
		{"1900", false},
		{"2000", false},
		{"2024", false},
		{"9999", false},

		// Invalid year formats
		{"000", true},   // 3 digits
		{"99", true},    // 2 digits
		{"1", true},     // 1 digit
		{"10000", true}, // 5 digits
		{"", true},      // empty
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("year_%s", tt.year), func(t *testing.T) {
			date := fmt.Sprintf("%s-01-01", tt.year)
			if tt.year == "" {
				date = "-01-01"
			}
			err := Birthdate().Validate(date)
			if tt.expectErr && err == nil {
				t.Errorf("expected error for year %s but got nil", tt.year)
			} else if !tt.expectErr && err != nil {
				t.Errorf("expected no error for year %s but got: %v", tt.year, err)
			}
		})
	}
}

func TestBirthdateRule_RegexPattern(t *testing.T) {
	// Test the regex pattern boundaries
	tests := []struct {
		name      string
		date      string
		expectErr bool
	}{
		// Month boundaries ([01]\d allows 00-19)
		{"month 00", "2000-00-01", false}, // Allowed by regex
		{"month 01", "2000-01-01", false},
		{"month 09", "2000-09-01", false},
		{"month 10", "2000-10-01", false},
		{"month 12", "2000-12-01", false},
		{"month 13", "2000-13-01", false}, // Allowed by regex
		{"month 19", "2000-19-01", false}, // Allowed by regex
		{"month 20", "2000-20-01", true},

		// Day boundaries ([0-3]\d allows 00-39)
		{"day 00", "2000-01-00", false}, // Allowed by regex
		{"day 01", "2000-01-01", false},
		{"day 09", "2000-01-09", false},
		{"day 10", "2000-01-10", false},
		{"day 29", "2000-01-29", false},
		{"day 30", "2000-01-30", false},
		{"day 31", "2000-01-31", false},
		{"day 39", "2000-01-39", false}, // Regex allows up to 39
		{"day 40", "2000-01-40", true},
		{"day 99", "2000-01-99", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Birthdate().Validate(tt.date)
			if tt.expectErr && err == nil {
				t.Errorf("expected error for %s but got nil", tt.date)
			} else if !tt.expectErr && err != nil {
				t.Errorf("expected no error for %s but got: %v", tt.date, err)
			}
		})
	}
}

func BenchmarkBirthdateValidation(b *testing.B) {
	dates := []string{
		"2000-01-01",
		"1990-12-31",
		"2024-06-15",
		"invalid-date",
		"2000/01/01",
		"01-01-2000",
		"2000-13-01",
		"2000-01-32",
		"",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, date := range dates {
			_ = Birthdate().Validate(date)
		}
	}
}

// Helper function
func birthdateStrPtr(s string) *string {
	return &s
}
