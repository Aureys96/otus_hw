package internalhttp

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/Aureys96/hw12_13_14_15_calendar/internal/server/config"
	"github.com/Aureys96/hw12_13_14_15_calendar/internal/storage"
	"github.com/justinas/alice"
	"go.uber.org/zap"
)

type Server struct {
	logger *zap.Logger
	config config.ServerConfig
	srv    *http.Server
	app    Application
}

type Application interface {
	CreateEvent(ctx context.Context, event storage.Event) error
	Get(ctx context.Context, id int) (storage.Event, error)
	Update(ctx context.Context, id int, event storage.Event) error
	Delete(ctx context.Context, id int)
	ListEvents(ctx context.Context, start, end time.Time) ([]storage.Event, error)
}

func NewServer(logger *zap.Logger, app Application, config config.ServerConfig) *Server {
	srv := &http.Server{
		Addr: fmt.Sprintf("%s:%d", config.Host, config.Port),
	}
	return &Server{logger, config, srv, app}
}

func (s *Server) Start() error {
	s.srv.Handler = s.routes()
	s.logger.Info("start server",
		zap.String("host", s.config.Host),
		zap.Int("port", s.config.Port),
	)

	return s.srv.ListenAndServe()
}

func (s *Server) Stop(ctx context.Context) error {
	return s.srv.Shutdown(ctx)
}

func (s *Server) routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/", s.home())
	return alice.New(handlePanic(s.logger), setContent, logRequest(s.logger)).Then(mux)
}

func (s *Server) home() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "hello world")
		w.WriteHeader(200)
	}
}
