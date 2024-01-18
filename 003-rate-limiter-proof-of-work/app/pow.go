package app

import (
	"bytes"
	"crypto/sha256"
	"errors"
	"math/bits"
	"proof-of-work/types"

	"github.com/google/uuid"
)

func getNonce() []byte {
	return []byte(uuid.NewString())
}

func getChecksum(nonce []byte, secret string) []byte {
	merge := append(nonce, []byte(secret)...)
	h := sha256.Sum256(merge)
	return h[:]
}

func verify(pi types.ParsedIssue) (bool, error) {
	pass := verifyDifficulty(pi.Hash, difficulty) // wrong algo, check later
	if !pass {
		return pass, errors.New("failed verify Difficulty")
	}

	pass = verifyNonce(pi.Nonce, pi.Checksum)
	if !pass {
		return pass, errors.New("failed verify Nonce")
	}

	pass = verifyHash(pi.Nonce, pi.Counter, pi.Hash)
	if !pass {
		return pass, errors.New("failed verify Hash")
	}
	return pass, nil
}

func verifyDifficulty(hash []byte, difficulty int) bool {
	diff := difficulty
	for _, b := range hash {
		lead := bits.LeadingZeros8(uint8(b))
		diff -= lead

		if lead < 8 {
			return diff <= 0
		}

		if diff <= 0 {
			return true
		}
	}
	return false
}

func verifyNonce(nonce, checksum []byte) bool {
	h := getChecksum(nonce, secret)
	return bytes.Equal(h, checksum)
}

func verifyHash(nonce, counter, hash []byte) bool {
	h := getChecksum(counter, string(nonce))
	return bytes.Equal(h, hash)
}
