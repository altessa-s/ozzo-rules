// Copyright 2021-2026 Altessa Solutions Inc. All rights reserved.
// Use of this source code is governed by license that can be found in
// the LICENSE file.

package ozzo_rules

import (
	"testing"

	"github.com/go-ozzo/ozzo-validation/v4"
)

type ruInnEntry struct {
	value     any
	expectErr bool
	name      string
}

func TestRuINN(t *testing.T) {
	inns := []ruInnEntry{
		// Valid
		{value: "7736207543", expectErr: false, name: "Valid INN 1"},
		{value: int64(7736207543), expectErr: false, name: "Valid INN 2"},
		{value: "212501691178", expectErr: false, name: "Valid INN 3"},
		{value: ptrWrap("212501691178"), expectErr: false, name: "Valid INN 7"},
		{value: ptrWrap(int64(212501691178)), expectErr: false, name: "Valid INN 8"},

		// Invalid
		{value: int64(2222200000), expectErr: true, name: "Invalid INN 2"},
		{value: "999999999999999", expectErr: true, name: "Invalid INN 4"},
		{value: "111111", expectErr: true, name: "Invalid INN 5"},
		{value: "abc", expectErr: true, name: "Invalid INN 7"},
	}

	for _, tt := range inns {
		runRuInnTest(t, &tt, RuINN())
	}
}

func ptrWrap[T ~int64 | ~string](s T) *T {
	return &s
}

// Test RuINNLegal() validation method.
func TestRuINNLegal(t *testing.T) {
	inns := []ruInnEntry{
		// Valid
		{value: "7736207543", expectErr: false, name: "Valid INN 1"},
		{value: int64(7736207543), expectErr: false, name: "Valid INN 2"},

		// Invalid
		{value: "212501691178", expectErr: true, name: "Invalid INN 1"},
		{value: int64(2222200000), expectErr: true, name: "Invalid INN 2"},
		{value: "999999999999999", expectErr: true, name: "Invalid INN 4"},
		{value: "111111", expectErr: true, name: "Invalid INN 5"},
		{value: "abc", expectErr: true, name: "Invalid INN 7"},
	}

	for _, tt := range inns {
		runRuInnTest(t, &tt, RuINNLegal())
	}
}

// Test RuINNPersonal() validation method.
func TestRuINNPersonal(t *testing.T) {
	inns := []ruInnEntry{
		// Valid
		{value: "212501691178", expectErr: false, name: "Valid INN 1"},
		{value: int64(212501691178), expectErr: false, name: "Valid INN 2"},

		// Invalid
		{value: "1832021661o3", expectErr: true, name: "Invalid INN 1"},
		{value: "7736207543", expectErr: true, name: "Invalid INN 2"}, // legal valid
		{value: "999999999999999", expectErr: true, name: "Invalid INN 4"},
		{value: "111111", expectErr: true, name: "Invalid INN 5"},
		{value: int64(111111999999), expectErr: true, name: "Invalid INN 6"},
		{value: "abc", expectErr: true, name: "Invalid INN 7"},
	}

	for _, tt := range inns {
		runRuInnTest(t, &tt, RuINNPersonal())
	}
}

func runRuInnTest(t *testing.T, entry *ruInnEntry, rule RuINNRule) {
	err := validation.Validate(entry.value, rule)

	if entry.expectErr && err == nil {
		t.Errorf("Expected error, got nil")
	}

	if !entry.expectErr && err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}
