package socket

import (
	"net"

	"golang.org/x/sys/unix"
)

// SMC protocol definitions
const (
	protoIPv4 = 0
	protoIPv6 = 1
)

// parseIP parses the ip address in the string address and returns an IPv4 and
// an IPv6 address
func parseIP(address string) (net.IP, net.IP) {
	ip := net.ParseIP(address)
	if ip == nil {
		return nil, nil
	}
	return ip.To4(), ip.To16()
}

// createSockAddr constructs a socket address from address and port
func createSockaddr(address string, port int) (typ string, s unix.Sockaddr) {
	ipv4, ipv6 := parseIP(address)
	if ipv4 != nil {
		sockaddr4 := &unix.SockaddrInet4{}
		sockaddr4.Port = port
		copy(sockaddr4.Addr[:], ipv4[:net.IPv4len])
		return "ipv4", sockaddr4
	}
	if ipv6 != nil {
		sockaddr6 := &unix.SockaddrInet6{}
		sockaddr6.Port = port
		copy(sockaddr6.Addr[:], ipv6[:net.IPv6len])
		return "ipv6", sockaddr6
	}

	return "err", nil
}
