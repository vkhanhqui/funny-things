package app

import (
	"encoding/base64"
	"encoding/hex"
	"errors"
	"strings"
)

type Issue struct {
	Nonce      string `json:"nonce"`
	Difficulty int    `json:"difficulty"`
}

type VerifyParser struct {
	Nonce, Hash, Bin []byte
	Counter          int
	Difficulty       int
}

type VerifyIssueReq struct {
	Nonce      string `json:"nonce"`
	Counter    int    `json:"counter"`
	Hash       string `json:"hash"`
	Difficulty int    `json:"difficulty"`
}

func (v *VerifyIssueReq) toParser() (VerifyParser, error) {
	nonceBts, err := base64.StdEncoding.DecodeString(v.Nonce)
	if err != nil {
		return VerifyParser{}, errors.New("invalid nonce")
	}

	splitted := strings.Split(v.Hash, "#")
	if len(splitted) != 2 {
		return VerifyParser{}, errors.New("invalid hash")
	}

	hash, bin := splitted[0], splitted[1]
	hashBts, err := hex.DecodeString(hash)
	if err != nil {
		return VerifyParser{}, errors.New("invalid hash")
	}

	binBts, err := base64.StdEncoding.DecodeString(bin)
	if err != nil {
		return VerifyParser{}, errors.New("invalid hash")
	}
	return VerifyParser{Nonce: nonceBts, Hash: hashBts, Bin: binBts, Counter: v.Counter, Difficulty: v.Difficulty}, nil
}
