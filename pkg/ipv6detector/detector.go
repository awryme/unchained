package ipv6detector

import (
	"context"
	"fmt"
	"net"
	"strings"
	"time"
)

const TestServer = "google.com"

const DefaultTimeout = time.Second * 10

type Logf func(format string, args ...any)

func Detect(ctx context.Context, logf Logf) bool {
	ok := HasInterfaceIPv6(logf)
	if !ok {
		return false
	}

	return DialTestServer(ctx, logf)
}

// HasInterfaceIPv6 checks if network has public (not loopback) ipv6 interface
func HasInterfaceIPv6(logf Logf) bool {
	if logf == nil {
		logf = func(format string, args ...any) {}
	}

	ifs, err := net.Interfaces()
	if err != nil {
		logf("list interfaces: %v", err)
		return false
	}

	for _, iface := range ifs {
		// ignore loopbacks
		if iface.Flags&net.FlagLoopback == 1 {
			continue
		}
		addrs, err := iface.Addrs()
		if err != nil {
			logf("get interface %s addrs: %v", iface.Name, err)
			return false
		}
		for _, addr := range addrs {
			// just count colons in a string
			// even with port ipv4 can only have 1
			// ipv6 has to have at least 2
			if strings.Count(addr.String(), ":") >= 2 {
				return true
			}
		}
	}
	return false
}

func DialTestServer(ctx context.Context, logf Logf) bool {
	if logf == nil {
		logf = func(format string, args ...any) {}
	}
	dialer := &net.Dialer{
		Timeout: DefaultTimeout,
	}
	r := &net.Resolver{
		Dial: dialer.DialContext,
	}
	addrs, err := r.LookupIP(ctx, "ip6", TestServer)
	if err != nil {
		logf("resolve ipv6 for %s: %v", TestServer, err)
		return false
	}

	if len(addrs) == 0 {
		logf("no ipv6 address detected for %s", TestServer)
		return false
	}

	a := addrs[0]
	port := 80
	if err := pingConn(ctx, dialer, a, port); err != nil {
		logf("test connection to %s: %v", TestServer, err)
		return false
	}
	return true
}

func pingConn(ctx context.Context, d *net.Dialer, addr net.IP, port int) error {
	addrstr := fmt.Sprintf("[%s]:%d", addr.String(), port)
	c, err := d.DialContext(ctx, "tcp", addrstr)
	if err != nil {
		return fmt.Errorf("dial '%s': %w", addrstr, err)
	}
	return c.Close()
}
