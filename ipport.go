package webutil

import (
	"net"
	"strconv"
)

// ParseIPPort will parse an IP with optionally a port
func ParseIPPort(ip string) *net.TCPAddr {
	host, port, err := net.SplitHostPort(ip)
	if err != nil {
		// fallback on ip only?
		if ip := net.ParseIP(ip); ip != nil {
			return &net.TCPAddr{IP: ip}
		}
		// can't parse something that small
		return nil
	}
	portN, err := strconv.ParseUint(port, 10, 16)
	if err != nil {
		return nil
	}

	res := &net.TCPAddr{
		IP:   net.ParseIP(host),
		Port: int(portN),
	}

	return res
}
