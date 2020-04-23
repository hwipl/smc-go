package socket

import (
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"
	"testing"
)

// Specific port tests could fail because address/port is already in use.
// Run them in individual test functions and skip in case of listen error.

func TestListenSpecificIPv4(t *testing.T) {
	var want, got string
	var l net.Listener
	var err error

	// test specific ip, specific port, ipv4
	l, err = Listen("127.0.0.1", 50001)
	if err != nil {
		t.Skip(err)
	}
	want = "127.0.0.1:50001"
	got = l.Addr().String()
	if got != want {
		t.Errorf("Addr() = %s; want %s", got, want)
	}
	l.Close()
}

func TestListenSpecificAllIPv4(t *testing.T) {
	var want, got string
	var l net.Listener
	var err error

	// test all ips, specific port, ipv4
	l, err = Listen("0.0.0.0", 50002)
	if err != nil {
		t.Skip(err)
	}
	want = "0.0.0.0:50002"
	got = l.Addr().String()
	if got != want {
		t.Errorf("Addr() = %s; want %s", got, want)
	}
	l.Close()
}

func TestListenSpecificIPv6(t *testing.T) {
	var want, got string
	var l net.Listener
	var err error

	// test specific ip, specific port, ipv6
	l, err = Listen("::1", 50003)
	if err != nil {
		t.Skip(err)
	}
	want = "[::1]:50003"
	got = l.Addr().String()
	if got != want {
		t.Errorf("Addr() = %s; want %s", got, want)
	}
	l.Close()
}

func TestListenSpecificAllIPv6(t *testing.T) {
	var want, got string
	var l net.Listener
	var err error

	// test all ips, specific port, ipv6
	l, err = Listen("::", 50004)
	if err != nil {
		t.Skip(err)
	}
	want = "[::]:50004"
	got = l.Addr().String()
	if got != want {
		t.Errorf("Addr() = %s; want %s", got, want)
	}
	l.Close()
}

func TestListen(t *testing.T) {
	var want, got string
	var l net.Listener
	var err error

	// test specific ip, random port, ipv4
	l, err = Listen("127.0.0.1", 0)
	if err != nil {
		log.Fatal(err)
	}
	want = "127.0.0.1:"          // ignore port
	got = l.Addr().String()[:10] // ignore port
	if got != want {
		t.Errorf("Addr() = %s; want %s", got, want)
	}
	l.Close()

	// test all ips, random port, ipv4
	l, err = Listen("0.0.0.0", 0)
	if err != nil {
		log.Fatal(err)
	}
	want = "0.0.0.0:"           // ignore port
	got = l.Addr().String()[:8] // ignore port
	if got != want {
		t.Errorf("Addr() = %s; want %s", got, want)
	}
	l.Close()

	// test specific ip, random port, ipv6
	l, err = Listen("::1", 0)
	if err != nil {
		log.Fatal(err)
	}
	want = "[::1]:"             // ignore port
	got = l.Addr().String()[:6] // ignore port
	if got != want {
		t.Errorf("Addr() = %s; want %s", got, want)
	}
	l.Close()

	// test all ips, random port, ipv6
	l, err = Listen("::", 0)
	if err != nil {
		log.Fatal(err)
	}
	want = "[::]:"              // ignore port
	got = l.Addr().String()[:5] // ignore port
	if got != want {
		t.Errorf("Addr() = %s; want %s", got, want)
	}
	l.Close()
}

func TestDial(t *testing.T) {
	var want, got string
	var cc, cs net.Conn
	var l net.Listener
	var err error
	var port int

	// test ipv4
	l, err = Listen("127.0.0.1", 0)
	if err != nil {
		log.Fatal(err)
	}
	port, err = strconv.Atoi(strings.Split(l.Addr().String(), ":")[1])
	if err != nil {
		log.Fatal(err)
	}
	cc, err = Dial("127.0.0.1", port)
	if err != nil {
		log.Fatal(err)
	}
	cs, err = l.Accept()
	if err != nil {
		log.Fatal(err)
	}
	want = fmt.Sprintf("127.0.0.1:%d", port)
	got = cc.RemoteAddr().String()
	if got != want {
		t.Errorf("RemoteAddr() = %s; want %s", got, want)
	}
	cc.Close()
	cs.Close()
	l.Close()
}
