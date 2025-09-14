package headers

import (
	"bytes"
	"errors"
	"slices"
	"strings"
	"unicode"
)

const newLine = "\r\n"

var validSpecialCharacters = []rune{'!', '#', '$', '%', '&', '\'', '*', '+', '-', '.', '^', '_', '`', '|', '~'}

type Headers map[string]string

func NewHeaders() Headers {
	return make(Headers)
}

func (h Headers) Parse(data []byte) (n int, done bool, err error) {
	idx := bytes.Index(data, []byte(newLine))
	if idx == -1 {
		return 0, false, nil
	}

	lineStr := string(data[:idx])

	if strings.HasPrefix(lineStr, newLine) || lineStr == "" {
		return idx + len([]byte(newLine)), true, nil
	}

	parts := strings.SplitN(lineStr, ":", 2)

	if len(parts) != 2 {
		return 0, false, errors.New("invalid header line")
	}
	if parts[0] == "" || strings.HasSuffix(parts[0], " ") || strings.HasSuffix(parts[0], "\t") {
		return 0, false, errors.New("invalid header name")
	}

	f := func(r rune) bool {
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) && !slices.Contains(validSpecialCharacters, r) {
			return true
		}
		return false
	}

	if strings.IndexFunc(strings.TrimSpace(parts[0]), f) != -1 {
		return 0, false, errors.New("invalid header name")
	}

	if len(parts) == 2 {
		key := strings.TrimSpace(parts[0])
		key = strings.ToLower(key)
		value := strings.TrimSpace(parts[1])

		if existingValue, exists := h[key]; exists {
			h[key] = existingValue + ", " + value
		} else {
			h[key] = value
		}
	}

	return idx + len([]byte(newLine)), false, nil
}
