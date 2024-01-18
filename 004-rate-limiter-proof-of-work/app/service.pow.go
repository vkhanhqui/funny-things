package app

import (
	"context"
	"encoding/hex"
	"errors"
	"proof-of-work/types"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	secret     = "secret"
	difficulty = 10
)

type PowService interface {
	GetIssue() types.Issue
	VerifyIssue(ctx context.Context, req types.VerifyIssueReq) (bool, error)
}

type powService struct {
	cache *redis.Client
}

func NewPowSvc(cache *redis.Client) PowService {
	return &powService{cache}
}

func (p *powService) GetIssue() types.Issue {
	n := getNonce()
	checksum := hex.EncodeToString(getChecksum(n, secret))

	return types.Issue{
		Nonce:      string(n),
		Checksum:   checksum,
		Difficulty: difficulty,
	}
}

func (p *powService) VerifyIssue(ctx context.Context, req types.VerifyIssueReq) (bool, error) {
	pass := false

	exist, err := p.cache.Get(ctx, req.Nonce).Result()
	if err != nil && err != redis.Nil {
		return pass, err
	}

	if exist != "" {
		return pass, errors.New("already verified")
	}

	hashBts, err := hex.DecodeString(req.Hash)
	if err != nil {
		return pass, err
	}

	checksumBts, err := hex.DecodeString(req.Checksum)
	if err != nil {
		return pass, err
	}

	pass, err = verify(types.ParsedIssue{
		Nonce:      []byte(req.Nonce),
		Checksum:   checksumBts,
		Counter:    []byte(strconv.Itoa(req.Counter)),
		Hash:       hashBts,
		Difficulty: req.Difficulty,
	})

	if err != nil {
		return pass, err
	}

	if pass {
		p.cache.Set(ctx, req.Nonce, 1, time.Minute)
	}
	return pass, err
}
