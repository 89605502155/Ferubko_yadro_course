package rate_limiter

import (
	"sync"
	"time"
)

type SlidingLog struct {
	logs time.Time
	hard int
}

type SlidindLogLimiter struct {
	limit   int
	inteval time.Duration
	logs    []SlidingLog
	mutex   sync.Mutex
}

func NewSlidingLogLimiter(limit int, inteval time.Duration) *SlidindLogLimiter {
	return &SlidindLogLimiter{
		limit:   limit,
		inteval: inteval,
		logs:    make([]SlidingLog, 0),
	}
}

func (l *SlidindLogLimiter) Allow(hard int, dominantus bool) bool {
	// у запросов на поиск будет hard=0, а у запросов на обновление пусть будет 250, чтобы нельзя было запустить много обновлений.
	l.mutex.Lock()
	defer l.mutex.Unlock()

	lastPeriod := time.Now().Add(-l.inteval)
	for len(l.logs) != 0 && l.logs[0].logs.Add(time.Duration(l.logs[0].hard)*time.Millisecond).Before(lastPeriod) {
		l.logs = l.logs[1:]
	}

	newRequest := SlidingLog{
		logs: time.Now(),
		hard: hard,
	}
	l.logs = append(l.logs, newRequest)
	s := 0
	for i := range l.logs {
		s += (l.logs[i].hard + 1)
	}
	if dominantus {
		if s <= l.limit*2 {
			return true
		} else {
			return false
		}
	} else {
		return s <= l.limit
	}

}
