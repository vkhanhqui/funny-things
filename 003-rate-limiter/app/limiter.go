package app

type limiter struct {
	pool map[string]bool
	algo LimiterAlgo
}

func NewLimiter(algo LimiterAlgo) *limiter {
	return &limiter{
		pool: make(map[string]bool),
		algo: algo,
	}
}

func (l *limiter) SendRequest(userID string) bool {
	l.pool[userID] = l.algo.SendRequest()
	return l.pool[userID]
}

type LimiterAlgo interface {
	SendRequest() bool
}
