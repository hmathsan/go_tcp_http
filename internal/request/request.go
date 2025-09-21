package request

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"

	"http.server/internal/headers"
	"http.server/internal/utils"
)

type ParseState int

const (
	Initialized ParseState = iota
	ParsingHeaders
	ParsingBody
	ParsingBodyValidation
	Done
)

type Request struct {
	RequestLine RequestLine
	state       ParseState
	Headers     headers.Headers
	Body        []byte
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

const bufferSize = 8
const newLine = "\r\n"

func RequestFromReader(reader io.Reader) (*Request, error) {
	buf := make([]byte, bufferSize)
	readToIndex := 0

	request := &Request{
		state:   Initialized,
		Headers: headers.NewHeaders(),
		Body:    []byte{},
	}

	for request.state != Done {
		if readToIndex >= len(buf) {
			newBuf := make([]byte, len(buf)*2)
			copy(newBuf, buf)
			buf = newBuf
		}

		numBytesRead, err := reader.Read(buf[readToIndex:])
		if err != nil {
			if errors.Is(err, io.EOF) {
				request.state = ParsingBodyValidation
			} else {
				return nil, err
			}
		}
		readToIndex += numBytesRead

		numBytesParsed, err := request.parse(buf[:readToIndex])
		if err != nil {
			return nil, err
		}

		copy(buf, buf[numBytesParsed:])
		readToIndex -= numBytesParsed
	}

	return request, nil
}

func parseRequestLine(data []byte) (*RequestLine, int, error) {
	idx := bytes.Index(data, []byte(newLine))
	if idx == -1 {
		return nil, 0, nil
	}

	lineStr := string(data[:idx])
	lineParts := strings.Split(lineStr, " ")

	requestLine := RequestLine{}

	if len(lineParts) != 3 {
		return nil, 0, errors.New("invalid request line")
	}

	if !utils.IsUpper(lineParts[0]) {
		return nil, 0, errors.New("invalid method")
	}
	requestLine.Method = lineParts[0]

	if lineParts[1][0] != '/' {
		return nil, 0, errors.New("invalid request path")
	}
	requestLine.RequestTarget = lineParts[1]

	httpVersion := strings.Split(lineParts[2], "/")
	if len(httpVersion) != 2 || httpVersion[0] != "HTTP" || httpVersion[1] != "1.1" {
		return nil, 0, errors.New("invalid HTTP version")
	}
	requestLine.HttpVersion = httpVersion[1]

	return &requestLine, idx + 2, nil
}

func (r *Request) parse(data []byte) (int, error) {
	switch r.state {
	case Initialized:
		requestLine, n, err := parseRequestLine(data)
		if err != nil {
			return 0, err
		}
		if n == 0 {
			return 0, nil
		}
		r.RequestLine = *requestLine
		r.state = ParsingHeaders
		return n, nil
	case ParsingHeaders:
		n, done, err := r.Headers.Parse(data)
		if err != nil {
			return 0, err
		}
		if n == 0 {
			return 0, nil
		}
		if done {
			r.state = ParsingBody
		}
		totalBytesParsed := n

		return totalBytesParsed, nil
	case ParsingBody:
		if len(r.Body) == 0 {
			contentLength := r.Headers.Get("content-length")
			if (contentLength == "") || (contentLength == "0") {
				r.state = Done
				return 0, nil
			}
			contentLen, err := strconv.Atoi(contentLength)
			if err != nil || contentLen < 0 {
				return 0, fmt.Errorf("invalid content-length: %v", contentLength)
			}
		}

		r.Body = append(r.Body, data...)

		return len(data), nil
	case ParsingBodyValidation:
		contentLength := r.Headers.Get("content-length")
		if (contentLength == "") || (contentLength == "0") {
			r.state = Done
			return 0, nil
		}

		contentLen, err := strconv.Atoi(contentLength)
		if err != nil {
			return 0, fmt.Errorf("invalid content-length: %v", err)
		}
		if len(r.Body) > contentLen {
			return 0, fmt.Errorf("body length exceeds Content-Length")
		}
		if len(r.Body) < contentLen {
			return 0, fmt.Errorf("body length less than Content-Length")
		}
		if len(r.Body) == contentLen {
			r.state = Done
		}

		return 0, nil
	case Done:
		return 0, nil
	default:
		return 0, fmt.Errorf("error: unknown state %d", r.state)
	}
}
