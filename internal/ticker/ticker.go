package ticker

import (
	"context"
	"sync"
	"time"

	"gabe565.com/trmnl-nightscout/internal/config"
	"gabe565.com/trmnl-nightscout/internal/fetch"
)

func New(conf *config.Config) *Ticker {
	return &Ticker{
		config: conf,
		fetch:  fetch.NewFetch(conf),
	}
}

type Ticker struct {
	cancel context.CancelFunc

	config *config.Config
	fetch  *fetch.Fetch

	fetchTicker *time.Ticker

	mu    sync.RWMutex
	last  *fetch.Response
	error error
}

func (t *Ticker) Start(ctx context.Context) *Ticker {
	ctx, t.cancel = context.WithCancel(ctx)
	t.beginFetch(ctx)
	return t
}

func (t *Ticker) Close() {
	if t.cancel != nil {
		t.cancel()
	}
}

func (t *Ticker) Last() (*fetch.Response, error) {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return t.last, t.error
}
