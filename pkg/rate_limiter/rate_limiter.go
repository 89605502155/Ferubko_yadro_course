package rate_limiter

import (
	"sync"
	"time"
)

const (
	timeBetween = 5 * time.Minute
)

type SlidingLog struct {
	logs time.Time
	hard int
}

type SlidindLogLimiter struct {
	limit          int
	inteval        time.Duration
	logs           []SlidingLog
	mutex          sync.Mutex
	timeDominantus time.Time
}

func NewSlidingLogLimiter(limit int, inteval time.Duration) *SlidindLogLimiter {
	// t, _ := strconv.Atoi(os.Getenv("TIME_BETWEEN"))
	// t = 5
	// timeBetwe := time.Duration(t) * time.Minute
	// logrus.Info(t)
	return &SlidindLogLimiter{
		limit:          limit,
		inteval:        inteval,
		logs:           make([]SlidingLog, 0),
		timeDominantus: time.Now().Add(-timeBetween),
	}
}

func (l *SlidindLogLimiter) Allow(hard int, dominantus bool) bool {
	// у запросов на поиск будет hard=0, а у запросов на обновление пусть будет 250 000, чтобы нельзя было запустить много обновлений.
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
		if time.Since(l.timeDominantus) >= timeBetween {
			l.timeDominantus = time.Now()
			return true
		}
		return false
	} else {
		return s <= l.limit
	}

}
