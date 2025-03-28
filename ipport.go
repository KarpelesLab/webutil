package webutil

import (
	"net"
	"strconv"
)

// ParseIPPort parses a string containing an IP address with an optional port.
//
// The input can be in the following formats:
//   - "127.0.0.1:80" -> IP with port
//   - "127.0.0.1" -> IP only
//   - "[::1]:80" -> IPv6 with port
//   - "::1" -> IPv6 only
//   - ":80" -> Empty host with port
//
// Returns nil if the input cannot be parsed as a valid IP/port combination.
func ParseIPPort(ipStr string) *net.TCPAddr {
	// Try to split into host and port parts
	host, port, err := net.SplitHostPort(ipStr)
	if err != nil {
		// No port in the string, try to parse as IP only
		if ip := net.ParseIP(ipStr); ip != nil {
			return &net.TCPAddr{IP: ip}
		}
		// Can't parse as IP only
		return nil
	}

	// Parse port number
	portN, err := strconv.ParseUint(port, 10, 16)
	if err != nil {
		return nil
	}

	// Parse IP (host may be empty for ":port" format)
	var ip net.IP
	if host != "" {
		ip = net.ParseIP(host)
	}

	return &net.TCPAddr{
		IP:   ip,
		Port: int(portN),
	}
}
