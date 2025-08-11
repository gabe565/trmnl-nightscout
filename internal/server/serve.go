package server

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"image/png"
	"log/slog"
	"net/http"
	"net/url"
	"path"
	"time"

	"gabe565.com/trmnl-nightscout/internal/config"
	"gabe565.com/trmnl-nightscout/internal/ticker"
	"gabe565.com/trmnl-nightscout/internal/trmnl"
	"gabe565.com/utils/bytefmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/httprate"
	"golang.org/x/sync/errgroup"
)

func New(conf *config.Config) *Server {
	return &Server{
		conf:    conf,
		started: time.Now(),
	}
}

type Server struct {
	conf    *config.Config
	ticker  *ticker.Ticker
	started time.Time
}

func (s *Server) ListenAndServe(ctx context.Context) error {
	var err error
	if s.ticker, err = ticker.New(s.conf); err != nil {
		return err
	}

	s.ticker.Start(ctx)

	r := chi.NewRouter()

	r.Use(middleware.Heartbeat("/ping"))
	if s.conf.RealIPHeader {
		r.Use(middleware.RealIP)
	}
	r.Use(middleware.Logger)
	r.Use(middleware.GetHead)
	r.Use(middleware.Recoverer)
	r.Use(httprate.LimitByIP(10, time.Minute))
	r.Use(middleware.Timeout(60 * time.Second))
	if s.conf.AccessToken != "" {
		r.Use(Token(s.conf.AccessToken))
	}

	r.Get("/", s.json)
	r.Get("/image.png", s.image)

	server := &http.Server{
		Addr:           s.conf.ListenAddress,
		Handler:        r,
		ReadTimeout:    5 * time.Second,
		MaxHeaderBytes: 100 * bytefmt.KiB,
	}

	group, ctx := errgroup.WithContext(ctx)

	group.Go(func() error {
		log := slog.With("address", s.conf.ListenAddress)
		if s.conf.TLSCertPath != "" && s.conf.TLSKeyPath != "" {
			log.Info("Listening for https connections")
			return server.ListenAndServeTLS(s.conf.TLSCertPath, s.conf.TLSKeyPath)
		}
		log.Info("Listening for http connections")
		return server.ListenAndServe()
	})

	group.Go(func() error {
		<-ctx.Done()

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		slog.Info("Gracefully shutting down server")
		return server.Shutdown(ctx)
	})

	err = group.Wait()
	if errors.Is(err, context.Canceled) || errors.Is(err, http.ErrServerClosed) {
		err = nil
	}
	return err
}

func (s *Server) json(w http.ResponseWriter, r *http.Request) {
	last, err := s.ticker.Last()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	stamp := last.Properties.Bgnow.Mills.Time
	if stamp.Before(s.started) {
		stamp = s.started
	} else {
		age := time.Since(stamp)
		if age > s.conf.UpdateInterval+s.conf.FetchDelay {
			missed := int(age / s.conf.UpdateInterval)
			stamp = stamp.Add(time.Duration(missed) * s.conf.UpdateInterval)
		}
	}

	var u *url.URL
	switch {
	case s.conf.ImageURL != "":
		if u, err = url.Parse(s.conf.ImageURL); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	case r.Host != "":
		u = &url.URL{
			Host:   r.Host,
			Scheme: "http",
		}
		if r.TLS != nil || r.Header.Get("X-Forwarded-Proto") == "https" {
			u.Scheme = "https"
		}
	default:
		http.Error(w, "Image URL unknown. Either set a request host or the IMAGE_URL env.", http.StatusBadRequest)
		return
	}

	u.Path = path.Join(u.Path, "image.png")
	u.RawQuery = r.URL.RawQuery

	refreshRate := time.Until(last.Properties.Bgnow.Mills.Time) + s.conf.UpdateInterval + s.conf.FetchDelay
	refreshRate = max(refreshRate, 60*time.Second)

	b, err := json.Marshal(trmnl.Redirect{
		Filename:    "nightscout-" + stamp.Format(time.RFC3339),
		URL:         u.String(),
		RefreshRate: refreshRate,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.ServeContent(w, r, "image.json", time.Time{}, bytes.NewReader(b))
}

func (s *Server) image(w http.ResponseWriter, r *http.Request) {
	last, err := s.ticker.Last()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	renderConf := s.conf.Render
	if err := renderConf.UnmarshalQuery(r.URL.Query()); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	img, err := trmnl.NewRenderer(renderConf, last).Render()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	buf := bytes.NewBuffer(make([]byte, 0, 8192))
	if err := png.Encode(buf, img); err != nil {
		slog.Error("Failed to encode image", "error", err)
	}

	http.ServeContent(w, r, "image.png", last.Properties.Bgnow.Mills.Time, bytes.NewReader(buf.Bytes()))
}
