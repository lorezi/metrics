package main

import (
	"crypto/sha256"
	"encoding/hex"
)

// computeSum generates a SHA256 hash of the body
func computeSum(body []byte) []byte {
	h := sha256.New()
	h.Write(body)
	hashed := hex.EncodeToString(h.Sum(nil))
	return []byte(hashed)
}

func main() {

}
