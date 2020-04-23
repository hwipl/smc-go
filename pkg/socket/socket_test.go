package socket

import (
	"log"
	"testing"
)

func TestListen(t *testing.T) {
	l, err := Listen("127.0.0.1", 0)
	if err != nil {
		log.Fatal(err)
	}
	l.Close()
}
