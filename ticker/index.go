package ticker_test

import (
	"sync"
	"time"
)

type Worker struct {
	ticker *time.Ticker
	GoFunc
}

type GoFunc struct {
	swg sync.WaitGroup
}

func (w *Worker) New(duration time.Duration, f func()) {
	w.ticker = time.NewTicker(duration)
	select {
	case v := <-w.ticker.C:
		go func(params time.Time) {

		}(v)
	}
}
