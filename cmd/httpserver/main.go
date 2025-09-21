package main

import (
	"io"
	"log"
	"os"
	"os/signal"
	"syscall"

	"http.server/internal/request"
	"http.server/internal/server"
)

const port = 42069

func main() {
	server, err := server.Serve(port, func(w io.Writer, req *request.Request) *server.HandlerError {
		if req.RequestLine.RequestTarget == "/yourproblem" {
			log.Println("Client says they have a problem")
			w.Write([]byte("Your problem is not my problem\n"))
			return &server.HandlerError{StatusCode: 400}
		}
		if req.RequestLine.RequestTarget == "/myproblem" {
			log.Println("Client says I have a problem")
			w.Write([]byte("Woopsie, my bad\n"))
			return &server.HandlerError{StatusCode: 500}
		}

		log.Println("Client says all is good")
		w.Write([]byte("All good, frfr\n"))
		return nil
	})
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
	defer server.Close()
	log.Println("Server started on port", port)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Println("Server gracefully shutting down")
}
