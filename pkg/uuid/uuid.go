package uuid

import (
	"crypto/rand"
	"fmt"
	"io"
)

// New generates a cryptographically secure RFC 4122 compliant UUID v4 string.
func New() string {
	var b [16]byte
	_, err := io.ReadFull(rand.Reader, b[:])
	if err != nil {
		return fmt.Sprintf("id-%d", randInt())
	}
	b[6] = (b[6] & 0x0f) | 0x40 // Version 4
	b[8] = (b[8] & 0x3f) | 0x80 // Variant RFC 4122
	return fmt.Sprintf("%08x-%04x-%04x-%04x-%012x",
		b[0:4], b[4:6], b[6:8], b[8:10], b[10:16])
}

func randInt() int64 {
	var b [8]byte
	_, _ = io.ReadFull(rand.Reader, b[:])
	return int64(b[0])<<56 | int64(b[1])<<48 | int64(b[2])<<40 | int64(b[3])<<32 |
		int64(b[4])<<24 | int64(b[5])<<16 | int64(b[6])<<8 | int64(b[7])
}
