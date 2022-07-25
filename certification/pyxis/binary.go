package pyxis

import (
	"crypto/sha256"
	"fmt"
	"io"
)

func BinarySHA(r io.Reader) (string, error) {
	exe, err := io.ReadAll(r)
	if err != nil {
		return "", fmt.Errorf("could not read all bytes: %v", err)
	}
	sha := sha256.Sum256(exe)
	return fmt.Sprintf("%x", sha), nil
}
