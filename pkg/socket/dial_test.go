package socket

import (
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"
	"testing"
)

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

	// test ipv6
	l, err = Listen("::1", 0)
	if err != nil {
		log.Fatal(err)
	}
	port, err = strconv.Atoi(strings.Split(l.Addr().String(), "]:")[1])
	if err != nil {
		log.Fatal(err)
	}
	cc, err = Dial("::1", port)
	if err != nil {
		log.Fatal(err)
	}
	cs, err = l.Accept()
	if err != nil {
		log.Fatal(err)
	}
	want = fmt.Sprintf("[::1]:%d", port)
	got = cc.RemoteAddr().String()
	if got != want {
		t.Errorf("RemoteAddr() = %s; want %s", got, want)
	}
	cc.Close()
	cs.Close()
	l.Close()
}
