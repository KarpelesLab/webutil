package webutil_test

import (
	"net"
	"testing"

	"github.com/KarpelesLab/webutil"
)

func TestParseIPPort(t *testing.T) {
	testCases := []struct {
		name   string
		input  string
		expect *net.TCPAddr
	}{
		{
			name:   "IPv4 with port",
			input:  "127.0.0.1:80",
			expect: &net.TCPAddr{IP: net.ParseIP("127.0.0.1"), Port: 80},
		},
		{
			name:   "IPv6 with port",
			input:  "[::1]:80",
			expect: &net.TCPAddr{IP: net.ParseIP("::1"), Port: 80},
		},
		{
			name:   "Port only",
			input:  ":80",
			expect: &net.TCPAddr{Port: 80},
		},
		{
			name:   "IPv4 only",
			input:  "127.0.0.1",
			expect: &net.TCPAddr{IP: net.ParseIP("127.0.0.1")},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := webutil.ParseIPPort(tc.input)

			if result == nil {
				t.Fatalf("ParseIPPort returned nil for %q", tc.input)
			}

			if result.String() != tc.expect.String() {
				t.Errorf("ParseIPPort(%q): got %s, want %s",
					tc.input, result.String(), tc.expect.String())
			}
		})
	}

	// Test invalid inputs
	invalidCases := []struct {
		name  string
		input string
	}{
		{
			name:  "Empty string",
			input: "",
		},
		{
			name:  "Invalid IP format",
			input: "not-an-ip",
		},
		{
			name:  "Invalid port",
			input: "127.0.0.1:xyz",
		},
	}

	for _, tc := range invalidCases {
		t.Run(tc.name, func(t *testing.T) {
			result := webutil.ParseIPPort(tc.input)

			if result != nil && tc.name != "Empty string" {
				t.Errorf("Expected nil for invalid input %q, got %v", tc.input, result)
			}
		})
	}
}
