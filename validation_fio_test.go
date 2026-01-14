// Copyright 2021-2026 Altessa Solutions Inc. All rights reserved.
// Use of this source code is governed by license that can be found in
// the LICENSE file.

package ozzo_rules

import (
	"testing"

	"github.com/go-ozzo/ozzo-validation/v4"
)

type fioEntry struct {
	value     any
	expectErr bool
	name      string
}

func TestRussianFirstName(t *testing.T) {
	tests := []fioEntry{
		{"Алексей", false, "Valid First Name 1"},
		{"Мария", false, "Valid First Name 2"},
		{"Иван", false, "Valid First Name 3"},
		{"Алексей-Сергей", false, "Valid First Name 4"},
		{"Юрий", false, "Valid First Name 5"},
		{"Эльвира", false, "Valid First Name 6"},
		{"Ольга", false, "Valid First Name 7"},
		{"Анна-Мария", false, "Valid First Name 8"},
		{fioStrPtr("Александр"), false, "Valid First Name Pointer 1"},
		{fioStrPtr("Николай"), false, "Valid First Name Pointer 2"},

		{"алексей", true, "Invalid First Name 1"},
		{"Alexey", true, "Invalid First Name 2"},
		{"А", true, "Invalid First Name 3"},
		{"Алексей!", true, "Invalid First Name 4"},
		{"", true, "Invalid First Name 5"},
		{fioStrPtr("12345"), true, "Invalid First Name Pointer 1"},
		{fioStrPtr("@#$%^"), true, "Invalid First Name Pointer 2"},
	}

	for _, tt := range tests {
		runFioTest(t, &tt, RussianFirstName())
	}
}

func TestRussianMiddleName(t *testing.T) {
	tests := []fioEntry{
		{"Александрович", false, "Valid Middle Name 1"},
		{"Ивановна", false, "Valid Middle Name 2"},
		{"Сергеевич", false, "Valid Middle Name 3"},
		{"Владимирович", false, "Valid Middle Name 4"},
		{fioStrPtr("Викторович"), false, "Valid Middle Name Pointer 1"},

		{"Ив", true, "Invalid Middle Name 1"},
		{"Иван", true, "Invalid Middle Name 2"},
		{"ivanovich", true, "Invalid Middle Name 3"},
		{"", true, "Invalid Middle Name 4"},
		{fioStrPtr("@@@@@@"), true, "Invalid Middle Name Pointer 1"},
	}

	for _, tt := range tests {
		runFioTest(t, &tt, RussianMiddleName())
	}
}

func TestRussianLastName(t *testing.T) {
	tests := []fioEntry{
		{"Иванов", false, "Valid Last Name 1"},
		{"Петров", false, "Valid Last Name 2"},
		{"Смирнов", false, "Valid Last Name 3"},
		{"Кузнецов-Петров", false, "Valid Last Name 4"},
		{"Сидоров", false, "Valid Last Name 5"},
		{"Александров", false, "Valid Last Name 6"},
		{"Капитонов", false, "Valid Last Name 7"},
		{fioStrPtr("Фёдоров"), false, "Valid Last Name Pointer 1"},

		{"иванов", true, "Invalid Last Name 1"},
		{"Smith", true, "Invalid Last Name 2"},
		{"И", true, "Invalid Last Name 3"},
		{"Иванов!", true, "Invalid Last Name 4"},
		{"", true, "Invalid Last Name 5"},
		{fioStrPtr("123456"), true, "Invalid Last Name Pointer 1"},
		{fioStrPtr("%%%%%%"), true, "Invalid Last Name Pointer 2"},
	}

	for _, tt := range tests {
		runFioTest(t, &tt, RussianLastName())
	}
}

func runFioTest(t *testing.T, entry *fioEntry, rule RussianNameRule) {
	err := validation.Validate(entry.value, rule)

	if entry.expectErr && err == nil {
		t.Errorf("%s: Expected error, got nil", entry.name)
	}

	if !entry.expectErr && err != nil {
		t.Errorf("%s: Expected no error, got %v", entry.name, err)
	}
}

func fioStrPtr(s string) *string {
	return &s
}
