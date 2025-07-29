package fetch

import (
	"context"
	"crypto/sha1" //nolint:gosec
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"path"
	"strconv"
	"time"

	"gabe565.com/trmnl-nightscout/internal/config"
	"gabe565.com/trmnl-nightscout/internal/nightscout"
	"gabe565.com/trmnl-nightscout/internal/util"
	"golang.org/x/sync/errgroup"
)

var (
	ErrHTTP        = errors.New("unexpected HTTP error")
	ErrNotModified = errors.New("not modified")
)

func NewFetch(conf *config.Config) *Fetch {
	f := &Fetch{
		config: conf,
		client: &http.Client{
			Transport: util.NewUserAgentTransport("trmnl-nightscout", conf.Version),
			Timeout:   time.Minute,
		},
	}

	if token := conf.NightscoutToken; token != "" {
		rawChecksum := sha1.Sum([]byte(token)) //nolint:gosec
		f.tokenChecksum = hex.EncodeToString(rawChecksum[:])
		slog.Debug("Generated token checksum", "value", f.tokenChecksum)
	}

	return f
}

type Fetch struct {
	config         *config.Config
	client         *http.Client
	tokenChecksum  string
	propertiesEtag string
	entriesEtag    string
}

type Response struct {
	Properties nightscout.Properties
	Entries    []nightscout.SGVv1
}

func (f *Fetch) Do(ctx context.Context) (*Response, error) {
	response := &Response{}

	group, ctx := errgroup.WithContext(ctx)

	group.Go(func() error {
		properties, err := f.fetchProperties(ctx)
		if err != nil {
			return err
		}
		response.Properties = *properties
		return nil
	})

	group.Go(func() error {
		entries, err := f.fetchEntries(ctx)
		if err != nil {
			return err
		}
		response.Entries = entries
		return err
	})

	err := group.Wait()
	if err != nil {
		return nil, err
	}
	return response, err
}

func (f *Fetch) request(ctx context.Context, url, etag string, target any) (*http.Response, error) {
	start := time.Now()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	if etag != "" {
		req.Header.Set("If-None-Match", etag)
	}

	if f.tokenChecksum != "" {
		req.Header.Set("Api-Secret", f.tokenChecksum)
	}

	slog.Debug("Fetching data",
		"etag", etag != "",
		"secret", f.tokenChecksum != "",
	)

	res, err := f.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() {
		_, _ = io.Copy(io.Discard, res.Body)
		_ = res.Body.Close()
	}()

	switch res.StatusCode {
	case http.StatusNotModified:
		slog.Debug("Data was not modified", "took", time.Since(start))
		return res, ErrNotModified
	case http.StatusOK:
		// Decode JSON
		if err := json.NewDecoder(res.Body).Decode(target); err != nil {
			return nil, err
		}

		slog.Debug("Parsed response", "took", time.Since(start), "data", target)
		return res, nil
	default:
		return res, fmt.Errorf("%w: %d", ErrHTTP, res.StatusCode)
	}
}

func (f *Fetch) fetchProperties(ctx context.Context) (*nightscout.Properties, error) {
	u, err := url.Parse(f.config.NightscoutURL)
	if err != nil {
		return nil, err
	}

	u.Path = path.Join(u.Path, "/api/v2/properties/bgnow,delta,direction")

	var properties *nightscout.Properties
	res, err := f.request(ctx, u.String(), f.propertiesEtag, &properties)
	if err != nil {
		if !errors.Is(err, ErrNotModified) {
			f.propertiesEtag = ""
		}
		return nil, err
	}
	_ = res.Body.Close()

	f.propertiesEtag = res.Header.Get("etag")
	return properties, nil
}

func (f *Fetch) fetchEntries(ctx context.Context) ([]nightscout.SGVv1, error) {
	u, err := url.Parse(f.config.NightscoutURL)
	if err != nil {
		return nil, err
	}

	u.Path = path.Join(u.Path, "/api/v1/entries.json")

	q := u.Query()
	count := strconv.FormatInt(int64(f.config.Render.GraphDuration/f.config.UpdateInterval), 10)
	q.Set("count", count)
	u.RawQuery = q.Encode()
	find := strconv.FormatInt(time.Now().Add(-f.config.Render.GraphDuration).Unix(), 10)
	u.RawQuery += "&find[date][$gte]=" + find

	var entries []nightscout.SGVv1
	res, err := f.request(ctx, u.String(), f.entriesEtag, &entries)
	if err != nil {
		if !errors.Is(err, ErrNotModified) {
			f.entriesEtag = ""
		}
		return nil, err
	}
	_ = res.Body.Close()

	f.entriesEtag = res.Header.Get("etag")
	return entries, nil
}
