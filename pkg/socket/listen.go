package socket

import (
	"fmt"
	"net"
	"os"

	"golang.org/x/sys/unix"
)

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
			protoIPv4)
	} else {
		fd, err = unix.Socket(unix.AF_SMC, unix.SOCK_STREAM,
			protoIPv6)
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
