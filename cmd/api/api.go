package main

import (
	"net/http"
	"time"

	"github.com/AmiyoKm/book_store/internal/store"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
)

type Application struct {
	cfg    Config
	logger *zap.SugaredLogger
	store  store.Storage
}

type Config struct {
	addr   string
	apiUrl string
	env    string
	db     DbConfig
}
type DbConfig struct {
	addr        string
	maxConnOpen int
	maxIdleConn int
	maxIdleTime string
}

func (app *Application) mount() http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Use(middleware.Timeout(time.Second * 60))

	r.Route("/api/v1", func(r chi.Router) {
		r.Get("/health", app.healthCheckHandler)

		r.Route("/books", func(r chi.Router) {
			r.Post("/", app.createBookHandler)

			r.Route("/{bookID}", func(r chi.Router) {
				r.Use(app.bookContextMiddleware)

				r.Get("/", app.getBookHandler)
				r.Patch("/" ,app.updateBookHandler)
				r.Delete("/" , app.deleteBookHandler)
			})
		})
	})

	return r
}

func (app *Application) run(mux http.Handler) error {
	server := http.Server{
		Addr:         app.cfg.addr,
		Handler:      mux,
		ReadTimeout:  time.Second * 30,
		WriteTimeout: time.Second * 30,
		IdleTimeout:  time.Second * 30,
	}
	app.logger.Infow("server has started", "addr", app.cfg.addr, "env", app.cfg.env)
	if err := server.ListenAndServe(); err != nil {
		switch err {
		case http.ErrServerClosed:
			return err
		}
	}

	app.logger.Infow("server has stopped", "addr", app.cfg.addr)
	return nil
}

