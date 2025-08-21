package main

import (
	"fmt"
	"log"
	"net"

	"http.server/internal/request"
)

func main() {
	l, err := net.Listen("tcp", ":42069")
	if err != nil {
		log.Fatal(err)
	}
	defer l.Close()

	for {
		conn, err := l.Accept()
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("Connection accepted")
		request, err := request.RequestFromReader(conn)
		if err != nil {
			log.Fatal("Failed to parse request:", err)
		}
		fmt.Printf("Request line:\n- Method: %s\n- Target: %s\n- Version: %s\n",
			request.RequestLine.Method,
			request.RequestLine.RequestTarget,
			request.RequestLine.HttpVersion,
		)
	}
}
