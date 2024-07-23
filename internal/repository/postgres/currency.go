package postgres

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kartochnik010/test-kmf/internal/domain"
	"github.com/kartochnik010/test-kmf/internal/domain/models"
	"github.com/kartochnik010/test-kmf/internal/pkg/logger"
)

type CurrencyRepo struct {
	db *pgxpool.Pool
}

func NewCurrencyRepo(db *pgxpool.Pool) *CurrencyRepo {
	return &CurrencyRepo{
		db: db,
	}
}

func (c *CurrencyRepo) SaveRates(ctx context.Context, date string, rates []models.RateItem) error {
	log := logger.GetLoggerFromCtx(ctx).WithField("op", "CurrencyRepo.SaveItem")
	log.Debugf("saving %v items", len(rates))

	query := `
		INSERT INTO r_currency (title, code, value, a_date)
		VALUES ($1, $2, $3, $4)
	`
	tx, err := c.db.Begin(ctx)
	if err != nil {
		log.WithError(err).Error("failed to begin transaction")
		return domain.ErrInternal
	}

	for _, r := range rates {
		_, err := tx.Exec(ctx, query, r.FullName, r.Title, r.Change, date)
		if err != nil {
			log.WithError(err).Error("failed to save item")
			return domain.ErrInternal
		}
	}

	return tx.Commit(ctx)
}

func (c *CurrencyRepo) GetRatesByDate(ctx context.Context, date string) ([]models.Rate, error) {
	log := logger.GetLoggerFromCtx(ctx).WithField("op", "CurrencyRepo.GetRates")
	log.Debugf("getting rates from %v", date)

	query := `
	SELECT title, code, value, a_date
	FROM r_currency
	WHERE a_date = $1`

	rows, err := c.db.Query(ctx, query, date)
	if err != nil {
		log.WithError(err).Error("failed to get rates")
		return nil, err
	}
	defer rows.Close()

	rates := []models.Rate{}
	for rows.Next() {
		r := models.Rate{}
		err := rows.Scan(&r.Title, &r.Code, &r.Value, &r.Date)
		if err != nil {
			log.WithError(err).Error("failed to scan row")
			return nil, err
		}
		rates = append(rates, r)
	}
	return rates, nil
}

func (c *CurrencyRepo) GetRatesByDateAndCode(ctx context.Context, date string, code string) ([]models.Rate, error) {
	log := logger.GetLoggerFromCtx(ctx).WithField("op", "CurrencyRepo.GetRates")
	log.Debugf("getting rates from %v with code %v", date, code)

	query := `
	SELECT title, code, value, a_date
	FROM r_currency
	WHERE a_date = $1 AND code = $2`

	rows, err := c.db.Query(ctx, query, date, code)
	if err != nil {
		log.WithError(err).Error("failed to get rates")
		return nil, err
	}
	defer rows.Close()

	rates := []models.Rate{}
	for rows.Next() {
		r := models.Rate{}
		err := rows.Scan(&r.Title, &r.Code, &r.Value, &r.Date)
		if err != nil {
			log.WithError(err).Error("failed to scan row")
			return nil, err
		}
		rates = append(rates, r)
	}

	return rates, nil
}
