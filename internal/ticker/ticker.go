package ticker

import (
	"context"
	"errors"
	"sync"
	"time"

	"gabe565.com/trmnl-nightscout/internal/config"
	"gabe565.com/trmnl-nightscout/internal/fetch"
)

func New(conf *config.Config) (*Ticker, error) {
	f, err := fetch.New(conf)
	if err != nil {
		return nil, err
	}

	return &Ticker{
		config: conf,
		fetch:  f,
	}, nil
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

var ErrNotLoaded = errors.New("reading not loaded")

func (t *Ticker) Last() (*fetch.Response, error) {
	t.mu.RLock()
	last, err := t.last, t.error
	defer t.mu.RUnlock()
	if last == nil && err == nil {
		err = ErrNotLoaded
	}
	return last, err
}
