package http

import (
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
)

var (
	httpOutput Server
)

// Server is returned by StartHTTPOutput and contains an output buffer for
// the http server and the listener of the http server
type Server struct {
	Buffer   Buffer
	Listener net.Listener
}

// printHTTP prints the content of the http server's Buffer to http clients
func printHTTP(w http.ResponseWriter, r *http.Request) {
	b := httpOutput.Buffer.CopyBuffer()
	if _, err := io.Copy(w, b); err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
	flush := r.URL.Query().Get("flush")
	if flush == "true" {
		httpOutput.Buffer.Reset()
	}
}

// StartHTTPOutput starts a http server that listens on address, and returns
// Server that contains the output Buffer and Listener
func StartHTTPOutput(address string) *Server {
	// create listener
	listener, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatal(err)
	}
	httpOutput.Listener = listener

	// start listening
	http.HandleFunc("/", printHTTP)
	go http.Serve(listener, nil)

	return &httpOutput
}
