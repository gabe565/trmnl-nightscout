package ticker

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"gabe565.com/trmnl-nightscout/internal/fetch"
)

func (t *Ticker) beginFetch(ctx context.Context) {
	go func() {
		t.fetchTicker = time.NewTicker(t.config.FallbackInterval)
		defer t.fetchTicker.Stop()

		for {
			next := t.Fetch(ctx)
			t.fetchTicker.Reset(next)
			slog.Debug("Scheduled next fetch", "in", next)

			select {
			case <-ctx.Done():
				return
			case <-t.fetchTicker.C:
			}
		}
	}()
}

func (t *Ticker) Fetch(ctx context.Context) time.Duration {
	res, err := t.fetch.Do(ctx)
	if err != nil {
		if !errors.Is(err, fetch.ErrNotModified) {
			t.mu.Lock()
			slog.Error("Failed to fetch", "err", err)
			t.error = err
			t.mu.Unlock()
		}
		return t.config.FallbackInterval
	}

	t.mu.Lock()
	t.last = res
	t.error = nil
	t.mu.Unlock()

	next := time.Until(res.Properties.Bgnow.Mills.Time) + t.config.UpdateInterval + t.config.FetchDelay
	return max(next, t.config.FallbackInterval)
}
