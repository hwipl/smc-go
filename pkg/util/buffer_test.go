package util

import (
	"bytes"
	"testing"
)

func TestBuffer(t *testing.T) {
	var buf buffer
	var want []byte
	var got []byte

	// test empty buffer
	got = buf.copyBuffer().Bytes()
	if !bytes.Equal(want, got) {
		t.Errorf("buf = %s; want %s", got, want)
	}

	// test writing to empty buffer
	want = []byte("hello world")
	buf.Write(want)
	got = buf.copyBuffer().Bytes()
	if !bytes.Equal(want, got) {
		t.Errorf("buf = %s; want %s", got, want)
	}

	// test appending to buffer
	buf.Write(want)
	want = []byte("hello worldhello world")
	got = buf.copyBuffer().Bytes()
	if !bytes.Equal(want, got) {
		t.Errorf("buf = %s; want %s", got, want)
	}

	// test resetting buffer
	buf.reset()
	want = []byte("")
	got = buf.copyBuffer().Bytes()
	if !bytes.Equal(want, got) {
		t.Errorf("buf = %s; want %s", got, want)
	}
}
