package webutil_test

import (
	"net"
	"testing"

	"github.com/KarpelesLab/webutil"
)

func TestParseIPPort(t *testing.T) {
	var cases = []struct {
		A string
		R *net.TCPAddr
	}{
		{"127.0.0.1:80", &net.TCPAddr{IP: net.ParseIP("127.0.0.1"), Port: 80}},
		{"[::1]:80", &net.TCPAddr{IP: net.ParseIP("::1"), Port: 80}},
		{":80", &net.TCPAddr{Port: 80}},
		{"127.0.0.1", &net.TCPAddr{IP: net.ParseIP("127.0.0.1")}},
	}

	for _, c := range cases {
		r := webutil.ParseIPPort(c.A)

		if r.String() != c.R.String() {
			t.Errorf("test failed, %s result %s, expected %s", c.A, r, c.R)
		}
	}
}
