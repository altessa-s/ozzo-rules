// Copyright 2021-2026 Altessa Solutions Inc. All rights reserved.
// Use of this source code is governed by license that can be found in
// the LICENSE file.

package ozzo_rules

import (
	"testing"

	"github.com/go-ozzo/ozzo-validation/v4"
)

func TestRegexRule(t *testing.T) {
	tests := []struct {
		name      string
		value     any
		expectErr bool
	}{
		{"valid regex simple", "^test$", false},
		{"valid regex complex", `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`, false},
		{"valid regex groups", `(\d{3})-(\d{3})-(\d{4})`, false},
		{"valid regex unicode", `^[\p{L}\p{N}]+$`, false},
		{"invalid regex unclosed bracket", "[a-z", true},
		{"invalid regex bad parentheses", "((a)", true},
		{"invalid regex bad escape", `\k`, true},
		{"empty string", "", false},
		{"nil value", nil, true},
		{"non-string value", 123, true},
		{"pointer to valid regex", regexStrPtr("^test$"), false},
		{"pointer to invalid regex", regexStrPtr("[invalid"), true},
		{"pointer to nil", (*string)(nil), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Regex().Validate(tt.value)
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

func TestRegexRule_When(t *testing.T) {
	// Test with condition false - should skip validation
	err := Regex().When(false).Validate("[invalid")
	if err != nil {
		t.Errorf("expected no error when condition is false, got: %v", err)
	}

	// Test with condition true - should validate
	err = Regex().When(true).Validate("[invalid")
	if err == nil {
		t.Error("expected error when condition is true and regex is invalid")
	}
}

func TestRegexRule_Error(t *testing.T) {
	customMsg := "custom error message"
	err := Regex().Error(customMsg).Validate("[invalid")
	if err == nil {
		t.Fatal("expected error but got nil")
	}
	if err.Error() != customMsg {
		t.Errorf("expected error message %q, got %q", customMsg, err.Error())
	}
}

func TestRegexRule_ErrorObject(t *testing.T) {
	customErr := validation.NewError("custom_code", "custom message")
	err := Regex().ErrorObject(customErr).Validate("[invalid")
	if err == nil {
		t.Fatal("expected error but got nil")
	}
	// Check error message instead of comparing error objects
	if err.Error() != "custom message" {
		t.Errorf("expected custom error message, got: %v", err.Error())
	}
}

func TestRegexCache(t *testing.T) {
	// Clear cache before testing
	globalRegexCache.Purge()

	// Test that cache starts empty
	if size := globalRegexCache.Len(); size != 0 {
		t.Errorf("expected cache size 0, got %d", size)
	}

	// Compile a regex pattern
	pattern := "^test[0-9]+$"
	re1, err := compileRegex(pattern)
	if err != nil {
		t.Fatalf("failed to compile regex: %v", err)
	}

	// Check cache size increased
	if size := globalRegexCache.Len(); size != 1 {
		t.Errorf("expected cache size 1, got %d", size)
	}

	// Compile same pattern again - should return cached version
	re2, err := compileRegex(pattern)
	if err != nil {
		t.Fatalf("failed to get cached regex: %v", err)
	}

	// Verify same instance returned (pointer comparison)
	if re1 != re2 {
		t.Error("expected same regex instance from cache")
	}

	// Cache size should remain 1
	if size := globalRegexCache.Len(); size != 1 {
		t.Errorf("expected cache size still 1, got %d", size)
	}

	// Test invalid regex
	_, err = compileRegex("[invalid")
	if err == nil {
		t.Error("expected error for invalid regex")
	}

	// Cache size should still be 1 (invalid regex not cached)
	if size := globalRegexCache.Len(); size != 1 {
		t.Errorf("expected cache size still 1 after invalid regex, got %d", size)
	}
}

func TestRegexCacheConcurrency(t *testing.T) {
	globalRegexCache.Purge()

	// Test concurrent access to cache
	done := make(chan bool)
	pattern := "^concurrent[0-9]+$"

	// Launch multiple goroutines
	for i := 0; i < 10; i++ {
		go func() {
			for j := 0; j < 100; j++ {
				_, err := compileRegex(pattern)
				if err != nil {
					t.Errorf("concurrent compile failed: %v", err)
				}
			}
			done <- true
		}()
	}

	// Wait for all goroutines
	for i := 0; i < 10; i++ {
		<-done
	}

	// Should only have one entry in cache
	if size := globalRegexCache.Len(); size != 1 {
		t.Errorf("expected cache size 1 after concurrent access, got %d", size)
	}
}

func TestRegexLRUEviction(t *testing.T) {
	// Create a small cache for testing eviction
	smallCache := mustNewRegexCache(2)

	// Save original cache and restore after test
	originalCache := globalRegexCache
	globalRegexCache = smallCache
	defer func() {
		globalRegexCache = originalCache
	}()

	// Add patterns to fill cache
	patterns := []string{"^test1$", "^test2$", "^test3$"}

	for _, pattern := range patterns {
		_, err := compileRegex(pattern)
		if err != nil {
			t.Fatalf("failed to compile pattern %s: %v", pattern, err)
		}
	}

	// Cache should have only 2 items (LRU limit)
	if size := smallCache.Len(); size != 2 {
		t.Errorf("expected cache size 2, got %d", size)
	}

	// First pattern should be evicted
	if _, exists := smallCache.Get("^test1$"); exists {
		t.Error("expected first pattern to be evicted from LRU cache")
	}

	// Last two patterns should still be in cache
	if _, exists := smallCache.Get("^test2$"); !exists {
		t.Error("expected second pattern to be in cache")
	}
	if _, exists := smallCache.Get("^test3$"); !exists {
		t.Error("expected third pattern to be in cache")
	}
}

func BenchmarkRegexValidation(b *testing.B) {
	// Clear cache to ensure consistent benchmark
	globalRegexCache.Purge()

	patterns := []string{
		"^test$",
		`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`,
		`(\d{3})-(\d{3})-(\d{4})`,
		`^[\p{L}\p{N}]+$`,
	}

	b.Run("WithCache", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			for _, pattern := range patterns {
				_ = Regex().Validate(pattern)
			}
		}
	})

	b.Run("WithoutCache", func(b *testing.B) {
		// Clear cache before each iteration to simulate no caching
		for i := 0; i < b.N; i++ {
			globalRegexCache.Purge()
			for _, pattern := range patterns {
				_ = Regex().Validate(pattern)
			}
		}
	})
}

// Helper function for creating string pointers
func regexStrPtr(s string) *string {
	return &s
}
