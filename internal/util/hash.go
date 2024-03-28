package util

import (
	"crypto/md5"
	"encoding/base64"
	"math/rand"
)

// HashMD5 hashes the input byte slice using MD5 algorithm and returns the first 12 characters of the base64 encoded hash.
func HashMD5(data []byte) []byte {
	hash := md5.New()
	hash.Write(data)
	hashSum := hash.Sum(nil)
	// Encode the hash sum to base64 and return the first 12 characters
	return []byte(base64.StdEncoding.EncodeToString(hashSum))[:12]
}

// Randomize randomizes the input byte slice using XOR operations, shuffling, and resizing, and then hashes the result using the HashMD5 function.
func Randomize(data []byte, offset int) []byte {
	// XOR each byte with the offset - 1
	for i := range data {
		data[i] = byte(ClosestCharCode((int(data[i]) ^ offset - 1) & 0xFF))
	}

	// Shuffle the data
	Shuffle(data, offset)

	// Resize the data to a fixed size
	Resize(&data, 15)

	// Hash the randomized data and return the result
	return HashMD5(data)
}

// Shuffle shuffles the elements of the input byte slice using Fisher-Yates algorithm.
func Shuffle(data []byte, offset int) {
	r := rand.New(rand.NewSource(int64(offset)))
	r.Shuffle(len(data), func(i, j int) {
		data[i], data[j] = data[j], data[i]
	})
}

// Resize resizes the input byte slice to the specified size by randomly selecting elements.
func Resize(data *[]byte, size int) {
	source := rand.New(rand.NewSource(int64(len(*data) + size%17)))

	// Limit the size to avoid resizing beyond half of the original length
	if size >= len(*data)/2 {
		size = len(*data) / 2
	}

	// Generate random positions and resize the data
	positions := source.Perm(size)
	*data = (*data)[:size]
	for i, pos := range positions {
		(*data)[i] = (*data)[pos]
	}
}

// ClosestCharCode returns the closest valid ASCII character code for the given input code within the range [45, 58] or [97, 122].
func ClosestCharCode(code int) int {
	// Ensure the code is within the valid ASCII range
	if code >= 45 && code <= 58 || code >= 97 && code <= 122 {
		return code
	}

	// Adjust the code to the nearest valid ASCII value
	return int(min(max(float64(code), 45), 122))
}
