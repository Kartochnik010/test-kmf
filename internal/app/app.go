package app

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kartochnik010/test-kmf/internal/config"
	"github.com/kartochnik010/test-kmf/internal/handler"
	"github.com/kartochnik010/test-kmf/internal/repository"
	"github.com/kartochnik010/test-kmf/internal/service"
	"github.com/sirupsen/logrus"
)

type App struct {
	Server *http.Server
	Logger *logrus.Logger
	Repo   repository.Repository
}

func NewApp(cfg *config.Config, db *pgxpool.Pool, l *logrus.Logger) *App {
	const op = "app.NewApp"

	repo := repository.NewRepository(db, l)

	service := service.NewService(repo, &http.Client{})

	h := handler.NewHandler(&http.Client{}, service, l, cfg)

	s := &http.Server{
		Addr:    fmt.Sprintf(":%v", cfg.Port),
		Handler: handler.Routes(h),
	}
	return &App{
		Server: s,
		Logger: l,
		Repo:   repo,
	}
}

func (a *App) Run() error {
	return a.Server.ListenAndServe()
}

func OpenDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("sqlserver", dsn)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
