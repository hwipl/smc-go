package socket

import (
	"fmt"
	"net"
	"os"

	"golang.org/x/sys/unix"
)

// SMC definitions
const (
	SMCProtoIPv4 = 0
	SMCProtoIPv6 = 1
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

// Listen creates a SMC listener that listens on address and port
func Listen(address string, port int) (net.Listener, error) {
	var l net.Listener
	var err error
	var fd int

	// construct socket address from address and port
	typ, sockaddr := createSockaddr(address, port)
	if typ == "err" {
		return l, fmt.Errorf("Error parsing IP")
	}

	// create socket
	if typ == "ipv4" {
		fd, err = unix.Socket(unix.AF_SMC, unix.SOCK_STREAM,
			SMCProtoIPv4)
	} else {
		fd, err = unix.Socket(unix.AF_SMC, unix.SOCK_STREAM,
			SMCProtoIPv6)
	}
	if err != nil {
		return l, err
	}
	defer unix.Close(fd)

	// bind socket address
	err = unix.Bind(fd, sockaddr)
	if err != nil {
		return l, err
	}

	// start listening
	err = unix.Listen(fd, 1)
	if err != nil {
		return l, err
	}

	// create a listener from listening socket
	file := os.NewFile(uintptr(fd), "")
	l, err = net.FileListener(file)
	return l, err
}

// smcDial is a SMC version of Dial()
func smcDial(address string, port int) (net.Conn, error) {
	var conn net.Conn
	var err error
	var fd int

	// construct socket address from address and port
	typ, sockaddr := createSockaddr(address, port)
	if typ == "err" {
		return conn, fmt.Errorf("Error parsing IP")
	}

	// create socket
	if typ == "ipv4" {
		fd, err = unix.Socket(unix.AF_SMC, unix.SOCK_STREAM,
			SMCProtoIPv4)
	} else {
		fd, err = unix.Socket(unix.AF_SMC, unix.SOCK_STREAM,
			SMCProtoIPv6)
	}

	if err != nil {
		return conn, err
	}
	defer unix.Close(fd)

	// connect to server
	err = unix.Connect(fd, sockaddr)
	if err != nil {
		return conn, err
	}

	// create a connection from connected socket
	file := os.NewFile(uintptr(fd), "")
	conn, err = net.FileConn(file)
	return conn, err
}
