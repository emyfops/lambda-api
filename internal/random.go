package internal

import (
	"math/rand"
	"unsafe"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const (
	letterIdxBits = 6                    // log2(52), six bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // Number of letter indices fitting in 63 bits
)

// RandBytesMaskSrc generates a random string of n characters.
// It is not cryptographically secure.
func RandBytesMaskSrc(n int) []byte {
	b := make([]byte, n)

	// The index is n-1 because we are 0-indexed.
	// The cache is the current integer we are generating characters from.
	for i, cache, remain :=
		n-1, rand.Int63(), letterIdxMax; i >= 0; {

		// If the remain is 0, we generate a new random int 63 bits
		// and start over.
		if remain == 0 {
			cache, remain = rand.Int63(), letterIdxMax
		}

		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}

		// With 63 bits, we can generate 10 characters because 63/6 = 10.5.
		// We shift the cache 6 bits to the right to get the next character.
		// We also decrement the remains because we generated a character.
		// This is the most clever random string generation
		// algorithm I've ever seen.
		cache >>= letterIdxBits
		remain--
	}

	return b
}

func RandString(n int) string {
	b := RandBytesMaskSrc(n)

	// The length of b is n, so we can safely cast it to a string
	// and avoid the overhead of a copy.
	return *(*string)(unsafe.Pointer(&b))
}
