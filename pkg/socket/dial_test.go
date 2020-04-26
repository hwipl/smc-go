package socket

import (
	"log"
	"net"
	"testing"
)

func TestDialHostname(t *testing.T) {
	var want, got string
	var cc, cs net.Conn
	var l net.Listener
	var err error

	// test hostname
	l, err = Listen("localhost:50101")
	if err != nil {
		t.Skip(err)
	}
	cc, err = Dial("localhost:50101")
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

func TestDialIPv6Zone(t *testing.T) {
	var want, got string
	var cc, cs net.Conn
	var l net.Listener
	var err error

	// test ipv6 with zone
	l, err = Listen("[::1%lo]:50102")
	if err != nil {
		t.Skip(err)
	}
	cc, err = Dial("[::1%lo]:50102")
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

func TestDialNoHost(t *testing.T) {
	var want, got string
	var cc, cs net.Conn
	var l net.Listener
	var err error

	// test no host
	l, err = Listen("127.0.0.1:50103")
	if err != nil {
		t.Skip(err)
	}
	cc, err = Dial(":50103")
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
