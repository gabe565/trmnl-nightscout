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
	"strings"
	"time"

	"gabe565.com/trmnl-nightscout/internal/config"
	"gabe565.com/trmnl-nightscout/internal/ticker"
	"gabe565.com/trmnl-nightscout/internal/trmnl"
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
	s.ticker = ticker.New(s.conf).Start(ctx)

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
		Addr:        s.conf.ListenAddress,
		Handler:     r,
		ReadTimeout: 5 * time.Second,
	}

	group, ctx := errgroup.WithContext(ctx)

	group.Go(func() error {
		slog.Info("Listening for connections", "address", s.conf.ListenAddress)
		return server.ListenAndServe()
	})

	group.Go(func() error {
		<-ctx.Done()

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		slog.Info("Gracefully shutting down server")
		return server.Shutdown(ctx)
	})

	err := group.Wait()
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
	}

	u, err := url.Parse(s.conf.PublicURL)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	u.Path = path.Join(u.Path, "image.png")

	if token := r.URL.Query().Get("token"); token != "" {
		q := u.Query()
		q.Set("token", token)
		u.RawQuery = q.Encode()
	}

	buf := bytes.NewBuffer(make([]byte, 0, 256))
	if err := json.NewEncoder(buf).Encode(trmnl.Redirect{
		Filename:    "nightscout-" + stamp.Format(time.RFC3339),
		URL:         u.String(),
		RefreshRate: time.Minute,
	}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.ServeContent(w, r, "image.json", stamp, strings.NewReader(buf.String()))
}

func (s *Server) image(w http.ResponseWriter, r *http.Request) {
	last, err := s.ticker.Last()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	img, err := trmnl.Render(s.conf, last)
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
