package personal_limiter

import (
	"sync"
	"time"

	"xkcd/pkg/rate_limiter"
)

type PersonalLimiter struct {
	limit   int
	inteval time.Duration
	list    map[string][]rate_limiter.SlidingLog
	mutex   sync.Mutex
}

func NewPersonalLimiter(limit int, interval time.Duration) *PersonalLimiter {
	return &PersonalLimiter{
		limit:   limit,
		inteval: interval,
		list:    make(map[string][]rate_limiter.SlidingLog, 0),
	}
}

func (p *PersonalLimiter) Allow(userName string, hard int) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	lastPeriod := time.Now().Add(-p.inteval)
	// var logs []rate_limiter.SlidingLog
	// var ok bool
	if _, ok := p.list[userName]; !ok {
		p.list[userName] = make([]rate_limiter.SlidingLog, 0)
	}
	logs := p.list[userName]

	for len(logs) != 0 && p.list[userName][0].logs.Add(time.Duration(p.list[userName][0].hard)*time.Second).Before(lastPeriod) {
		l.logs = l.logs[1:]
	}

	newRequest := SlidingLog{
		logs: time.Now(),
		hard: hard,
	}

}
