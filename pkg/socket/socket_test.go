package socket

import (
	"log"
	"net"
	"testing"
)

func TestListen(t *testing.T) {
	var want, got string
	var l net.Listener
	var err error

	// test specific ip, specific port, ipv4
	l, err = Listen("127.0.0.1", 50000)
	if err != nil {
		log.Fatal(err)
	}
	want = "127.0.0.1:50000"
	got = l.Addr().String()
	if got != want {
		t.Errorf("Addr() = %s; want %s", got, want)
	}
	l.Close()

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

	// test all ips, specific port, ipv4
	l, err = Listen("0.0.0.0", 50000)
	if err != nil {
		log.Fatal(err)
	}
	want = "0.0.0.0:50000"
	got = l.Addr().String()
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

	// test specific ip, specific port, ipv6
	l, err = Listen("::1", 50000)
	if err != nil {
		log.Fatal(err)
	}
	want = "[::1]:50000"
	got = l.Addr().String()
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

	// test all ips, specific port, ipv6
	l, err = Listen("::", 50000)
	if err != nil {
		log.Fatal(err)
	}
	want = "[::]:50000"
	got = l.Addr().String()
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
