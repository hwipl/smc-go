package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"net"

	"github.com/hwipl/smc-go/pkg/socket"
)

// variable definitions
var (
	server  bool   = false
	client  bool   = false
	address string = "127.0.0.1"
	port    int    = 50000
)

// run as a server
func runServer(address string, port int) {
	// listen for smc connections
	l, err := socket.Listen(address, port)
	if err != nil {
		log.Fatal(err)
	}
	defer l.Close()

	// accept new connections from listener and handle them
	for {
		// accept new connection
		conn, err := l.Accept()
		if err != nil {
			log.Fatal(err)
		}

		// handle new connection: send data back to client
		go func(c net.Conn) {
			fmt.Printf("New client connection\n")
			written, err := io.Copy(c, c)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf("Echoed %d bytes to client\n", written)
			c.Close()
		}(conn)
	}
}

// run as a client
func runClient(address string, port int) {
	// connect via smc
	conn, err := socket.Dial(address, port)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	fmt.Printf("Connected to server\n")

	// sent text, read reply an
	text := "Hello, world\n"
	fmt.Fprintf(conn, text)
	fmt.Printf("Sent %d bytes to server: %s", len(text), text)
	reply, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Read %d bytes from server: %s", len(reply), reply)
}

func main() {
	// parse command line arguments
	flag.BoolVar(&server, "s", server, "run server")
	flag.BoolVar(&client, "c", client, "run client")
	flag.StringVar(&address, "a", address, "server/client address")
	flag.IntVar(&port, "p", port, "server/client port")

	flag.Parse()

	// run server
	if server {
		runServer(address, port)
		return
	}

	// run client
	if client {
		runClient(address, port)
		return
	}
}
