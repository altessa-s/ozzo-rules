// Copyright 2021-2026 Altessa Solutions Inc. All rights reserved.
// Use of this source code is governed by license that can be found in
// the LICENSE file.

package ozzo_rules

import (
	"strings"
	"testing"

	"github.com/go-ozzo/ozzo-validation/v4"
)

func TestServerAddressRule(t *testing.T) {
	tests := []struct {
		name      string
		value     any
		expectErr bool
	}{
		// Valid server addresses
		{"valid localhost with port", "localhost:8080", false},
		{"valid IP with port", "192.168.1.1:80", false},
		{"valid domain with port", "example.com:443", false},
		{"valid subdomain with port", "api.example.com:8080", false},
		{"valid with http prefix", "http://example.com:80", false},
		{"valid with https prefix", "https://example.com:443", false},
		{"valid with ftp prefix", "ftp://example.com:21", false},
		{"valid with custom prefix", "tcp://example.com:1234", false},
		{"valid IPv4 full", "255.255.255.255:65534", false},
		{"valid port 1", "example.com:1", false},
		{"valid port 65534", "example.com:65534", false},
		{"valid numeric domain", "192.168.1.1:3000", false},
		{"valid hyphenated domain", "my-server.example.com:8080", false},
		{"valid long domain", "very.long.subdomain.example.com:443", false},

		// Invalid server addresses - missing port
		{"missing port", "example.com", true},
		{"missing port with protocol", "http://example.com", true},
		{"missing port localhost", "localhost", true},
		{"missing port IP", "192.168.1.1", true},

		// Invalid server addresses - invalid port
		{"port 0", "example.com:0", true},
		{"port 65535", "example.com:65535", true},
		{"port above max", "example.com:65536", true},
		{"port way above max", "example.com:99999", true},
		{"negative port", "example.com:-80", true},
		{"non-numeric port", "example.com:http", true},
		{"decimal port", "example.com:80.5", true},
		{"empty port", "example.com:", true},

		// Invalid server addresses - invalid host
		{"empty host", ":8080", true},
		{"invalid host chars", "exam ple.com:8080", true},
		{"invalid host special", "example@com:8080", true},
		{"invalid host brackets", "example[com]:8080", true},
		{"double dots in host", "example..com:8080", true},
		{"starts with dot", ".example.com:8080", true},
		{"ends with dot", "example.com.:8080", false}, // govalidator allows trailing dot
		{"invalid IP", "999.999.999.999:8080", false}, // govalidator doesn't validate IP ranges
		{"partial IP", "192.168.1:8080", false},       // govalidator allows this as hostname
		{"IPv6 without brackets", "::1:8080", true},
		{"IPv6 with brackets", "[::1]:8080", true}, // Note: govalidator.IsHost doesn't support IPv6

		// Invalid server addresses - multiple colons
		{"multiple colons", "example.com:80:443", true},
		{"triple colon", "example.com:::8080", true},

		// Invalid server addresses - with path
		{"with path", "example.com:8080/path", true},
		{"with query", "example.com:8080?query", true},
		{"with fragment", "example.com:8080#fragment", true},

		// Edge cases
		{"nil value", nil, true},
		{"empty string", "", true},
		{"non-string value", 8080, true},
		{"pointer to valid", serverAddressStrPtr("example.com:8080"), false},
		{"pointer to invalid", serverAddressStrPtr("example.com"), true},
		{"pointer to nil", (*string)(nil), true},

		// Protocol edge cases
		{"uppercase protocol", "HTTP://example.com:8080", false},
		{"mixed case protocol", "HtTp://example.com:8080", false},
		{"protocol only", "http://", true},
		{"double protocol", "http://https://example.com:8080", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ServerAddress().Validate(tt.value)
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

func TestServerAddressRule_When(t *testing.T) {
	err := ServerAddress().When(false).Validate("invalid-no-port")
	if err != nil {
		t.Errorf("expected no error when condition is false, got: %v", err)
	}

	err = ServerAddress().When(true).Validate("invalid-no-port")
	if err == nil {
		t.Error("expected error when condition is true and address is invalid")
	}
}

func TestServerAddressRule_Error(t *testing.T) {
	customMsg := "custom server address error"
	err := ServerAddress().Error(customMsg).Validate("no-port")
	if err == nil {
		t.Fatal("expected error but got nil")
	}
	// Note: The actual error includes additional info, so we check if it contains our message
	if !strings.Contains(err.Error(), customMsg) {
		t.Errorf("expected error to contain %q, got %q", customMsg, err.Error())
	}
}

func TestServerAddressRule_ErrorObject(t *testing.T) {
	customErr := validation.NewError("custom_code", "custom message")
	err := ServerAddress().ErrorObject(customErr).Validate("no-port")
	if err == nil {
		t.Fatal("expected error but got nil")
	}
	// The error is wrapped with additional info
	if !strings.Contains(err.Error(), "custom message") {
		t.Errorf("expected custom error message in error: %v", err.Error())
	}
}

func TestServerAddressRule_CommonScenarios(t *testing.T) {
	// Common server address scenarios
	scenarios := []struct {
		name        string
		description string
		address     string
		expectErr   bool
	}{
		// Web servers
		{
			"web server HTTP",
			"standard HTTP server",
			"web.example.com:80",
			false,
		},
		{
			"web server HTTPS",
			"standard HTTPS server",
			"secure.example.com:443",
			false,
		},
		{
			"web server custom port",
			"web server on custom port",
			"app.example.com:8080",
			false,
		},

		// Database servers
		{
			"MySQL default",
			"MySQL on default port",
			"mysql.example.com:3306",
			false,
		},
		{
			"PostgreSQL default",
			"PostgreSQL on default port",
			"postgres.example.com:5432",
			false,
		},
		{
			"Redis default",
			"Redis on default port",
			"redis.example.com:6379",
			false,
		},
		{
			"MongoDB default",
			"MongoDB on default port",
			"mongo.example.com:27017",
			false,
		},

		// Development servers
		{
			"localhost dev",
			"local development server",
			"localhost:3000",
			false,
		},
		{
			"localhost high port",
			"local server on high port",
			"localhost:9999",
			false,
		},
		{
			"127.0.0.1 dev",
			"local IP development",
			"127.0.0.1:8000",
			false,
		},

		// Microservices
		{
			"microservice internal",
			"internal microservice",
			"user-service:50051",
			false,
		},
		{
			"microservice k8s",
			"Kubernetes service",
			"api-gateway.default.svc.cluster.local:8080",
			false,
		},

		// Load balancers
		{
			"load balancer",
			"AWS ELB address",
			"my-lb-1234567890.us-west-2.elb.amazonaws.com:443",
			false,
		},

		// Invalid common mistakes
		{
			"URL with path",
			"full URL instead of address",
			"http://example.com:8080/api/v1",
			true,
		},
		{
			"missing port common",
			"forgot to add port",
			"api.example.com",
			true,
		},
		{
			"localhost no port",
			"localhost without port",
			"localhost",
			true,
		},
	}

	for _, sc := range scenarios {
		t.Run(sc.name, func(t *testing.T) {
			err := ServerAddress().Validate(sc.address)
			if sc.expectErr && err == nil {
				t.Errorf("%s: expected error but got nil", sc.description)
			} else if !sc.expectErr && err != nil {
				t.Errorf("%s: expected no error but got: %v", sc.description, err)
			}
		})
	}
}

func TestServerAddressRule_ProtocolHandling(t *testing.T) {
	// Test various protocol prefixes
	protocols := []struct {
		name      string
		address   string
		expectErr bool
	}{
		// Standard protocols
		{"http lowercase", "http://example.com:80", false},
		{"https lowercase", "https://example.com:443", false},
		{"ftp lowercase", "ftp://ftp.example.com:21", false},
		{"ssh lowercase", "ssh://ssh.example.com:22", false},

		// Custom protocols
		{"tcp", "tcp://example.com:1234", false},
		{"udp", "udp://example.com:1234", false},
		{"grpc", "grpc://api.example.com:50051", false},
		{"custom", "myprotocol://example.com:9999", false},

		// Case variations
		{"HTTP uppercase", "HTTP://example.com:80", false},
		{"HtTpS mixed", "HtTpS://example.com:443", false},

		// Invalid protocol usage
		{"double slash no protocol", "//example.com:80", true},
		{"protocol no slashes", "http:example.com:80", true},
		{"numbers in protocol", "http2://example.com:80", true},          // protocol regex only allows letters
		{"underscore in protocol", "my_protocol://example.com:80", true}, // RFC 3986: scheme = ALPHA *( ALPHA / DIGIT / "+" / "-" / "." )
	}

	for _, tc := range protocols {
		t.Run(tc.name, func(t *testing.T) {
			err := ServerAddress().Validate(tc.address)
			if tc.expectErr && err == nil {
				t.Error("expected error but got nil")
			} else if !tc.expectErr && err != nil {
				t.Errorf("expected no error but got: %v", err)
			}
		})
	}
}

func TestServerAddressRule_PortBoundaries(t *testing.T) {
	// Test port boundary values
	ports := []struct {
		name      string
		port      string
		expectErr bool
	}{
		{"port 1", "1", false},
		{"port 80", "80", false},
		{"port 443", "443", false},
		{"port 1024", "1024", false},
		{"port 8080", "8080", false},
		{"port 32768", "32768", false},
		{"port 65534", "65534", false},

		// Invalid ports
		{"port 0", "0", true},
		{"port -1", "-1", true},
		{"port 65535", "65535", true},
		{"port 65536", "65536", true},
		{"port 100000", "100000", true},
	}

	for _, tc := range ports {
		t.Run(tc.name, func(t *testing.T) {
			address := "example.com:" + tc.port
			err := ServerAddress().Validate(address)
			if tc.expectErr && err == nil {
				t.Errorf("expected error for port %s but got nil", tc.port)
			} else if !tc.expectErr && err != nil {
				t.Errorf("expected no error for port %s but got: %v", tc.port, err)
			}
		})
	}
}

func BenchmarkServerAddressValidation(b *testing.B) {
	addresses := []string{
		"example.com:8080",
		"http://example.com:80",
		"https://secure.example.com:443",
		"localhost:3000",
		"192.168.1.1:22",
		"invalid-no-port",
		"example.com:0",
		"example.com:99999",
		":8080",
		"example.com:8080/path",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, addr := range addresses {
			_ = ServerAddress().Validate(addr)
		}
	}
}

// Helper function to avoid conflicts
func serverAddressStrPtr(s string) *string {
	return &s
}
