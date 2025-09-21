package response

import (
	"fmt"
	"io"

	"http.server/internal/headers"
)

type StatusCode int

const (
	OK                  StatusCode = 200
	BadRequest          StatusCode = 400
	NotFound            StatusCode = 404
	InternalServerError StatusCode = 500
)

func WriteStatusLine(w io.Writer, statusCode StatusCode) string {
	switch statusCode {
	case OK:
		statusLine := "HTTP/1.1 200 OK\r\n"
		w.Write([]byte(statusLine))
		return statusLine
	case BadRequest:
		statusLine := "HTTP/1.1 400 Bad Request\r\n"
		w.Write([]byte(statusLine))
		return statusLine
	case NotFound:
		statusLine := "HTTP/1.1 404 Not Found\r\n"
		w.Write([]byte(statusLine))
		return statusLine
	case InternalServerError:
		statusLine := "HTTP/1.1 500 Internal Server Error\r\n"
		w.Write([]byte(statusLine))
		return statusLine
	default:
		statusLine := fmt.Sprintf("HTTP/1.1 %d \r\n", statusCode)
		w.Write([]byte(statusLine))
		return statusLine
	}
}

func GetDefaultHeaders(contentLength int) headers.Headers {
	h := headers.NewHeaders()
	h["Content-Length"] = fmt.Sprintf("%d", contentLength)
	h["Content-Type"] = "text/plain"
	h["Connection"] = "close"
	return h
}

func WriteHeaders(w io.Writer, headers headers.Headers) error {
	for key, value := range headers {
		_, err := fmt.Fprintf(w, "%s: %s\r\n", key, value)
		if err != nil {
			return err
		}
	}
	_, err := w.Write([]byte("\r\n"))
	return err
}

func WriteBody(w io.Writer, body []byte) (int, error) {
	return w.Write(body)
}
