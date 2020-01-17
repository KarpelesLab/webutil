package webutil

import (
	"net"
	"strconv"
	"strings"
)

// ParseIPPort will parse an IP with optionally a port
func ParseIPPort(ip string) *net.TCPAddr {
	if len(ip) < 2 {
		// can't parse something that small
		return nil
	}

	res := &net.TCPAddr{}
	pos := strings.LastIndex(ip, ":") // is there a port?

	if ip[0] == '[' {
		ip_end := strings.Index(ip, "]")
		if ip_end == -1 {
			return nil
		}
		if ip_end != pos-1 {
			res.IP = net.ParseIP(ip[1 : len(ip)-1])
			if res.IP != nil {
				return res
			}
			return nil // :(
		}

		res.IP = net.ParseIP(ip[1:ip_end])
		if res.IP == nil {
			return nil
		}
		ip = ip[pos:]
	} else if pos > 0 {
		res.IP = net.ParseIP(ip[0:pos])
		if res.IP == nil {
			res.IP = net.ParseIP(ip)
			if res.IP != nil {
				return res
			}
			return nil
		}
		ip = ip[pos:]
	} else if pos == -1 {
		res.IP = net.ParseIP(ip)
		if res.IP == nil {
			return nil
		}
		return res
	}

	// only remains is ":port"
	if ip[0] != ':' {
		return nil
	}

	port, err := strconv.ParseUint(ip[1:], 10, 16)
	if err != nil {
		return nil
	}

	res.Port = int(port)
	return res
}
