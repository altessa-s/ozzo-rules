// Copyright 2021-2026 Altessa Solutions Inc. All rights reserved.
// Use of this source code is governed by license that can be found in
// the LICENSE file.

package ozzo_rules

import (
	"testing"

	"github.com/go-ozzo/ozzo-validation/v4"
)

// TestConsulServiceAddressRule tests the ConsulServiceAddress validator
func TestConsulServiceAddressRule(t *testing.T) {
	tests := []struct {
		name      string
		value     any
		expectErr bool
	}{
		// Valid IP addresses
		{"valid IPv4 localhost", "127.0.0.1", false},
		{"valid IPv4 private", "192.168.1.1", false},
		{"valid IPv4 public", "8.8.8.8", false},
		{"valid IPv4 max", "255.255.255.255", false},

		// Valid DNS names
		{"valid simple hostname", "consul", false},
		{"valid FQDN", "consul.service.consul", false},
		{"valid with hyphen", "consul-server", false},
		{"valid with underscore", "consul_server", false},
		{"valid numeric start", "1consul", false},
		{"valid long segment", "verylonghostnamethatisalmostthemaximumlengthallowed123456", false},
		{"valid multiple segments", "consul.service.dc1.internal", false},
		{"valid ending with dot", "consul.service.", false},

		// Invalid IP addresses
		{"invalid unspecified IPv4", "0.0.0.0", true},
		{"invalid unspecified IPv6", "::", true},
		{"invalid IPv4 format", "256.256.256.256", true},
		{"invalid IPv4 incomplete", "192.168.1", true},
		{"invalid IPv4 with port", "192.168.1.1:8500", true},

		// Invalid DNS names
		{"invalid empty", "", true},
		{"invalid starts with hyphen", "-consul", true},
		{"invalid ends with hyphen", "consul-", true},
		{"invalid double hyphen start", "--consul", true},
		{"invalid special chars", "consul@service", true},
		{"invalid space", "consul service", true},
		{"invalid double dot", "consul..service", true},
		{"invalid starts with dot", ".consul", true},
		{"invalid segment too long", "verylonghostnamethatexceedsthemaximumlengthallowedforasegment12345", true},

		// Edge cases
		{"nil value", nil, true},
		{"integer value", 12345, true},
		{"boolean value", true, true},
		{"slice value", []string{"consul"}, true},

		// Pointer handling
		{"pointer to valid IP", consulStrPtr("192.168.1.1"), false},
		{"pointer to valid DNS", consulStrPtr("consul.service"), false},
		{"pointer to invalid", consulStrPtr("0.0.0.0"), true},
		{"pointer to nil", (*string)(nil), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ConsulServiceAddress().Validate(tt.value)
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

// TestConsulServiceAddressOrTemplateRule tests the ConsulServiceAddressOrTemplate validator
func TestConsulServiceAddressOrTemplateRule(t *testing.T) {
	tests := []struct {
		name      string
		value     any
		expectErr bool
	}{
		// Valid addresses (same as regular validator)
		{"valid IP", "192.168.1.1", false},
		{"valid DNS", "consul.service", false},

		// Valid templates
		{"valid template simple", "{{.Address}}", false},
		{"valid template with text", "consul-{{.Datacenter}}", false},
		{"valid template complex", "{{.Node}}.{{.Datacenter}}.consul", false},
		{"valid template IP", "10.0.{{.Index}}.1", false},
		{"valid double braces", "service-{{.Name}}-{{.Port}}", false},

		// Invalid templates (still need valid structure outside templates)
		{"invalid template only dots", "{{.}}.{{.}}", true},
		{"invalid unspecified even with template", "0.0.0.0", true},

		// Edge cases with templates
		{"single open brace", "{.Address}", true},
		{"single close brace", "}.Address{", true},
		{"incomplete template", "{{.Address", true},
		{"empty template", "{{}}", false},

		// Standard edge cases
		{"nil value", nil, true},
		{"empty string", "", true},
		{"integer value", 12345, true},

		// Pointer handling
		{"pointer to valid template", consulStrPtr("{{.Address}}"), false},
		{"pointer to invalid", consulStrPtr("invalid.."), true},
		{"pointer to nil", (*string)(nil), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ConsulServiceAddressOrTemplate().Validate(tt.value)
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

// TestConsulServiceNameRule tests the ConsulServiceName validator
func TestConsulServiceNameRule(t *testing.T) {
	tests := []struct {
		name      string
		value     any
		expectErr bool
	}{
		// Valid service names
		{"valid simple", "consul", false},
		{"valid with hyphen", "consul-server", false},
		{"valid alphanumeric", "consul2", false},
		{"valid starts with number", "1consul", false},
		{"valid FQDN style", "api.v1.service", false},
		{"valid multiple segments", "web.frontend.prod.service", false},
		{"valid uppercase", "ConsulService", false},
		{"valid mixed case", "MyConsulService", false},
		{"valid long name", "very-long-service-name-that-is-still-valid-123", false},

		// Invalid service names
		{"invalid empty", "", true},
		{"invalid starts with hyphen", "-consul", true},
		{"invalid ends with hyphen", "consul-", true},
		{"invalid starts with dot", ".consul", true},
		{"invalid ends with dot", "consul.", true},
		{"invalid double hyphen", "consul--service", true},
		{"invalid double dot", "consul..service", true},
		{"invalid special chars @", "consul@service", true},
		{"invalid special chars _", "consul_service", true},
		{"invalid space", "consul service", true},
		{"invalid hyphen at segment start", "consul.-service", true},
		{"invalid hyphen at segment end", "consul-.service", true},

		// Edge cases
		{"nil value", nil, true},
		{"integer value", 12345, true},
		{"boolean value", true, true},
		{"slice value", []string{"consul"}, true},

		// Pointer handling
		{"pointer to valid", consulStrPtr("consul-service"), false},
		{"pointer to invalid", consulStrPtr("-invalid"), true},
		{"pointer to nil", (*string)(nil), true},

		// Single character segments
		{"single char", "a", false},
		{"single char segments", "a.b.c", false},
		{"single digit", "1", false},
		{"single digit segments", "1.2.3", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ConsulServiceName().Validate(tt.value)
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

// TestConsulServiceNameOrNilRule tests the ConsulServiceNameOrNil validator
func TestConsulServiceNameOrNilRule(t *testing.T) {
	tests := []struct {
		name      string
		value     any
		expectErr bool
	}{
		// Valid cases including nil/empty
		{"valid service name", "consul", false},
		{"valid with hyphen", "consul-server", false},
		{"valid nil", nil, false},
		{"valid empty string", "", false},
		{"valid pointer nil", (*string)(nil), false},

		// Invalid cases (same as regular validator except nil/empty)
		{"invalid starts with hyphen", "-consul", true},
		{"invalid special chars", "consul@service", true},
		{"invalid space", "consul service", true},

		// Type errors
		{"invalid integer", 12345, true},
		{"invalid boolean", true, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ConsulServiceNameOrNil().Validate(tt.value)
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

// Test When method for all Consul validators
func TestConsulValidators_When(t *testing.T) {
	testCases := []struct {
		name      string
		validator validation.Rule
		value     string
	}{
		{"ConsulServiceAddress", ConsulServiceAddress(), "invalid.."},
		{"ConsulServiceAddressOrTemplate", ConsulServiceAddressOrTemplate(), "invalid.."},
		{"ConsulServiceName", ConsulServiceName(), "-invalid"},
		{"ConsulServiceNameOrNil", ConsulServiceNameOrNil(), "-invalid"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var err error
			switch v := tc.validator.(type) {
			case ConsulServiceAddressRule:
				err = v.When(false).Validate(tc.value)
			case ConsulServiceNameRule:
				err = v.When(false).Validate(tc.value)
			}

			if err != nil {
				t.Errorf("expected no error when condition is false, got: %v", err)
			}
		})
	}
}

// Test Error method for all Consul validators
func TestConsulValidators_Error(t *testing.T) {
	customMsg := "custom consul error"

	testCases := []struct {
		name      string
		validator validation.Rule
		value     string
	}{
		{"ConsulServiceAddress", ConsulServiceAddress(), "0.0.0.0"},
		{"ConsulServiceAddressOrTemplate", ConsulServiceAddressOrTemplate(), "invalid.."},
		{"ConsulServiceName", ConsulServiceName(), "-invalid"},
		{"ConsulServiceNameOrNil", ConsulServiceNameOrNil(), "-invalid"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var err error
			switch v := tc.validator.(type) {
			case ConsulServiceAddressRule:
				err = v.Error(customMsg).Validate(tc.value)
			case ConsulServiceNameRule:
				err = v.Error(customMsg).Validate(tc.value)
			}

			if err == nil {
				t.Fatal("expected error but got nil")
			}
			if err.Error() != customMsg {
				t.Errorf("expected error message %q, got %q", customMsg, err.Error())
			}
		})
	}
}

// Test ErrorObject method for all Consul validators
func TestConsulValidators_ErrorObject(t *testing.T) {
	customErr := validation.NewError("custom_code", "custom message")

	testCases := []struct {
		name      string
		validator validation.Rule
		value     string
	}{
		{"ConsulServiceAddress", ConsulServiceAddress(), "0.0.0.0"},
		{"ConsulServiceAddressOrTemplate", ConsulServiceAddressOrTemplate(), "invalid.."},
		{"ConsulServiceName", ConsulServiceName(), "-invalid"},
		{"ConsulServiceNameOrNil", ConsulServiceNameOrNil(), "-invalid"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var err error
			switch v := tc.validator.(type) {
			case ConsulServiceAddressRule:
				err = v.ErrorObject(customErr).Validate(tc.value)
			case ConsulServiceNameRule:
				err = v.ErrorObject(customErr).Validate(tc.value)
			}

			if err == nil {
				t.Fatal("expected error but got nil")
			}
			if err.Error() != "custom message" {
				t.Errorf("expected custom error message, got: %v", err.Error())
			}
		})
	}
}

// Test common Consul service scenarios
func TestConsulValidators_CommonScenarios(t *testing.T) {
	t.Run("Service Discovery", func(t *testing.T) {
		// Common Consul service names
		validNames := []string{
			"web",
			"api",
			"database",
			"cache-service",
			"auth-service",
			"user-api-v2",
			"payment-gateway",
			"notification.service",
			"backend.prod.service",
		}

		for _, name := range validNames {
			err := ConsulServiceName().Validate(name)
			if err != nil {
				t.Errorf("expected valid service name %q, got error: %v", name, err)
			}
		}
	})

	t.Run("Service Addresses", func(t *testing.T) {
		// Common Consul service addresses
		validAddresses := []string{
			"consul.service.consul",
			"web.service.consul",
			"10.0.0.100",
			"192.168.1.50",
			"consul-server-1",
			"consul-server-1.dc1",
			"api.prod.internal",
		}

		for _, addr := range validAddresses {
			err := ConsulServiceAddress().Validate(addr)
			if err != nil {
				t.Errorf("expected valid service address %q, got error: %v", addr, err)
			}
		}
	})

	t.Run("Template Addresses", func(t *testing.T) {
		// Common Consul template patterns
		validTemplates := []string{
			"{{.Address}}",
			"{{.Node}}.node.consul",
			"10.0.{{.Index}}.{{.Port}}",
			"{{.Service}}.service.{{.Datacenter}}.consul",
			"{{range service \"web\"}}{{.Address}}{{end}}",
		}

		for _, tmpl := range validTemplates {
			err := ConsulServiceAddressOrTemplate().Validate(tmpl)
			if err != nil {
				t.Errorf("expected valid template %q, got error: %v", tmpl, err)
			}
		}
	})
}

// Benchmark tests
func BenchmarkConsulServiceAddressValidation(b *testing.B) {
	values := []string{
		"192.168.1.1",
		"consul.service.consul",
		"0.0.0.0",
		"invalid..name",
		"",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, v := range values {
			_ = ConsulServiceAddress().Validate(v)
		}
	}
}

func BenchmarkConsulServiceNameValidation(b *testing.B) {
	values := []string{
		"consul",
		"web-service",
		"api.v1.service",
		"-invalid",
		"",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, v := range values {
			_ = ConsulServiceName().Validate(v)
		}
	}
}

// Helper function
func consulStrPtr(s string) *string {
	return &s
}
