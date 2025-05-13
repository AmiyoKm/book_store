package main

import (
	"time"

	"github.com/AmiyoKm/book_store/internal/db"
	"github.com/AmiyoKm/book_store/internal/env"
	"github.com/AmiyoKm/book_store/internal/mail"
	"github.com/AmiyoKm/book_store/internal/store"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

func main() {
	logger := zap.Must(zap.NewProduction()).Sugar()
	defer logger.Sync()

	if err := godotenv.Load(); err != nil {
		logger.Fatal(err)
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
		apiKey: env.GetString("APP_PASSWORD", ""),
	}

	config := Config{
		db:          dbConfig,
		env:         env.GetString("ENVIRONMENT", "DEVELOPMENT"),
		addr:        env.GetString("ADDR", ":8080"),
		apiUrl:      env.GetString("API_URL", "localhost:8080"),
		frontendURL: env.GetString("FRONTEND_URL", "http://localhost:5173"),
		mail: mailCgf,
	}

	db, err := db.New(config.db.addr, config.db.maxConnOpen, config.db.maxIdleConn, config.db.maxIdleTime)
	if err != nil {
		logger.Fatal(err)
	}
	defer db.Close()
	logger.Info("DB connection pool established")

	store := store.NewStorage(db)
	mailClient , err := mailer.NewGoMailClient(config.mail.apiKey , config.mail.fromEmail)
	if err != nil {
		logger.Fatal(err)
	}
	app := &Application{
		cfg:    config,
		logger: logger,
		store:  store,
		mail: mailClient,
	}
	mux := app.mount()
	logger.Fatal(app.run(mux))

}
