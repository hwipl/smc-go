package http

import (
	"bytes"
	"testing"
)

func TestBuffer(t *testing.T) {
	var buf Buffer
	var want []byte
	var got []byte

	// test empty buffer
	got = buf.CopyBuffer().Bytes()
	if !bytes.Equal(want, got) {
		t.Errorf("buf = %s; want %s", got, want)
	}

	// test writing to empty buffer
	want = []byte("hello world")
	buf.Write(want)
	got = buf.CopyBuffer().Bytes()
	if !bytes.Equal(want, got) {
		t.Errorf("buf = %s; want %s", got, want)
	}

	// test appending to buffer
	buf.Write(want)
	want = []byte("hello worldhello world")
	got = buf.CopyBuffer().Bytes()
	if !bytes.Equal(want, got) {
		t.Errorf("buf = %s; want %s", got, want)
	}

	// test resetting buffer
	buf.Reset()
	want = []byte("")
	got = buf.CopyBuffer().Bytes()
	if !bytes.Equal(want, got) {
		t.Errorf("buf = %s; want %s", got, want)
	}
}
