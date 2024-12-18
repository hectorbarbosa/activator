package main

import (
	"activator/internal/app/service"
	"activator/internal/config"
	"activator/internal/logging"
	"activator/internal/mailer"
	"activator/internal/rest"
	"activator/internal/storage/postgresql"
	"database/sql"

	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"

	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

func main() {
	// Find path for env file
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		fmt.Fprintln(os.Stderr, "abs path", err)
		os.Exit(1)
	}

	// config file must be in project root dir, compiled bin must be in /bin dir!!!
	configPath, err := filepath.Abs(dir + "/../.env")
	if err != nil {
		fmt.Fprintln(os.Stderr, "composing path", err)
		os.Exit(1)
	}
	// fmt.Println("Config path:", configPath)

	// Create config
	cfg, err := config.NewConfig(configPath)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error loading config", err)
		os.Exit(1)
	}

	// Create db connection
	db, err := newDB(cfg)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error connecting db", err)
		os.Exit(1)
	}

	// Create logger
	logger, err := logging.GetLogger(cfg.LogLevel)
	if err != nil {
		fmt.Printf("getting logger: %s\n", err)
		os.Exit(1)
	}
	// fmt.Println("Log Level:", cfg.LogLevel)

	r := mux.NewRouter()
	repoUser := postgresql.NewUserRepo(db, logger)
	svcUser := service.NewUserService(cfg, logger, repoUser)
	repoToken := postgresql.NewTokenRepo(db, logger)
	svcToken := service.NewTokenService(cfg, logger, repoToken)
	mailer := mailer.New(cfg, logger, cfg.SmtpHost, cfg.SmtpPort, "mk@mk9")

	rest.NewUserHandler(cfg, logger, svcUser, svcToken, mailer).Register(r)

	server := http.Server{
		Handler:           r,
		Addr:              cfg.ServerAddr,
		ReadTimeout:       1 * time.Second,
		ReadHeaderTimeout: 1 * time.Second,
		WriteTimeout:      1 * time.Second,
		IdleTimeout:       1 * time.Second,
	}
	logger.Info("Server start", "listening on address:", cfg.ServerAddr)
	server.ListenAndServe()
}

func newDB(cfg config.Config) (*sql.DB, error) {
	db, err := sql.Open("postgres", cfg.DbUrl)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
