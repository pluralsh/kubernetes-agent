package ioz

import (
	"encoding/base64"
	"fmt"
	"io"
	"os"
)

func LoadBase64Secret(filename string) ([]byte, error) {
	encodedAuthSecret, err := os.ReadFile(filename) // nolint: gosec
	if err != nil {
		return nil, fmt.Errorf("read file: %w", err)
	}
	decodedAuthSecret := make([]byte, len(encodedAuthSecret))

	n, err := base64.StdEncoding.Decode(decodedAuthSecret, encodedAuthSecret)
	if err != nil {
		return nil, fmt.Errorf("decoding: %w", err)
	}
	return decodedAuthSecret[:n], nil
}

// NewReceiveReader turns receiver into an io.Reader. Errors from the receiver
// function are passed on unmodified. This means receiver should emit
// io.EOF when done.
func NewReceiveReader(receiver func() ([]byte, error)) io.Reader {
	return &receiveReader{receiver: receiver}
}

type receiveReader struct {
	receiver func() ([]byte, error)
	data     []byte
	err      error
}

func (rr *receiveReader) Read(p []byte) (int, error) {
	if len(rr.data) == 0 && rr.err == nil {
		rr.data, rr.err = rr.receiver()
	}

	n := copy(p, rr.data)
	rr.data = rr.data[n:]

	// We want to return any potential error only in case we have no
	// buffered data left. Otherwise, it can happen that we do not relay
	// bytes when the reader returns both data and an error.
	if len(rr.data) == 0 {
		return n, rr.err
	}

	return n, nil
}
