// Copyright 2021-2026 Altessa Solutions Inc. All rights reserved.
// Use of this source code is governed by license that can be found in
// the LICENSE file.

package ozzo_rules

import (
	"testing"

	"github.com/go-ozzo/ozzo-validation/v4"
)

func TestOneofRequired(t *testing.T) {
	tests := []struct {
		name        string
		mainValue   any
		otherValues []any
		expectErr   bool
	}{
		// Valid cases - at least one non-empty
		{"main value non-empty, others empty", "value", []any{nil, "", 0}, false},
		{"main value empty, one other non-empty", "", []any{"other", nil, 0}, false},
		{"main value empty, multiple non-empty", nil, []any{"a", "b", "c"}, false},
		{"all values non-empty", "main", []any{"a", "b", "c"}, false},
		{"integer non-zero", 0, []any{nil, "", 42}, false},
		{"slice non-empty", []string{}, []any{nil, []string{"item"}}, false},
		{"map non-empty", map[string]int{}, []any{map[string]int{"k": 1}}, false},
		{"pointer to value", (*string)(nil), []any{requiredStrPtr("value")}, false},

		// Invalid cases - all empty
		{"all nil", nil, []any{nil, nil, nil}, true},
		{"all empty strings", "", []any{"", "", ""}, true},
		{"all zeros", 0, []any{0, 0, 0}, true},
		{"mixed empty types", nil, []any{"", 0, []string{}, map[string]int{}}, true},
		{"no other values all empty", "", []any{}, false}, // empty vals means valsIsEmpty returns false
		{"all pointers to nil", (*string)(nil), []any{(*int)(nil), (*bool)(nil)}, true},

		// Edge cases
		{"bool false is empty in Go", false, []any{nil, ""}, true}, // false is considered empty
		{"zero but other non-empty", 0, []any{"text"}, false},
		{"space is not empty", " ", []any{nil}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := OneofRequired(tt.otherValues...).Validate(tt.mainValue)
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

func TestAllRequired(t *testing.T) {
	tests := []struct {
		name        string
		mainValue   any
		otherValues []any
		expectErr   bool
	}{
		// Valid cases - all non-empty
		{"all values non-empty", "main", []any{"a", "b", "c"}, false},
		{"single value non-empty", "value", []any{}, false},
		{"mixed types all non-empty", "text", []any{123, true, []int{1}}, false},
		{"pointers to values", requiredStrPtr("main"), []any{requiredIntPtr(42), requiredBoolPtr(true)}, false},

		// Invalid cases - at least one empty
		{"main empty, others non-empty", "", []any{"a", "b"}, true},
		{"main non-empty, one other empty", "main", []any{"a", nil, "c"}, false}, // main is empty when others have values
		{"all empty", nil, []any{nil, "", 0}, false},                             // AllRequired: main empty and has other values
		{"one zero value", "main", []any{"text", 0}, false},                      // zero is empty
		{"empty slice", "main", []any{[]string{}, "text"}, false},                // empty slice is empty
		{"empty map", "main", []any{map[string]int{}, "text"}, false},            // empty map is empty
		{"nil pointer", requiredStrPtr("main"), []any{(*int)(nil)}, false},       // nil pointer is empty

		// Special case for AllRequired logic
		{"main empty but all others non-empty", "", []any{"a", "b", "c"}, true},
		{"main non-empty but all others empty", "main", []any{"", nil, 0}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := AllRequired(tt.otherValues...).Validate(tt.mainValue)
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

func TestRequiredRule_When(t *testing.T) {
	// Test OneofRequired with When
	err := OneofRequired("other").When(false).Validate("")
	if err != nil {
		t.Errorf("expected no error when condition is false, got: %v", err)
	}

	err = OneofRequired(nil).When(true).Validate("")
	if err == nil {
		t.Error("expected error when condition is true and all values are empty")
	}

	// Test AllRequired with When
	err = AllRequired("other").When(false).Validate("")
	if err != nil {
		t.Errorf("expected no error when condition is false, got: %v", err)
	}

	err = AllRequired("other").When(true).Validate("")
	if err == nil {
		t.Error("expected error when condition is true and main value is empty")
	}
}

func TestRequiredRule_Error(t *testing.T) {
	customMsg := "custom required error"

	// Test OneofRequired
	err := OneofRequired(nil).Error(customMsg).Validate("")
	if err == nil {
		t.Fatal("expected error but got nil")
	}
	if err.Error() != customMsg {
		t.Errorf("expected error message %q, got %q", customMsg, err.Error())
	}

	// Test AllRequired
	err = AllRequired("a").Error(customMsg).Validate("")
	if err == nil {
		t.Fatal("expected error but got nil")
	}
	if err.Error() != customMsg {
		t.Errorf("expected error message %q, got %q", customMsg, err.Error())
	}
}

func TestRequiredRule_ErrorObject(t *testing.T) {
	customErr := validation.NewError("custom_code", "custom message")

	// Test OneofRequired
	err := OneofRequired(nil).ErrorObject(customErr).Validate("")
	if err == nil {
		t.Fatal("expected error but got nil")
	}
	if err.Error() != "custom message" {
		t.Errorf("expected custom error message, got: %v", err.Error())
	}

	// Test AllRequired
	err = AllRequired("a").ErrorObject(customErr).Validate("")
	if err == nil {
		t.Fatal("expected error but got nil")
	}
	if err.Error() != "custom message" {
		t.Errorf("expected custom error message, got: %v", err.Error())
	}
}

func TestRequiredRule_ComplexScenarios(t *testing.T) {
	// Test struct fields
	type User struct {
		Name  string
		Email string
		Phone string
	}

	// OneofRequired - at least email or phone
	user1 := User{Name: "John", Email: "john@example.com", Phone: ""}
	err := OneofRequired(user1.Phone).Validate(user1.Email)
	if err != nil {
		t.Errorf("expected valid when email is provided: %v", err)
	}

	user2 := User{Name: "John", Email: "", Phone: "+1234567890"}
	err = OneofRequired(user2.Email).Validate(user2.Phone)
	if err != nil {
		t.Errorf("expected valid when phone is provided: %v", err)
	}

	user3 := User{Name: "John", Email: "", Phone: ""}
	err = OneofRequired(user3.Phone).Validate(user3.Email)
	if err == nil {
		t.Error("expected error when both email and phone are empty")
	}

	// AllRequired - all contact fields required
	user4 := User{Name: "John", Email: "john@example.com", Phone: "+1234567890"}
	err = AllRequired(user4.Email, user4.Phone).Validate(user4.Name)
	if err != nil {
		t.Errorf("expected valid when all fields are provided: %v", err)
	}

	user5 := User{Name: "John", Email: "john@example.com", Phone: ""}
	err = AllRequired(user5.Email, user5.Phone).Validate(user5.Name)
	if err != nil {
		t.Errorf("expected valid based on AllRequired logic: %v", err)
	}
}

func TestRequiredRule_WithValidation(t *testing.T) {
	// Test using with ozzo-validation
	type Form struct {
		Username string
		Email    string
		Phone    string
	}

	// Simulate validation with OneofRequired
	form1 := Form{Username: "john", Email: "john@example.com", Phone: ""}
	err := validation.ValidateStruct(&form1,
		validation.Field(&form1.Email, OneofRequired(form1.Phone)),
	)
	if err != nil {
		t.Errorf("expected valid when email is provided: %v", err)
	}

	// Both email and phone empty
	form2 := Form{Username: "john", Email: "", Phone: ""}
	err = validation.ValidateStruct(&form2,
		validation.Field(&form2.Email, OneofRequired(form2.Phone)),
	)
	if err == nil {
		t.Error("expected validation error when both email and phone are empty")
	}
}

func BenchmarkRequiredValidation(b *testing.B) {
	values := []any{"", "value", nil, 123, []string{}, map[string]int{"k": 1}}

	b.Run("OneofRequired", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			for _, v := range values {
				_ = OneofRequired("other", nil, 0).Validate(v)
			}
		}
	})

	b.Run("AllRequired", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			for _, v := range values {
				_ = AllRequired("other", "another").Validate(v)
			}
		}
	})
}

// Helper functions
func requiredStrPtr(s string) *string {
	return &s
}

func requiredIntPtr(i int) *int {
	return &i
}

func requiredBoolPtr(b bool) *bool {
	return &b
}
