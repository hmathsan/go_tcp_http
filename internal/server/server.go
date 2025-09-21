package server

import (
	"bytes"
	"io"
	"log"
	"net"
	"sync/atomic"

	"http.server/internal/headers"
	"http.server/internal/request"
	"http.server/internal/response"
)

type Server struct {
	listener net.Listener
	isClosed *atomic.Bool
	handler  Handler
}

type HandlerError struct {
	StatusCode response.StatusCode
}

type Handler func(w io.Writer, req *request.Request) *HandlerError

func Serve(port int, handler Handler) (*Server, error) {
	l, err := net.Listen("tcp", ":42069")
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	isClosed := &atomic.Bool{}
	isClosed.Store(false)
	server := &Server{listener: l, isClosed: isClosed, handler: handler}

	go server.listen()

	return server, nil
}

func (s *Server) Close() error {
	s.isClosed.Store(true)
	return s.listener.Close()
}

func (s *Server) listen() {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			if s.isClosed.Load() {
				return
			}
			log.Printf("Error accepting connection: %v", err)
			continue
		}

		log.Println("Connection accepted")

		go s.handleConnection(conn)
	}
}

func (s *Server) handleConnection(conn net.Conn) {
	defer conn.Close()
	log.Println("Connection accepted from", conn.RemoteAddr())
	req, err := request.RequestFromReader(conn)
	if err != nil {
		log.Println("Error reading request:", err)
		return
	}
	log.Println("Request processed successfully from", conn.RemoteAddr())

	responseBuffer := bytes.Buffer{}
	log.Println("Invoking handler for request from", conn.RemoteAddr())
	handlerError := s.handler(&responseBuffer, req)

	responseStatusCode := response.OK

	if handlerError != nil {
		log.Println("Handler error:", handlerError.StatusCode)
		responseStatusCode = handlerError.StatusCode
	}

	headers := headers.NewHeaders()

	response.WriteStatusLine(conn, responseStatusCode)
	response.WriteHeaders(conn, headers)
	response.WriteBody(conn, responseBuffer.Bytes())

	conn.Close()
}
