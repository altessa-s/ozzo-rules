// Copyright 2021-2026 Altessa Solutions Inc. All rights reserved.
// Use of this source code is governed by license that can be found in
// the LICENSE file.

package ozzo_rules

import (
	"testing"

	"github.com/go-ozzo/ozzo-validation/v4"
)

func TestTimezoneRule(t *testing.T) {
	tests := []struct {
		name      string
		value     any
		expectErr bool
	}{
		// Valid timezones - UTC and offsets
		{"valid UTC", "UTC", false},
		{"valid GMT", "GMT", false},
		{"valid Local", "Local", false},

		// Valid timezones - Americas
		{"valid America/New_York", "America/New_York", false},
		{"valid America/Chicago", "America/Chicago", false},
		{"valid America/Los_Angeles", "America/Los_Angeles", false},
		{"valid America/Toronto", "America/Toronto", false},
		{"valid America/Mexico_City", "America/Mexico_City", false},
		{"valid America/Sao_Paulo", "America/Sao_Paulo", false},
		{"valid America/Buenos_Aires", "America/Buenos_Aires", false},
		{"valid US/Eastern", "US/Eastern", false},
		{"valid US/Central", "US/Central", false},
		{"valid US/Mountain", "US/Mountain", false},
		{"valid US/Pacific", "US/Pacific", false},
		{"valid Canada/Eastern", "Canada/Eastern", false},

		// Valid timezones - Europe
		{"valid Europe/London", "Europe/London", false},
		{"valid Europe/Paris", "Europe/Paris", false},
		{"valid Europe/Berlin", "Europe/Berlin", false},
		{"valid Europe/Madrid", "Europe/Madrid", false},
		{"valid Europe/Rome", "Europe/Rome", false},
		{"valid Europe/Moscow", "Europe/Moscow", false},
		{"valid Europe/Amsterdam", "Europe/Amsterdam", false},
		{"valid Europe/Stockholm", "Europe/Stockholm", false},
		{"valid Europe/Warsaw", "Europe/Warsaw", false},
		{"valid Europe/Athens", "Europe/Athens", false},

		// Valid timezones - Asia
		{"valid Asia/Tokyo", "Asia/Tokyo", false},
		{"valid Asia/Shanghai", "Asia/Shanghai", false},
		{"valid Asia/Hong_Kong", "Asia/Hong_Kong", false},
		{"valid Asia/Singapore", "Asia/Singapore", false},
		{"valid Asia/Seoul", "Asia/Seoul", false},
		{"valid Asia/Kolkata", "Asia/Kolkata", false},
		{"valid Asia/Dubai", "Asia/Dubai", false},
		{"valid Asia/Bangkok", "Asia/Bangkok", false},
		{"valid Asia/Jakarta", "Asia/Jakarta", false},
		{"valid Asia/Manila", "Asia/Manila", false},

		// Valid timezones - Australia/Pacific
		{"valid Australia/Sydney", "Australia/Sydney", false},
		{"valid Australia/Melbourne", "Australia/Melbourne", false},
		{"valid Australia/Brisbane", "Australia/Brisbane", false},
		{"valid Australia/Perth", "Australia/Perth", false},
		{"valid Pacific/Auckland", "Pacific/Auckland", false},
		{"valid Pacific/Fiji", "Pacific/Fiji", false},
		{"valid Pacific/Honolulu", "Pacific/Honolulu", false},

		// Valid timezones - Africa
		{"valid Africa/Cairo", "Africa/Cairo", false},
		{"valid Africa/Johannesburg", "Africa/Johannesburg", false},
		{"valid Africa/Lagos", "Africa/Lagos", false},
		{"valid Africa/Nairobi", "Africa/Nairobi", false},
		{"valid Africa/Casablanca", "Africa/Casablanca", false},

		// Valid timezones - Other
		{"valid Antarctica/McMurdo", "Antarctica/McMurdo", false},
		{"valid Indian/Mauritius", "Indian/Mauritius", false},
		{"valid Atlantic/Reykjavik", "Atlantic/Reykjavik", false},

		// Invalid timezones
		{"invalid America/InvalidCity", "America/InvalidCity", true},
		{"invalid timezone", "Invalid/Timezone", true},
		{"valid lowercase", "america/new_york", false}, // Go accepts lowercase
		{"invalid with space", "America/New York", true},
		{"invalid with dash", "America/New-York", true},
		{"invalid empty", "", true},
		{"invalid single word", "NewYork", true},
		{"invalid numeric", "GMT+5", true},
		{"invalid offset", "+05:00", true},
		{"invalid abbreviation", "EST", false}, // Actually supported by Go
		{"invalid PST", "PST", false},          // Actually supported by Go
		{"invalid CST", "CST", false},          // Actually supported by Go

		// Invalid types
		{"nil value", nil, true},
		{"integer value", 123, true},
		{"float value", 123.45, true},
		{"boolean value", true, true},
		{"slice value", []string{"UTC"}, true},

		// Edge cases
		{"pointer to valid", timezoneStrPtr("UTC"), false},
		{"pointer to invalid", timezoneStrPtr("Invalid"), true},
		{"pointer to nil", (*string)(nil), true},

		// Case sensitivity
		{"valid lowercase utc", "utc", false},              // Go accepts case variations
		{"valid uppercase EUROPE", "EUROPE/LONDON", false}, // Go accepts case variations
		{"valid mixed case", "Europe/london", false},       // Go accepts case variations
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Timezone().Validate(tt.value)
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

func TestTimezoneOrNilRule(t *testing.T) {
	tests := []struct {
		name      string
		value     any
		expectErr bool
	}{
		// Valid cases including nil
		{"valid timezone", "UTC", false},
		{"valid nil", nil, false},
		{"valid pointer nil", (*string)(nil), false},
		{"valid America/New_York", "America/New_York", false},

		// Invalid cases
		{"invalid timezone", "Invalid/Timezone", true},
		{"empty string", "", true},
		{"integer", 123, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := TimezoneOrNil().Validate(tt.value)
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

func TestTimezoneRule_When(t *testing.T) {
	err := Timezone().When(false).Validate("Invalid")
	if err != nil {
		t.Errorf("expected no error when condition is false, got: %v", err)
	}

	err = Timezone().When(true).Validate("Invalid")
	if err == nil {
		t.Error("expected error when condition is true and timezone is invalid")
	}

	// Test with TimezoneOrNil
	err = TimezoneOrNil().When(false).Validate("Invalid")
	if err != nil {
		t.Errorf("expected no error when condition is false, got: %v", err)
	}
}

func TestTimezoneRule_Error(t *testing.T) {
	customMsg := "custom timezone error"
	err := Timezone().Error(customMsg).Validate("Invalid")
	if err == nil {
		t.Fatal("expected error but got nil")
	}
	if err.Error() != customMsg {
		t.Errorf("expected error message %q, got %q", customMsg, err.Error())
	}

	// Test with TimezoneOrNil
	err = TimezoneOrNil().Error(customMsg).Validate("Invalid")
	if err == nil {
		t.Fatal("expected error but got nil")
	}
	if err.Error() != customMsg {
		t.Errorf("expected error message %q, got %q", customMsg, err.Error())
	}
}

func TestTimezoneRule_ErrorObject(t *testing.T) {
	customErr := validation.NewError("custom_code", "custom message")
	err := Timezone().ErrorObject(customErr).Validate("Invalid")
	if err == nil {
		t.Fatal("expected error but got nil")
	}
	if err.Error() != "custom message" {
		t.Errorf("expected custom error message, got: %v", err.Error())
	}

	// Test with TimezoneOrNil
	err = TimezoneOrNil().ErrorObject(customErr).Validate("Invalid")
	if err == nil {
		t.Fatal("expected error but got nil")
	}
	if err.Error() != "custom message" {
		t.Errorf("expected custom error message, got: %v", err.Error())
	}
}

func TestTimezoneRule_CommonScenarios(t *testing.T) {
	scenarios := []struct {
		name        string
		description string
		timezone    string
		expectErr   bool
	}{
		// Business scenarios
		{
			"user profile timezone",
			"storing user's preferred timezone",
			"America/New_York",
			false,
		},
		{
			"server default timezone",
			"setting server timezone",
			"UTC",
			false,
		},
		{
			"event timezone",
			"timezone for scheduled event",
			"Europe/London",
			false,
		},
		{
			"meeting across timezones",
			"international meeting timezone",
			"Asia/Tokyo",
			false,
		},

		// Common mistakes
		{
			"abbreviation instead of full",
			"using EST instead of America/New_York",
			"EST",
			false, // Go actually supports abbreviations
		},
		{
			"offset instead of name",
			"using +05:00 instead of timezone name",
			"+05:00",
			true,
		},
		{
			"city name only",
			"using just city name",
			"London",
			true,
		},
		{
			"country/city wrong format",
			"using country-city format",
			"US-NewYork",
			true,
		},

		// Legacy formats
		{
			"legacy US format",
			"old US timezone format",
			"US/Eastern",
			false,
		},
		{
			"legacy Canada format",
			"old Canada timezone format",
			"Canada/Pacific",
			false,
		},

		// Special cases
		{
			"local timezone",
			"system local timezone",
			"Local",
			false,
		},
		{
			"GMT timezone",
			"Greenwich Mean Time",
			"GMT",
			false,
		},
	}

	for _, sc := range scenarios {
		t.Run(sc.name, func(t *testing.T) {
			err := Timezone().Validate(sc.timezone)
			if sc.expectErr && err == nil {
				t.Errorf("%s: expected error but got nil", sc.description)
			} else if !sc.expectErr && err != nil {
				t.Errorf("%s: expected no error but got: %v", sc.description, err)
			}
		})
	}
}

func TestTimezoneRule_RegionalCoverage(t *testing.T) {
	// Test coverage of different regions
	regions := map[string][]string{
		"Americas": {
			"America/Anchorage",
			"America/Chicago",
			"America/Denver",
			"America/Halifax",
			"America/Los_Angeles",
			"America/New_York",
			"America/Phoenix",
			"America/Regina",
			"America/St_Johns",
			"America/Vancouver",
		},
		"Europe": {
			"Europe/Amsterdam",
			"Europe/Athens",
			"Europe/Berlin",
			"Europe/Brussels",
			"Europe/Bucharest",
			"Europe/Budapest",
			"Europe/Copenhagen",
			"Europe/Dublin",
			"Europe/Helsinki",
			"Europe/Istanbul",
			"Europe/Kiev",
			"Europe/Lisbon",
			"Europe/London",
			"Europe/Madrid",
			"Europe/Moscow",
			"Europe/Oslo",
			"Europe/Paris",
			"Europe/Prague",
			"Europe/Rome",
			"Europe/Stockholm",
			"Europe/Vienna",
			"Europe/Warsaw",
			"Europe/Zurich",
		},
		"Asia": {
			"Asia/Baghdad",
			"Asia/Bangkok",
			"Asia/Dhaka",
			"Asia/Dubai",
			"Asia/Hong_Kong",
			"Asia/Jakarta",
			"Asia/Jerusalem",
			"Asia/Karachi",
			"Asia/Kolkata",
			"Asia/Kuala_Lumpur",
			"Asia/Manila",
			"Asia/Seoul",
			"Asia/Shanghai",
			"Asia/Singapore",
			"Asia/Taipei",
			"Asia/Tehran",
			"Asia/Tokyo",
		},
		"Pacific": {
			"Pacific/Auckland",
			"Pacific/Fiji",
			"Pacific/Guam",
			"Pacific/Honolulu",
			"Pacific/Port_Moresby",
		},
		"Australia": {
			"Australia/Adelaide",
			"Australia/Brisbane",
			"Australia/Darwin",
			"Australia/Hobart",
			"Australia/Melbourne",
			"Australia/Perth",
			"Australia/Sydney",
		},
		"Africa": {
			"Africa/Cairo",
			"Africa/Casablanca",
			"Africa/Johannesburg",
			"Africa/Lagos",
			"Africa/Nairobi",
		},
	}

	for region, timezones := range regions {
		t.Run(region, func(t *testing.T) {
			for _, tz := range timezones {
				err := Timezone().Validate(tz)
				if err != nil {
					t.Errorf("expected valid timezone %s in region %s, got error: %v", tz, region, err)
				}
			}
		})
	}
}

func BenchmarkTimezoneValidation(b *testing.B) {
	timezones := []string{
		"UTC",
		"America/New_York",
		"Europe/London",
		"Asia/Tokyo",
		"Australia/Sydney",
		"Invalid/Timezone",
		"",
		"EST",
		"Local",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, tz := range timezones {
			_ = Timezone().Validate(tz)
		}
	}
}

func BenchmarkTimezoneOrNilValidation(b *testing.B) {
	values := []any{
		"UTC",
		"America/New_York",
		nil,
		(*string)(nil),
		"Invalid/Timezone",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, v := range values {
			_ = TimezoneOrNil().Validate(v)
		}
	}
}

// Helper function
func timezoneStrPtr(s string) *string {
	return &s
}
