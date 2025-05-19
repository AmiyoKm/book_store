package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/AmiyoKm/book_store/docs"
	"github.com/AmiyoKm/book_store/internal/auth"
	mailer "github.com/AmiyoKm/book_store/internal/mail"
	"github.com/AmiyoKm/book_store/internal/store"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	httpSwagger "github.com/swaggo/http-swagger/v2"
	"go.uber.org/zap"
)

type Application struct {
	cfg    Config
	logger *zap.SugaredLogger
	store  store.Storage
	mail   mailer.Client
	auth   auth.Authenticator
}

type Config struct {
	addr        string
	auth        authConfig
	apiUrl      string
	env         string
	db          DbConfig
	mail        MailConfig
	frontendURL string
}
type authConfig struct {
	basic basicConfig
	token tokenConfig
}

type tokenConfig struct {
	secret string
	exp    time.Duration
	iss    string
}
type basicConfig struct {
	user string
	pass string
}
type DbConfig struct {
	addr        string
	maxConnOpen int
	maxIdleConn int
	maxIdleTime string
}

type MailConfig struct {
	apiKey    string
	fromEmail string
	exp       time.Duration
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

		docsURL := fmt.Sprintf("%s/swagger/doc.json", app.cfg.addr)
		r.Get("/swagger/*", httpSwagger.Handler(httpSwagger.URL(docsURL)))

		r.Route("/authentication", func(r chi.Router) {
			r.Post("/user", app.createUserHandler)
			r.Post("/token", app.createTokenHandler)
		})

		r.Route("/users", func(r chi.Router) {
			r.Use(app.AuthTokenMiddleware)

			r.Route("/me", func(r chi.Router) {
				r.Get("/", app.getUserHandler)
				r.Patch("/", app.updateUserHandler)
			})

			r.Route("/password", func(r chi.Router) {
				r.Post("/reset-request", app.passwordResetRequestHandler)
				r.Get("/request/verify" ,app.passwordRequestVerifyHandler)
				r.Post("/reset", app.passwordResetHandler)
			})

		})

		r.Route("/books", func(r chi.Router) {
			r.Use(app.AuthTokenMiddleware)

			r.Post("/", app.checkBookManipulationAuthority("moderator", app.createBookHandler))

			r.Route("/{bookID}", func(r chi.Router) {
				r.Use(app.bookContextMiddleware)

				r.Get("/", app.getBookHandler)

				r.Patch("/", app.checkBookManipulationAuthority("moderator", app.updateBookHandler))
				r.Delete("/", app.checkBookManipulationAuthority("moderator", app.deleteBookHandler))
			})
		})

		r.Route("/orders", func(r chi.Router) {
			r.Use(app.AuthTokenMiddleware)

			r.Post("/", app.createOrderHandler)
			r.Get("/", app.getAllOrdersHandler)

			r.Route("/{orderID}", func(r chi.Router) {
				r.Use(app.orderContextMiddleware)
				r.Get("/", app.getOrderHandler)

				r.Patch("/", app.updateOderHandler)

			})
		})
		r.Route("/admin", func(r chi.Router) {
			r.Use(app.AuthTokenMiddleware)
			r.Use(app.adminCheck)

			r.Route("/orders", func(r chi.Router) {

				r.Route("/{orderID}", func(r chi.Router) {
					r.Use(app.orderContextMiddleware)

					r.Patch("/", app.updateAdminOrderHandler)
				})
			})
		})
	})

	return r
}

func (app *Application) run(mux http.Handler) error {
	docs.SwaggerInfo.Version = version
	docs.SwaggerInfo.Host = app.cfg.apiUrl
	docs.SwaggerInfo.BasePath = "/api/v1"
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
