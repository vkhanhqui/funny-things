package app

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/base64"
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
		Nonce:      p.getNonce(),
		Difficulty: difficulty,
	}
}

func (p *powSvc) VerifyIssue(ctx context.Context, req VerifyIssueReq) (bool, error) {
	pass := false
	parser, err := req.toParser()
	if err != nil {
		return false, err
	}

	pass = p.verifyHash(parser.Nonce, parser.Hash, parser.Counter)
	if !pass {
		return pass, errors.New("failed to verify hash")
	}

	pass = p.verifyDifficulty(parser.Bin, req.Difficulty)
	if !pass {
		return pass, errors.New("failed to verify difficulty")
	}
	return pass, nil
}

func (p *powSvc) verifyDifficulty(binBts []byte, difficulty int) bool {
	bin := string(binBts)
	for _, c := range bin[:difficulty] {
		if c != '0' {
			return false
		}
	}
	return true
}

func (p *powSvc) getNonce() string {
	return base64.StdEncoding.EncodeToString([]byte(uuid.NewString()))
}

func (p *powSvc) verifyHash(nonce, hash []byte, counter int) bool {
	counterBts := []byte(strconv.Itoa(counter))
	h := p.getChecksum(counterBts, nonce)
	return bytes.Equal(h, hash)
}

func (p *powSvc) getChecksum(counter, nonce []byte) []byte {
	merge := append(counter, nonce...)
	h := sha256.Sum256(merge)
	return h[:]
}
