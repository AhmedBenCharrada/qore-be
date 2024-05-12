package server

import (
	"context"
	"qore-be/internal/config"
	"qore-be/internal/person"

	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// Server the http server.
type Server struct {
	cfg    *config.Config
	person *person.Controller

	router *gin.Engine
}

// Option http server option factory type.
type Option = func(*Server) error

// New creates a new http server.
func New(opts ...Option) *Server {
	srv := &Server{}

	for _, opt := range opts {
		if err := opt(srv); err != nil {
			panic(err)
		}
	}

	srv.setupRouter()

	return srv
}

// WithConfig initialize the http server with the project configuration.
func WithConfig(cfg *config.Config) Option {
	return func(svc *Server) error {
		svc.cfg = cfg
		return nil
	}
}

// WithPersonController initialize the server with a person controller.
func WithPersonController(c *person.Controller) Option {
	return func(svc *Server) error {
		if c == nil {
			return fmt.Errorf("nil store controller")
		}
		svc.person = c
		return nil
	}
}

// Start starts the http server.
func (s *Server) Start(ctx context.Context) error {
	server := &http.Server{
		Addr:    s.cfg.Host,
		Handler: s.router,
	}

	ch := make(chan error)
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			ch <- err
		}
	}()

	for {
		select {
		case <-ctx.Done():
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			err := server.Shutdown(ctx)
			cancel()
			if err != nil {
				slog.Error("failed to shutdown server", "error", err.Error())
				os.Exit(-1)
			}
		case err := <-ch:
			if err != http.ErrServerClosed {
				slog.Error("failed to start the server", "error", err.Error())
				os.Exit(-1)
			}
		}
	}
}

func (s *Server) setupRouter() {
	router := gin.Default()

	// CORS middleware
	router.Use(cors.Default())

	personCtrl := router.Group("/person")
	personCtrl.GET("", s.person.GetAll)
	personCtrl.POST("/create", s.person.Create)
	personCtrl.GET("/:id/info", s.person.GetByID)

	s.router = router
}
