package xnet

import (
	"fmt"
	"net"
	"net/netip"
)

func GetRandomListenAddr(addrstr string) (res netip.AddrPort, err error) {
	addr, err := netip.ParseAddr(addrstr)
	if err != nil {
		return res, err
	}

	a, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:0", addr))
	if err != nil {
		return res, err
	}

	l, err := net.ListenTCP("tcp", a)
	if err != nil {
		return res, err
	}

	if err := l.Close(); err != nil {
		return res, err
	}

	port := l.Addr().(*net.TCPAddr).Port

	return netip.AddrPortFrom(addr, uint16(port)), nil
}
