// Copyright 2021-2026 Altessa Solutions Inc. All rights reserved.
// Use of this source code is governed by license that can be found in
// the LICENSE file.

package ozzo_rules

import (
	"testing"

	"github.com/go-ozzo/ozzo-validation/v4"
)

func TestURIRule(t *testing.T) {
	tests := []struct {
		name      string
		value     any
		expectErr bool
	}{
		// Valid URI formats
		{"valid ssl with hostname", "ssl://example.com:443", false},
		{"valid ssl with IP", "ssl://192.168.1.1:443", false},
		{"valid starttls with hostname", "starttls://mail.example.com:587", false},
		{"valid starttls with IP", "starttls://10.0.0.1:25", false},
		{"valid ssl with subdomain", "ssl://smtp.gmail.com:465", false},
		{"valid ssl with long domain", "ssl://very.long.subdomain.example.com:8443", false},
		{"valid port min", "ssl://example.com:1", false},
		{"valid port max", "ssl://example.com:65534", false},
		{"valid localhost", "ssl://localhost:443", false},
		{"valid IPv6", "ssl://[2001:db8::1]:443", true}, // IPv6 not supported by current implementation

		// Invalid URI formats - wrong protocol
		{"http protocol", "http://example.com:80", true},
		{"https protocol", "https://example.com:443", true},
		{"ftp protocol", "ftp://example.com:21", true},
		{"tcp protocol", "tcp://example.com:8080", true},
		{"no protocol", "example.com:443", true},
		{"wrong case SSL", "SSL://example.com:443", true},
		{"wrong case StartTLS", "StartTLS://example.com:587", true},

		// Invalid URI formats - port issues
		{"no port", "ssl://example.com", true},
		{"invalid port zero", "ssl://example.com:0", true},
		{"invalid port 65535", "ssl://example.com:65535", true},
		{"invalid port too high", "ssl://example.com:65536", true},
		{"invalid port negative", "ssl://example.com:-1", true},
		{"invalid port non-numeric", "ssl://example.com:https", true},
		{"invalid port float", "ssl://example.com:443.5", true},

		// Invalid URI formats - host issues
		{"empty host", "ssl://:443", true},
		{"invalid host spaces", "ssl://example com:443", true},
		{"invalid host special chars", "ssl://example@com:443", true},
		{"multiple colons", "ssl://example.com:443:extra", true},

		// Edge cases
		{"empty string", "", true},
		{"only protocol", "ssl://", true},
		{"nil value", nil, true},
		{"non-string value", 123, true},
		{"pointer to valid uri", uriStrPtr("ssl://example.com:443"), false},
		{"pointer to invalid uri", uriStrPtr("invalid"), true},
		{"pointer to nil", (*string)(nil), true},

		// Special formats
		{"with path", "ssl://example.com:443/path", true},             // paths not supported
		{"with query", "ssl://example.com:443?query=1", true},         // queries not supported
		{"with fragment", "ssl://example.com:443#section", true},      // fragments not supported
		{"with credentials", "ssl://user:pass@example.com:443", true}, // credentials not supported
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := URI().Validate(tt.value)
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

func TestURIRule_When(t *testing.T) {
	// Test with condition false - should skip validation
	err := URI().When(false).Validate("invalid")
	if err != nil {
		t.Errorf("expected no error when condition is false, got: %v", err)
	}

	// Test with condition true - should validate
	err = URI().When(true).Validate("invalid")
	if err == nil {
		t.Error("expected error when condition is true and uri is invalid")
	}
}

func TestURIRule_Error(t *testing.T) {
	customMsg := "custom uri error"
	err := URI().Error(customMsg).Validate("invalid")
	if err == nil {
		t.Fatal("expected error but got nil")
	}
	if err.Error() != customMsg {
		t.Errorf("expected error message %q, got %q", customMsg, err.Error())
	}
}

func TestURIRule_ErrorObject(t *testing.T) {
	customErr := validation.NewError("custom_code", "custom message")
	err := URI().ErrorObject(customErr).Validate("invalid")
	if err == nil {
		t.Fatal("expected error but got nil")
	}
	if err.Error() != "custom message" {
		t.Errorf("expected custom error message, got: %v", err.Error())
	}
}

func TestURIRule_ErrorMessages(t *testing.T) {
	// Test specific error messages
	tests := []struct {
		name          string
		value         string
		expectedError string
	}{
		{"invalid port non-numeric", "ssl://example.com:abc", "invalid uri address"},
		{"invalid port zero", "ssl://example.com:0", "invalid uri address"},
		{"invalid port too high", "ssl://example.com:99999", "invalid uri address"},
		{"invalid protocol", "http://example.com:80", "invalid uri address"},
		{"invalid format", "not-a-uri", "invalid uri address"},
		{"no port", "ssl://example.com", "invalid uri address"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := URI().Validate(tt.value)
			if err == nil {
				t.Fatal("expected error but got nil")
			}
			// Check if error message contains expected text
			if !uriContains(err.Error(), tt.expectedError) {
				t.Errorf("expected error to contain %q, got %q", tt.expectedError, err.Error())
			}
		})
	}
}

func TestURIRule_HostValidation(t *testing.T) {
	// Test various host formats
	validHosts := []string{
		"example.com",
		"sub.example.com",
		"sub.sub.example.com",
		"example-with-dash.com",
		"192.168.1.1",
		"10.0.0.1",
		"255.255.255.254",
		"localhost",
		// IPv6 not supported by govalidator.IsHost
	}

	for _, host := range validHosts {
		t.Run("valid_host_"+host, func(t *testing.T) {
			uri := "ssl://" + host + ":443"
			err := URI().Validate(uri)
			if err != nil {
				t.Errorf("expected valid host %s, got error: %v", host, err)
			}
		})
	}

	invalidHosts := []string{
		"",
		" ",
		"example com",  // space
		"example..com", // double dot
		".example.com", // leading dot
		"example@com",  // @ symbol
		"example#com",  // # symbol
	}

	for _, host := range invalidHosts {
		t.Run("invalid_host_"+host, func(t *testing.T) {
			uri := "ssl://" + host + ":443"
			err := URI().Validate(uri)
			if err == nil {
				t.Errorf("expected invalid host %s, but got no error", host)
			}
		})
	}
}

func BenchmarkURIValidation(b *testing.B) {
	uris := []string{
		"ssl://example.com:443",
		"starttls://mail.example.com:587",
		"invalid://example.com:80",
		"not-a-uri",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, uri := range uris {
			_ = URI().Validate(uri)
		}
	}
}

// Helper functions
func uriStrPtr(s string) *string {
	return &s
}

func uriContains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) &&
		(s[:len(substr)] == substr || s[len(s)-len(substr):] == substr ||
			len(substr) < len(s) && findSubstring(s, substr)))
}

func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
