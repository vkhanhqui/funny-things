package app

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type svcWithCache struct {
	svc   PowService
	cache *redis.Client
}

func NewSvcWithCache(svc PowService, cache *redis.Client) PowService {
	return &svcWithCache{svc, cache}
}

func (s *svcWithCache) GetIssue() Issue {
	return s.svc.GetIssue()
}

func (s *svcWithCache) VerifyIssue(ctx context.Context, req VerifyIssueReq) (bool, error) {
	pass := false
	exist := s.cache.Exists(ctx, req.Nonce).Val()
	if exist != 0 {
		return pass, errors.New("already verified")
	}

	pass, err := s.svc.VerifyIssue(ctx, req)
	if err != nil || !pass {
		return pass, err
	}

	s.cache.Set(ctx, req.Nonce, 1, time.Minute)
	return pass, nil
}

func ConnectRedis(endpoint string, port int, password string) (client *redis.Client, err error) {
	client = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", endpoint, port),
		Password: password,
		DB:       0,
	})

	_, err = client.Ping(context.Background()).Result()
	return
}
