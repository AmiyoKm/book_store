package main

import (
	"github.com/AmiyoKm/book_store/internal/db"
	"github.com/AmiyoKm/book_store/internal/env"
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

	config := Config{
		db:     dbConfig,
		env:    env.GetString("ENVIRONMENT", "DEVELOPMENT"),
		addr:   env.GetString("ADDR", ":8080"),
		apiUrl: env.GetString("API_URL", "localhost:8080"),
	}

	db, err := db.New(config.db.addr, config.db.maxConnOpen, config.db.maxIdleConn, config.db.maxIdleTime)
	if err != nil {
		logger.Fatal(err)
	}
	defer db.Close()
	logger.Info("DB connection pool established")

	store := store.NewStorage(db)

	app := &Application{
		cfg:    config,
		logger: logger,
		store:  store,
	}
	mux := app.mount()
	logger.Fatal(app.run(mux))

}
