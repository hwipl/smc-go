package socket

import (
	"log"
	"net"
	"testing"
)

func TestDial(t *testing.T) {
	var want, got string
	var cc, cs net.Conn
	var l net.Listener
	var err error

	// test ipv4
	l, err = Listen("127.0.0.1:0")
	if err != nil {
		log.Fatal(err)
	}
	cc, err = Dial(l.Addr().String())
	if err != nil {
		log.Fatal(err)
	}
	cs, err = l.Accept()
	if err != nil {
		log.Fatal(err)
	}
	want = l.Addr().String()
	got = cc.RemoteAddr().String()
	if got != want {
		t.Errorf("RemoteAddr() = %s; want %s", got, want)
	}
	cc.Close()
	cs.Close()
	l.Close()

	// test ipv6
	l, err = Listen("[::1]:0")
	if err != nil {
		log.Fatal(err)
	}
	cc, err = Dial(l.Addr().String())
	if err != nil {
		log.Fatal(err)
	}
	cs, err = l.Accept()
	if err != nil {
		log.Fatal(err)
	}
	want = l.Addr().String()
	got = cc.RemoteAddr().String()
	if got != want {
		t.Errorf("RemoteAddr() = %s; want %s", got, want)
	}
	cc.Close()
	cs.Close()
	l.Close()
}
