package main

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type Limiter struct {
	mu       sync.Mutex
	duration time.Duration
	limit    int
	slots    int
}

func NewLimiter(req int, t time.Duration) Limiter {
	return Limiter{
		duration: t,
		limit:    req,
	}
}

func (l *Limiter) Start() chan<- bool {
	stop := make(chan bool)
	ticker := time.NewTicker(l.duration)
	l.mu.Lock()
	l.slots = l.limit
	defer l.mu.Unlock()
	go func() {
		for {
			select {
			case <-stop:
				break
			case <-ticker.C:
				l.mu.Lock()
				l.slots = l.limit
				l.mu.Unlock()
			}
		}
	}()
	return stop
}

func (l *Limiter) AllowAndTake() bool {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.slots > 0 {
		l.slots -= 1
		return true
	}
	return false
}

func main() {
	r := gin.Default()
	limiter := NewLimiter(5, time.Minute)
	stop := limiter.Start()
	r.GET("/", func(ctx *gin.Context) {
		if limiter.AllowAndTake() {
			ctx.String(http.StatusOK, "The resource is available")
			return
		}
		ctx.String(http.StatusServiceUnavailable, "The resource is closed")
	})
	r.Run()
	stop <- true
}
