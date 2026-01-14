// Copyright 2021-2026 Altessa Solutions Inc. All rights reserved.
// Use of this source code is governed by license that can be found in
// the LICENSE file.

package ozzo_rules

import (
	"testing"

	"github.com/go-ozzo/ozzo-validation/v4"
)

func TestListenAddressRule(t *testing.T) {
	tests := []struct {
		name      string
		value     any
		expectErr bool
	}{
		// Valid listen addresses
		{"valid localhost with port", "localhost:8080", false},
		{"valid IP with port", "192.168.1.1:80", false},
		{"valid 0.0.0.0 with port", "0.0.0.0:3000", false},
		{"valid domain with port", "example.com:443", false},
		{"valid subdomain with port", "api.example.com:8443", false},
		{"valid port min", "localhost:1", false},
		{"valid port max", "localhost:65534", false},

		// Valid with protocol prefix (should be stripped)
		{"valid http prefix", "http://localhost:8080", false},
		{"valid https prefix", "https://example.com:443", false},
		{"valid tcp prefix", "tcp://0.0.0.0:9000", false},
		{"valid udp prefix", "udp://127.0.0.1:5353", false},

		// Invalid addresses - port issues
		{"no port", "localhost", true},
		{"invalid port zero", "localhost:0", true},
		{"invalid port 65535", "localhost:65535", true},
		{"invalid port too high", "localhost:65536", true},
		{"invalid port negative", "localhost:-1", true},
		{"invalid port non-numeric", "localhost:http", true},
		{"invalid port float", "localhost:80.5", true},

		// Invalid addresses - host issues
		{"empty host", ":8080", true},
		{"invalid host spaces", "local host:8080", true},
		{"invalid host special chars", "local@host:8080", true},
		{"multiple colons", "localhost:8080:extra", true},

		// Edge cases
		{"empty string", "", true},
		{"only protocol", "http://", true},
		{"nil value", nil, true},
		{"non-string value", 123, true},
		{"pointer to valid address", listenAddressStrPtr("localhost:8080"), false},
		{"pointer to invalid address", listenAddressStrPtr("invalid"), true},
		{"pointer to nil", (*string)(nil), true},

		// IPv6 addresses (might not be supported by govalidator.IsHost)
		{"IPv6 localhost", "[::1]:8080", true},      // likely unsupported
		{"IPv6 address", "[2001:db8::1]:443", true}, // likely unsupported
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ListenAddress().Validate(tt.value)
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

func TestListenAddressRule_When(t *testing.T) {
	err := ListenAddress().When(false).Validate("invalid")
	if err != nil {
		t.Errorf("expected no error when condition is false, got: %v", err)
	}

	err = ListenAddress().When(true).Validate("invalid")
	if err == nil {
		t.Error("expected error when condition is true and address is invalid")
	}
}

func TestListenAddressRule_Error(t *testing.T) {
	customMsg := "custom address error"
	err := ListenAddress().Error(customMsg).Validate("invalid")
	if err == nil {
		t.Fatal("expected error but got nil")
	}
	if err.Error() != customMsg {
		t.Errorf("expected error message %q, got %q", customMsg, err.Error())
	}
}

func TestListenAddressRule_ErrorObject(t *testing.T) {
	customErr := validation.NewError("custom_code", "custom message")
	err := ListenAddress().ErrorObject(customErr).Validate("invalid")
	if err == nil {
		t.Fatal("expected error but got nil")
	}
	if err.Error() != "custom message" {
		t.Errorf("expected custom error message, got: %v", err.Error())
	}
}

func TestListenAddressRule_ProtocolStripping(t *testing.T) {
	// Test that various protocols are properly stripped
	protocols := []string{
		"http", "https", "tcp", "udp", "ws", "wss", "ftp", "ssh",
	}

	for _, proto := range protocols {
		t.Run("protocol_"+proto, func(t *testing.T) {
			address := proto + "://localhost:8080"
			err := ListenAddress().Validate(address)
			if err != nil {
				t.Errorf("expected valid address with %s protocol, got error: %v", proto, err)
			}
		})
	}

	// Test case sensitivity
	err := ListenAddress().Validate("HTTP://localhost:8080")
	if err != nil {
		t.Error("expected protocol stripping to be case insensitive")
	}
}

func TestListenAddressRule_HostValidation(t *testing.T) {
	validHosts := []string{
		"localhost",
		"example.com",
		"sub.example.com",
		"example-with-dash.com",
		"192.168.1.1",
		"10.0.0.1",
		"255.255.255.254",
		"0.0.0.0",
		"127.0.0.1",
	}

	for _, host := range validHosts {
		t.Run("valid_host_"+host, func(t *testing.T) {
			address := host + ":8080"
			err := ListenAddress().Validate(address)
			if err != nil {
				t.Errorf("expected valid host %s, got error: %v", host, err)
			}
		})
	}

	invalidHosts := []string{
		"",
		" ",
		"example com",
		"example..com",
		".example.com",
		"example@com",
		"example#com",
		"_example.com", // underscore in hostname
	}

	for _, host := range invalidHosts {
		t.Run("invalid_host_"+host, func(t *testing.T) {
			address := host + ":8080"
			err := ListenAddress().Validate(address)
			if err == nil {
				t.Errorf("expected invalid host %s, but got no error", host)
			}
		})
	}
}

func TestListenAddressRule_PortRanges(t *testing.T) {
	// Test specific port boundaries
	portTests := []struct {
		port  string
		valid bool
	}{
		{"0", false},     // too low
		{"1", true},      // minimum valid
		{"80", true},     // common HTTP
		{"443", true},    // common HTTPS
		{"3000", true},   // common dev port
		{"8080", true},   // common alt HTTP
		{"65534", true},  // maximum valid
		{"65535", false}, // too high
		{"65536", false}, // way too high
		{"99999", false}, // extremely high
		{"-1", false},    // negative
		{"-80", false},   // negative common
	}

	for _, pt := range portTests {
		t.Run("port_"+pt.port, func(t *testing.T) {
			address := "localhost:" + pt.port
			err := ListenAddress().Validate(address)
			if pt.valid && err != nil {
				t.Errorf("expected valid port %s, got error: %v", pt.port, err)
			} else if !pt.valid && err == nil {
				t.Errorf("expected invalid port %s, but got no error", pt.port)
			}
		})
	}
}

func BenchmarkListenAddressValidation(b *testing.B) {
	addresses := []string{
		"localhost:8080",
		"http://example.com:443",
		"192.168.1.1:3000",
		"invalid",
		":8080",
		"localhost:99999",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, addr := range addresses {
			_ = ListenAddress().Validate(addr)
		}
	}
}

// Helper function
func listenAddressStrPtr(s string) *string {
	return &s
}
