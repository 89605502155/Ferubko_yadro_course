package personal_limiter

import (
	"context"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

type PersonalLimiter struct {
	limit   int
	inteval time.Duration
	list    map[string][]int
	mutex   sync.Mutex
}

func NewPersonalLimiter(ctx context.Context, limit int, interval time.Duration) *PersonalLimiter {
	struct_ := &PersonalLimiter{
		limit:   limit,
		inteval: interval,
		list:    make(map[string][]int, 0),
	}
	go struct_.statrPerionForRefresh(ctx, interval)
	return struct_
}

func (p *PersonalLimiter) statrPerionForRefresh(ctx context.Context, duration time.Duration) {
	timer := time.NewTimer(duration)
	defer timer.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-timer.C:
			p.mutex.Lock()
			p.list = make(map[string][]int, 0)
			p.mutex.Unlock()
		}
	}

}

func (p *PersonalLimiter) Allow(userName string, hard int) bool {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	logrus.Println(p.limit, p.list[userName], hard)
	sl := p.list[userName]
	s := 0
	for j := 0; j < len(sl); j++ {
		s += (sl[j] + 1)
	}
	if s <= p.limit-hard-1 {
		s += hard
		s += 1
		sl = append(sl, hard)
		p.list[userName] = sl
		return true
	} else {
		p.list[userName] = sl
		return false
	}
}
