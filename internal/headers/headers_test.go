package headers

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_HeadersParse_CorrectParse(t *testing.T) {
	// Test: Valid single header
	headers := NewHeaders()
	data := []byte("Host: localhost:42069\r\n\r\n")
	n, done, err := headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, "localhost:42069", headers["host"])
	assert.Equal(t, len(data)-2, n)
	assert.False(t, done)

	// Test: Valid single header with extra white space
	headers = NewHeaders()
	data = []byte("		Host:  localhost:42069		  \r\n\r\n")
	n, done, err = headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, "localhost:42069", headers["host"])
	assert.Equal(t, len(data)-2, n)
	assert.False(t, done)
}

func Test_HeadersParse_CorrectParse_MultipleHeaders(t *testing.T) {
	// Test: Valid multiple headers
	headers := NewHeaders()
	data := []byte("Host: localhost:42069\r\n")
	n, done, err := headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, "localhost:42069", headers["host"])
	assert.Equal(t, len(headers), 1)
	assert.Equal(t, len(data), n)
	assert.False(t, done)

	data = []byte("User-Agent: curl/7.81.0\r\n")
	n, done, err = headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, "curl/7.81.0", headers["user-agent"])
	assert.Equal(t, len(headers), 2)
	assert.Equal(t, len(data), n)
	assert.False(t, done)
}

func Test_HeadersParse_CorrectParse_Done(t *testing.T) {
	// Test: Valid header parse with done
	headers := NewHeaders()
	data := []byte("\r\n")
	n, done, err := headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, 2, n)
	assert.True(t, done)
}

func Test_HeadersParse_CorrectParse_MultipleValues(t *testing.T) {
	headers := NewHeaders()
	data := []byte("Accept: text/plain\r\n")
	n, done, err := headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, "text/plain", headers["accept"])
	assert.Equal(t, len(headers), 1)
	assert.Equal(t, len(data), n)
	assert.False(t, done)

	data = []byte("Accept: application/json\r\n")
	n, done, err = headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, "text/plain, application/json", headers["accept"])
	assert.Equal(t, len(headers), 1)
	assert.Equal(t, len(data), n)
	assert.False(t, done)
}

func Test_HeadersParse_IncorrectParse(t *testing.T) {
	headers := NewHeaders()
	data := []byte("     Host : localhost:42069      \r\n\r\n")
	n, done, err := headers.Parse(data)
	require.Error(t, err)
	assert.Equal(t, 0, n)
	assert.False(t, done)
}
