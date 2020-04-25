package socket

import (
	"context"
	"log"
	"net"
	"strconv"

	"golang.org/x/sys/unix"
)

// SMC protocol definitions
const (
	protoIPv4 = 0
	protoIPv6 = 1
)

// parseAddress parses a host:port address string and returns the host as
// string and the port as int
func parseAddress(address string) (string, int) {
	host, p, err := net.SplitHostPort(address)
	if err != nil {
		log.Fatal(err)
	}
	port, err := strconv.Atoi(p)
	if err != nil {
		log.Fatal(err)
	}
	return host, port
}

// parseIP parses and resolves the ip address in the host string and returns an
// IPv4 or IPv6 address
func parseIP(address string) *net.IPAddr {
	ipaddrs, err := net.DefaultResolver.LookupIPAddr(context.Background(),
		address)
	if err != nil {
		return nil
	}
	return &ipaddrs[0]
}

// createSockAddr constructs a socket address from address
func createSockaddr(address string) (typ string, s unix.Sockaddr) {
	host, port := parseAddress(address)
	ipaddr := parseIP(host)
	if ipaddr == nil {
		return "err", nil
	}
	ipv4 := ipaddr.IP.To4()
	ipv6 := ipaddr.IP.To16()
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
