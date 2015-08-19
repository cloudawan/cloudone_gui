package random

import (
	"crypto/rand"
	"fmt"
	"io"
)

func UUID() string {
	// RFC4122
	uuid := make([]byte, 16)
	io.ReadFull(rand.Reader, uuid)
	// Section 4.1.1 variant bits
	uuid[8] = uuid[8]&^0xc0 | 0x80
	// Section 4.1.3 version 4 (pseudo-random);
	uuid[6] = uuid[6]&^0xf0 | 0x40
	return fmt.Sprintf("%x-%x-%x-%x-%x", uuid[0:4], uuid[4:6], uuid[6:8], uuid[8:10], uuid[10:])
}
