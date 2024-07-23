package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kartochnik010/test-kmf/internal/domain/models"
	"github.com/kartochnik010/test-kmf/internal/repository/postgres"
	"github.com/sirupsen/logrus"
)

type Currency interface {
	SaveRates(ctx context.Context, date string, rates []models.RateItem) error
	GetRatesByDate(ctx context.Context, date string) ([]models.Rate, error)
	GetRatesByDateAndCode(ctx context.Context, date string, code string) ([]models.Rate, error)
}

type Repository struct {
	Currency
}

func NewRepository(db *pgxpool.Pool, l *logrus.Logger) Repository {
	return Repository{
		Currency: postgres.NewCurrencyRepo(db),
	}
}
