package socket

import (
	"fmt"
	"log"
	"net"
	"os"
	"strconv"

	"golang.org/x/sys/unix"
)

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

// Listen creates a SMC listener that listens on address
func Listen(address string) (net.Listener, error) {
	var l net.Listener
	var err error
	var fd int

	// parse address
	host, port := parseAddress(address)

	// construct socket address from address and port
	typ, sockaddr := createSockaddr(host, port)
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
