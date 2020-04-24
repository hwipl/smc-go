package socket

import (
	"fmt"
	"net"
	"os"

	"golang.org/x/sys/unix"
)

// Dial creates a SMC connection to address and port
func Dial(address string) (net.Conn, error) {
	var conn net.Conn
	var err error
	var fd int

	// construct socket address from address
	typ, sockaddr := createSockaddr(address)
	if typ == "err" {
		return conn, fmt.Errorf("Error parsing IP")
	}

	// create socket
	if typ == "ipv4" {
		fd, err = unix.Socket(unix.AF_SMC, unix.SOCK_STREAM,
			protoIPv4)
	} else {
		fd, err = unix.Socket(unix.AF_SMC, unix.SOCK_STREAM,
			protoIPv6)
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
