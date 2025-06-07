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
		t.fetchTicker = time.NewTicker(time.Millisecond)
		defer t.fetchTicker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-t.fetchTicker.C:
				next := t.Fetch()
				t.fetchTicker.Reset(next)
				slog.Debug("Scheduled next fetch", "in", next)
			}
		}
	}()
}

func (t *Ticker) Fetch() time.Duration {
	res, err := t.fetch.Do(context.Background())
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
	t.mu.Unlock()

	if len(res.Properties.Buckets) != 0 {
		bucket := res.Properties.Buckets[0]
		lastDiff := bucket.ToMills.Sub(bucket.FromMills.Time)
		nextRead := res.Properties.Bgnow.Mills.Add(lastDiff + t.config.FetchDelay)
		if until := time.Until(nextRead); until > 0 {
			return until
		}
	}
	return t.config.FallbackInterval
}
