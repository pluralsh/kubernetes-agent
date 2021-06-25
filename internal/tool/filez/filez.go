package filez

import (
	"encoding/base64"
	"fmt"
	"os"
)

func LoadBase64Secret(filename string) ([]byte, error) {
	encodedAuthSecret, err := os.ReadFile(filename) // nolint: gosec
	if err != nil {
		return nil, fmt.Errorf("read file: %v", err)
	}
	decodedAuthSecret := make([]byte, len(encodedAuthSecret))

	n, err := base64.StdEncoding.Decode(decodedAuthSecret, encodedAuthSecret)
	if err != nil {
		return nil, fmt.Errorf("decoding: %v", err)
	}
	return decodedAuthSecret[:n], nil
}
