// Copyright 2021-2026 Altessa Solutions Inc. All rights reserved.
// Use of this source code is governed by license that can be found in
// the LICENSE file.

package ozzo_rules

import (
	"testing"

	"github.com/go-ozzo/ozzo-validation/v4"
)

func TestPasswordRule(t *testing.T) {
	tests := []struct {
		name      string
		value     any
		expectErr bool
	}{
		// Valid passwords
		{"valid password basic", "Pass123!", false},
		{"valid password complex", "MyP@ssw0rd!", false},
		{"valid password with space", "Pass 123!", false},
		{"valid password with symbols", "Pass#123$", false},
		{"valid password with underscore", "Pass_123!", false},
		{"valid password long", "VeryLongP@ssw0rd123456789", false},
		{"valid password exact 8 chars", "Pass123!", false},

		// Invalid passwords - too short
		{"too short 7 chars", "Pass12!", true},
		{"too short 1 char", "P", true},
		{"empty string", "", true},

		// Invalid passwords - missing requirements
		{"missing uppercase", "pass123!", true},
		{"missing number", "Password!", true},
		{"missing punctuation", "Password123", true},
		{"only letters", "Password", true},
		{"only numbers", "12345678", true},
		{"only symbols", "!@#$%^&*", true},

		// Invalid passwords - non-ASCII
		{"cyrillic characters", "Пароль123!", true},
		{"chinese characters", "密码Pass123!", true},
		{"emoji", "Pass123!😊", true},
		{"unicode symbols", "Pass123!™", true},

		// Edge cases
		{"nil value", nil, true},
		{"non-string value", 12345678, true},
		{"pointer to valid password", passwordStrPtr("Pass123!"), false},
		{"pointer to invalid password", passwordStrPtr("pass"), true},
		{"pointer to nil", (*string)(nil), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Password().Validate(tt.value)
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

func TestPasswordRule_When(t *testing.T) {
	// Test with condition false - should skip validation
	err := Password().When(false).Validate("invalid")
	if err != nil {
		t.Errorf("expected no error when condition is false, got: %v", err)
	}

	// Test with condition true - should validate
	err = Password().When(true).Validate("invalid")
	if err == nil {
		t.Error("expected error when condition is true and password is invalid")
	}
}

func TestPasswordRule_Error(t *testing.T) {
	customMsg := "custom password error"
	err := Password().Error(customMsg).Validate("invalid")
	if err == nil {
		t.Fatal("expected error but got nil")
	}
	if err.Error() != customMsg {
		t.Errorf("expected error message %q, got %q", customMsg, err.Error())
	}
}

func TestPasswordRule_ErrorObject(t *testing.T) {
	customErr := validation.NewError("custom_code", "custom message")
	err := Password().ErrorObject(customErr).Validate("invalid")
	if err == nil {
		t.Fatal("expected error but got nil")
	}
	if err.Error() != "custom message" {
		t.Errorf("expected custom error message, got: %v", err.Error())
	}
}

func TestPasswordRule_SpecialCharacters(t *testing.T) {
	// Test various special characters that should be accepted
	specialChars := []string{
		"Pass123!",  // exclamation
		"Pass123@",  // at
		"Pass123#",  // hash
		"Pass123$",  // dollar
		"Pass123%",  // percent
		"Pass123^",  // caret
		"Pass123&",  // ampersand
		"Pass123*",  // asterisk
		"Pass123(",  // parenthesis
		"Pass123)",  // parenthesis
		"Pass123-",  // hyphen
		"Pass123_",  // underscore
		"Pass123+",  // plus
		"Pass123=",  // equals
		"Pass123[",  // bracket
		"Pass123]",  // bracket
		"Pass123{",  // brace
		"Pass123}",  // brace
		"Pass123|",  // pipe
		"Pass123\\", // backslash
		"Pass123:",  // colon
		"Pass123;",  // semicolon
		"Pass123\"", // quote
		"Pass123'",  // apostrophe
		"Pass123<",  // less than
		"Pass123>",  // greater than
		"Pass123,",  // comma
		"Pass123.",  // period
		"Pass123?",  // question
		"Pass123/",  // slash
		"Pass123~",  // tilde
		"Pass123`",  // backtick
		"Pass 123",  // space (also counts as punctuation)
	}

	for _, password := range specialChars {
		t.Run("special char: "+password, func(t *testing.T) {
			err := Password().Validate(password)
			if err != nil {
				t.Errorf("expected valid password with special char, got error: %v", err)
			}
		})
	}
}

func BenchmarkPasswordValidation(b *testing.B) {
	passwords := []string{
		"Pass123!",
		"invalid",
		"VeryLongP@ssw0rd123456789",
		"Пароль123!",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, password := range passwords {
			_ = Password().Validate(password)
		}
	}
}

// Helper function for creating string pointers
func passwordStrPtr(s string) *string {
	return &s
}
