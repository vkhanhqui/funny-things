package app

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"strconv"

	"github.com/google/uuid"
)

const (
	difficulty = 10
)

type PowService interface {
	GetIssue() Issue
	VerifyIssue(ctx context.Context, req VerifyIssueReq) (bool, error)
}

type powSvc struct {
}

func NewPowSvc() PowService {
	return &powSvc{}
}

func (p *powSvc) GetIssue() Issue {
	return Issue{
		Nonce:      string(p.getNonce()),
		Difficulty: difficulty,
	}
}

func (p *powSvc) VerifyIssue(ctx context.Context, req VerifyIssueReq) (bool, error) {
	pass := false
	hashBts, err := hex.DecodeString(req.Hash)
	if err != nil {
		return pass, errors.New("invalid hash")
	}

	pass = p.verifyHash([]byte(req.Nonce), []byte(strconv.Itoa(req.Counter)), hashBts)
	if !pass {
		return pass, errors.New("failed to verify")
	}
	return pass, nil
}

func (p *powSvc) getNonce() []byte {
	return []byte(uuid.NewString())
}

func (p *powSvc) verifyHash(nonce, counter, hash []byte) bool {
	h := p.getChecksum(counter, nonce)
	return bytes.Equal(h, hash)
}

func (p *powSvc) getChecksum(counter, nonce []byte) []byte {
	merge := append(counter, nonce...)
	h := sha256.Sum256(merge)
	return h[:]
}

type Issue struct {
	Nonce      string `json:"nonce"`
	Difficulty int    `json:"difficulty"`
}

type VerifyIssueReq struct {
	Nonce   string `json:"nonce"`
	Counter int    `json:"counter"`
	Hash    string `json:"hash"`
}
