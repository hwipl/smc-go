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
	server Server
)

// Server is returned by StartHTTPOutput and contains an output buffer for
// the http server and the listener of the http server
type Server struct {
	Buffer   Buffer
	Listener net.Listener
}

// handleRequest prints the content of the http server's Buffer to http clients
func handleRequest(w http.ResponseWriter, r *http.Request) {
	b := server.Buffer.CopyBuffer()
	if _, err := io.Copy(w, b); err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
	flush := r.URL.Query().Get("flush")
	if flush == "true" {
		server.Buffer.Reset()
	}
}

// StartServer starts a http server that listens on address, and returns
// Server that contains the output Buffer and Listener
func StartServer(address string) *Server {
	// create listener
	listener, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatal(err)
	}
	server.Listener = listener

	// start listening
	http.HandleFunc("/", handleRequest)
	go http.Serve(listener, nil)

	return &server
}
