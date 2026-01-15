// Copyright 2021-2026 Altessa Solutions Inc. All rights reserved.
// Use of this source code is governed by license that can be found in
// the LICENSE file.

package ozzo_rules

import (
	"testing"

	"github.com/go-ozzo/ozzo-validation/v4"
)

type ruOgrnEntry struct {
	value     any
	expectErr bool
	name      string
}

func TestRuOGRN(t *testing.T) {
	validStrOGRN1 := "1026605606620"
	validIntOGRN2 := int64(1026605606620)

	invalidStrOGRN3 := "1037739877295999"
	invalidIntOGRN4 := int64(1037739877295999)
	ogrns := []ruOgrnEntry{
		// UL
		{&validStrOGRN1, false, "Valid OGRN 1"},
		{&validIntOGRN2, false, "Valid OGRN 2"},
		{"1026605606620", false, "Valid OGRN 3"},
		{"1037739877295", false, "Valid OGRN 4"},
		// IP
		{"304770001231454", false, "Valid OGRN 5"},
		{"304667010400037", false, "Valid OGRN 6"},
		// Invalid
		{"aa47700012314bb", true, "Invalid OGRN 1"},
		{"1026605606629", true, "Invalid OGRN 2"},
		{"304770001231450", true, "Invalid OGRN 4"},
		{"304770001231454444", true, "Invalid OGRN 5"},
		{&invalidStrOGRN3, true, "Invalid OGRN 6"},
		{&invalidIntOGRN4, true, "Invalid OGRN 7"},
	}

	for _, tt := range ogrns {
		runRuOgrnsTest(t, &tt, RuOGRN())
	}
}

// Test RuOGRNLegal() validation method.
func TestRuOGRNLegal(t *testing.T) {
	ogrns := []ruOgrnEntry{
		// Valid
		{"1026605606620", false, "Valid OGRN 1"},
		{"1037739877295", false, "Valid OGRN 2"},

		// Invalid
		{value: "2125016911780", expectErr: true, name: "Invalid OGRN 1"},
		{value: int64(2222200000), expectErr: true, name: "Invalid OGRN 2"},
		{value: "999999999999999", expectErr: true, name: "Invalid OGRN 3"},
		{value: "111111", expectErr: true, name: "Invalid OGRN 4"},
		{value: "abc", expectErr: true, name: "Invalid OGRN 5"},
		{"304770001231454", true, "Valid OGRN 6"},
		{"304667010400037", true, "Valid OGRN 7"},
	}

	for _, tt := range ogrns {
		runRuOgrnsTest(t, &tt, RuOGRNLegal())
	}
}

// Test RuOGRNEntrepreneur() validation method.
func TestRuOGRNEntrepreneur(t *testing.T) {
	ogrns := []ruOgrnEntry{
		// Valid
		{"304770001231454", false, "Valid OGRN 1"},
		{"304667010400037", false, "Valid OGRN 2"},
		{"305890118100052", false, "Valid OGRN 3"},

		// Invalid
		{value: "1832021661o3", expectErr: true, name: "Invalid OGRN 1"},
		{value: "7736207543", expectErr: true, name: "Invalid OGRN 2"},
		{value: "999999999999999", expectErr: true, name: "Invalid OGRN 4"},
		{value: "111111", expectErr: true, name: "Invalid OGRN 5"},
		{value: int64(111111999999), expectErr: true, name: "Invalid OGRN 6"},
		{value: "abc", expectErr: true, name: "Invalid OGRN 7"},
		{"1026605606620", true, "Valid OGRN 8"},
		{"1037739877295", true, "Valid OGRN 9"},
	}

	for _, tt := range ogrns {
		runRuOgrnsTest(t, &tt, RuOGRNEntrepreneur())
	}
}

func runRuOgrnsTest(t *testing.T, entry *ruOgrnEntry, rule RuOGRNRule) {
	err := validation.Validate(entry.value, rule)

	if entry.expectErr && err == nil {
		t.Errorf("Expected error, got nil")
	}

	if !entry.expectErr && err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}
