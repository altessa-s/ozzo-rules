// Copyright 2021-2026 Altessa Solutions Inc. All rights reserved.
// Use of this source code is governed by license that can be found in
// the LICENSE file.

package ozzo_rules

import (
	"strings"
	"testing"

	"github.com/go-ozzo/ozzo-validation/v4"
)

func TestUsernameRule(t *testing.T) {
	tests := []struct {
		name      string
		value     any
		expectErr bool
	}{
		// Valid usernames
		{"valid simple username", "john", false},
		{"valid with numbers", "john123", false},
		{"valid with underscore", "john_doe", false},
		{"valid with dash", "john-doe", false},
		{"valid mixed", "John_Doe-123", false},
		{"valid single char", "j", false},
		{"valid max length", strings.Repeat("a", 56), false},
		{"valid all numbers", "123456", false},
		{"valid all underscores", "___", false},
		{"valid all dashes", "---", false},
		{"valid complex", "User_123-test_ABC", false},

		// Invalid usernames
		{"empty string", "", true},
		{"too long", strings.Repeat("a", 57), true},
		{"with space", "john doe", true},
		{"with special char @", "john@doe", true},
		{"with special char .", "john.doe", true},
		{"with special char !", "john!", true},
		{"with parentheses", "john(doe)", true},
		{"with slash", "john/doe", true},
		{"with plus", "john+doe", true},
		{"with equals", "john=doe", true},
		{"starts with space", " john", true},
		{"ends with space", "john ", true},
		{"non-latin letters", "joão", true},
		{"cyrillic", "иван", true},
		{"chinese", "用户", true},
		{"emoji", "john😊", true},

		// Edge cases
		{"nil value", nil, true},
		{"non-string value", 123, true},
		{"pointer to valid username", usernameStrPtr("john_doe"), false},
		{"pointer to invalid username", usernameStrPtr("john@doe"), true},
		{"pointer to nil", (*string)(nil), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Username().Validate(tt.value)
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

func TestUsernameRule_When(t *testing.T) {
	err := Username().When(false).Validate("invalid@username")
	if err != nil {
		t.Errorf("expected no error when condition is false, got: %v", err)
	}

	err = Username().When(true).Validate("invalid@username")
	if err == nil {
		t.Error("expected error when condition is true and username is invalid")
	}
}

func TestUsernameRule_Error(t *testing.T) {
	customMsg := "custom username error"
	err := Username().Error(customMsg).Validate("invalid@")
	if err == nil {
		t.Fatal("expected error but got nil")
	}
	if err.Error() != customMsg {
		t.Errorf("expected error message %q, got %q", customMsg, err.Error())
	}
}

func TestUsernameRule_ErrorObject(t *testing.T) {
	customErr := validation.NewError("custom_code", "custom message")
	err := Username().ErrorObject(customErr).Validate("invalid@")
	if err == nil {
		t.Fatal("expected error but got nil")
	}
	if err.Error() != "custom message" {
		t.Errorf("expected custom error message, got: %v", err.Error())
	}
}

func TestUsernameRule_EdgeCases(t *testing.T) {
	// Test specific edge cases for the regex
	edgeCases := []struct {
		name      string
		username  string
		expectErr bool
	}{
		// Mixed case (regex accepts both upper and lower)
		{"lowercase only", "johndoe", false},
		{"uppercase only", "JOHNDOE", false},
		{"mixed case", "JohnDoe", false},
		{"camelCase", "johnDoe", false},
		{"PascalCase", "JohnDoe", false},

		// Number positions
		{"starts with number", "123john", false},
		{"ends with number", "john123", false},
		{"only numbers", "12345", false},

		// Special character positions
		{"starts with underscore", "_john", false},
		{"ends with underscore", "john_", false},
		{"starts with dash", "-john", false},
		{"ends with dash", "john-", false},
		{"multiple underscores", "john__doe", false},
		{"multiple dashes", "john--doe", false},
		{"underscore and dash", "john_doe-123", false},

		// Length boundaries
		{"exactly 56 chars", strings.Repeat("a", 56), false},
		{"exactly 57 chars", strings.Repeat("a", 57), true},
		{"single char letter", "a", false},
		{"single char number", "1", false},
		{"single char underscore", "_", false},
		{"single char dash", "-", false},
	}

	for _, tc := range edgeCases {
		t.Run(tc.name, func(t *testing.T) {
			err := Username().Validate(tc.username)
			if tc.expectErr && err == nil {
				t.Error("expected error but got nil")
			} else if !tc.expectErr && err != nil {
				t.Errorf("expected no error but got: %v", err)
			}
		})
	}
}

func TestUsernameRule_CommonPatterns(t *testing.T) {
	// Test common username patterns
	validPatterns := []string{
		// Common formats
		"user123",
		"first_last",
		"first-last",
		"user_2024",
		"admin",
		"root",
		"guest",
		"test_user",
		"dev-user",
		"api_key_123",

		// GitHub/GitLab style
		"john-doe",
		"jane_doe",
		"user123",
		"test-repo-owner",

		// System usernames
		"www-data",
		"mysql",
		"postgres",
		"nginx",
		"apache2",

		// Service accounts
		"svc_account",
		"bot_user",
		"ci_runner",
		"deploy_bot",
	}

	for _, pattern := range validPatterns {
		t.Run("valid_pattern_"+pattern, func(t *testing.T) {
			err := Username().Validate(pattern)
			if err != nil {
				t.Errorf("expected valid username pattern %q, got error: %v", pattern, err)
			}
		})
	}

	// Common invalid patterns
	invalidPatterns := []string{
		// Email-like
		"user@domain",
		"user@example.com",

		// Domain-like
		"user.name",
		"first.last",

		// With spaces
		"user name",
		"first last",

		// Special characters
		"user$name",
		"user%20",
		"user+plus",
		"user*star",
		"user(parens)",
		"user[brackets]",
		"user{braces}",

		// Path-like
		"user/name",
		"path\\user",

		// Too long
		"this_username_is_way_too_long_and_exceeds_the_maximum_allowed_length_of_56",
	}

	for _, pattern := range invalidPatterns {
		t.Run("invalid_pattern_"+pattern, func(t *testing.T) {
			err := Username().Validate(pattern)
			if err == nil {
				t.Errorf("expected invalid username pattern %q, but got no error", pattern)
			}
		})
	}
}

func BenchmarkUsernameValidation(b *testing.B) {
	usernames := []string{
		"validuser",
		"user_123",
		"john-doe",
		"invalid@user",
		"user with spaces",
		strings.Repeat("a", 56),
		strings.Repeat("a", 57),
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, username := range usernames {
			_ = Username().Validate(username)
		}
	}
}

// Helper function
func usernameStrPtr(s string) *string {
	return &s
}
