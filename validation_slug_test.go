// Copyright 2021-2026 Altessa Solutions Inc. All rights reserved.
// Use of this source code is governed by license that can be found in
// the LICENSE file.

package ozzo_rules

import (
	"strings"
	"testing"

	"github.com/go-ozzo/ozzo-validation/v4"
)

func TestSlugRule(t *testing.T) {
	tests := []struct {
		name      string
		value     any
		expectErr bool
	}{
		// Valid slugs
		{"valid simple slug", "hello", false},
		{"valid with numbers", "hello123", false},
		{"valid with underscores", "hello_world", false},
		{"valid all lowercase", "abcdefghij", false},
		{"valid min length 5", "abcde", false},
		{"valid max length 100", strings.Repeat("a", 100), false},
		{"valid complex", "my_awesome_slug_123", false},
		{"valid all numbers after min", "12345", false},
		{"valid all underscores after min", "_____", false},
		{"valid snake_case", "snake_case_example", false},

		// Invalid slugs - length
		{"too short 4 chars", "abcd", true},
		{"too short 1 char", "a", true},
		{"empty string", "", true},
		{"too long 101 chars", strings.Repeat("a", 101), true},

		// Invalid slugs - uppercase
		{"with uppercase", "Hello", true},
		{"all uppercase", "HELLO", true},
		{"mixed case", "HelloWorld", true},
		{"uppercase in middle", "hello_World", true},

		// Invalid slugs - special characters
		{"with dash", "hello-world", true},
		{"with space", "hello world", true},
		{"with dot", "hello.world", true},
		{"with special char", "hello@world", true},
		{"with slash", "hello/world", true},
		{"with parentheses", "hello(world)", true},
		{"starts with number but short", "123", true},

		// Edge cases
		{"nil value", nil, true},
		{"non-string value", 12345, true},
		{"pointer to valid slug", slugStrPtr("hello_world"), false},
		{"pointer to invalid slug", slugStrPtr("Hello"), true},
		{"pointer to nil", (*string)(nil), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Slug().Validate(tt.value)
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

func TestSlugRule_When(t *testing.T) {
	err := Slug().When(false).Validate("INVALID")
	if err != nil {
		t.Errorf("expected no error when condition is false, got: %v", err)
	}

	err = Slug().When(true).Validate("INVALID")
	if err == nil {
		t.Error("expected error when condition is true and slug is invalid")
	}
}

func TestSlugRule_Error(t *testing.T) {
	customMsg := "custom slug error"
	err := Slug().Error(customMsg).Validate("BAD")
	if err == nil {
		t.Fatal("expected error but got nil")
	}
	if err.Error() != customMsg {
		t.Errorf("expected error message %q, got %q", customMsg, err.Error())
	}
}

func TestSlugRule_ErrorObject(t *testing.T) {
	customErr := validation.NewError("custom_code", "custom message")
	err := Slug().ErrorObject(customErr).Validate("BAD")
	if err == nil {
		t.Fatal("expected error but got nil")
	}
	if err.Error() != "custom message" {
		t.Errorf("expected custom error message, got: %v", err.Error())
	}
}

func TestSlugRule_CommonPatterns(t *testing.T) {
	// Test common slug patterns
	validSlugs := []string{
		// Database/API slugs
		"user_profile",
		"blog_post",
		"product_category",
		"order_status",
		"payment_method",

		// Feature flags
		"feature_enabled",
		"beta_access",
		"new_ui_2024",

		// Configuration keys
		"max_upload_size",
		"cache_timeout",
		"api_rate_limit",

		// Event names
		"user_created",
		"order_completed",
		"payment_failed",

		// System identifiers
		"system_admin",
		"guest_user",
		"default_role",
	}

	for _, slug := range validSlugs {
		t.Run("valid_slug_"+slug, func(t *testing.T) {
			err := Slug().Validate(slug)
			if err != nil {
				t.Errorf("expected valid slug %q, got error: %v", slug, err)
			}
		})
	}

	// Common invalid patterns
	invalidSlugs := []string{
		// Too short
		"user",
		"api",
		"db",

		// With uppercase
		"userProfile",
		"BlogPost",
		"CONSTANT_VALUE",

		// With dashes (kebab-case)
		"user-profile",
		"blog-post",
		"api-endpoint",

		// With spaces
		"user profile",
		"blog post",

		// With dots
		"user.profile",
		"api.v2",

		// Special characters
		"user@profile",
		"blog#post",
		"api/v2",

		// Too long
		strings.Repeat("a", 101),
	}

	for _, slug := range invalidSlugs {
		t.Run("invalid_slug_"+slug, func(t *testing.T) {
			err := Slug().Validate(slug)
			if err == nil {
				t.Errorf("expected invalid slug %q, but got no error", slug)
			}
		})
	}
}

func TestSlugRule_BoundaryTests(t *testing.T) {
	// Test exact boundaries
	tests := []struct {
		name   string
		length int
		valid  bool
	}{
		{"length 4", 4, false},
		{"length 5", 5, true},
		{"length 50", 50, true},
		{"length 99", 99, true},
		{"length 100", 100, true},
		{"length 101", 101, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			slug := strings.Repeat("a", tt.length)
			err := Slug().Validate(slug)
			if tt.valid && err != nil {
				t.Errorf("expected valid slug of length %d, got error: %v", tt.length, err)
			} else if !tt.valid && err == nil {
				t.Errorf("expected invalid slug of length %d, but got no error", tt.length)
			}
		})
	}
}

func BenchmarkSlugValidation(b *testing.B) {
	slugs := []string{
		"valid_slug",
		"user_profile_settings",
		"a",                      // too short
		"Invalid_Slug",           // uppercase
		"invalid-slug",           // dash
		strings.Repeat("a", 100), // max length
		strings.Repeat("a", 101), // too long
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, slug := range slugs {
			_ = Slug().Validate(slug)
		}
	}
}

// Helper function
func slugStrPtr(s string) *string {
	return &s
}
