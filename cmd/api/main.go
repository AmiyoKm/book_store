package main

import (
	"time"

	"github.com/AmiyoKm/book_store/internal/auth"
	"github.com/AmiyoKm/book_store/internal/db"
	"github.com/AmiyoKm/book_store/internal/env"
	mailer "github.com/AmiyoKm/book_store/internal/mail"
	"github.com/AmiyoKm/book_store/internal/store"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

const version string = "0.0.1"

//	@title			BookBound API
//	@description	API for BookBound .
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	API Support
//	@contact.url	http://www.swagger.io/support
//	@contact.email	support@swagger.io

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

// @BasePath					/api/v1
//
// @securityDefinitions.apikey	ApiKeyAuth
// @in							header
// @name						Authorization
// @description
func main() {
	logger := zap.Must(zap.NewProduction()).Sugar()
	defer logger.Sync()

	if err := godotenv.Load(); err != nil {
		logger.Info(err)
	}
	dbConfig := DbConfig{
		addr:        env.GetString("DB_ADDR", "postgres://admin:adminpassword@localhost/book_store?sslmode=disable"),
		maxConnOpen: env.GetInt("MAX_CONN_OPEN", 30),
		maxIdleConn: env.GetInt("MAX_IDLE_CONN", 30),
		maxIdleTime: env.GetString("MAX_IDLE_TIME", "15m"),
	}
	mailCgf := MailConfig{
		exp:       time.Hour * 24 * 3,
		fromEmail: env.GetString("FROM_EMAIL", ""),
		apiKey:    env.GetString("APP_PASSWORD", ""),
	}

	authConfig := authConfig{
		basic: basicConfig{
			user: env.GetString("AUTH_BASIC_USER", "admin"),
			pass: env.GetString("AUTH_BASIC_PASS", "admin"),
		},
		token: tokenConfig{
			secret: env.GetString("AUTH_TOKEN_SECRET", "example"),
			exp:    time.Hour * 24 * 3,
			iss:    "BookBound",
		},
	}
	config := Config{
		db:          dbConfig,
		env:         env.GetString("ENVIRONMENT", "DEVELOPMENT"),
		addr:        env.GetString("ADDR", ":8080"),
		apiUrl:      env.GetString("API_URL", "localhost:8080"),
		frontendURL: env.GetString("FRONTEND_URL", "http://localhost:5173"),
		mail:        mailCgf,
		auth:        authConfig,
	}

	db, err := db.New(config.db.addr, config.db.maxConnOpen, config.db.maxIdleConn, config.db.maxIdleTime)
	if err != nil {
		logger.Fatal(err)
	}
	defer db.Close()
	logger.Info("DB connection pool established")

	store := store.NewStorage(db)
	mailClient, err := mailer.NewGoMailClient(config.mail.apiKey, config.mail.fromEmail)
	if err != nil {
		logger.Fatal(err)
	}
	JWTAuthenticator := auth.NewJWTAuthenticator(config.auth.token.secret, config.auth.token.iss, config.auth.token.iss)
	app := &Application{
		cfg:    config,
		logger: logger,
		store:  store,
		mail:   mailClient,
		auth:   JWTAuthenticator,
	}
	mux := app.mount()
	logger.Fatal(app.run(mux))

}
