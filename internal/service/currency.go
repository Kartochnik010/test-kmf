package service

import (
	"context"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"

	"github.com/kartochnik010/test-kmf/internal/domain"
	"github.com/kartochnik010/test-kmf/internal/domain/models"
	"github.com/kartochnik010/test-kmf/internal/pkg/logger"
	"github.com/kartochnik010/test-kmf/internal/repository"
	"github.com/sirupsen/logrus"
)

type CurrencyService struct {
	c    *http.Client
	repo repository.Currency
}

func NewCurrencyService(client *http.Client, repo repository.Currency) *CurrencyService {
	return &CurrencyService{
		c:    client,
		repo: repo,
	}
}

func (s *CurrencyService) FetchAndSaveRates(ctx context.Context, date string) error {
	log := logger.GetLoggerFromCtx(ctx).WithField("op", "CurrencyService.SaveRates")

	var rates models.RatesXML

	// fetch rates
	resp, err := s.c.Get(fmt.Sprintf("https://nationalbank.kz/rss/get_rates.cfm?fdate=%s", date))
	if err != nil {
		log.WithError(err).Error("failed to fetch rates")
		return domain.ErrInternal
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.WithError(err).Error("failed to read body")
		return domain.ErrInternal
	}

	err = xml.Unmarshal(body, &rates)
	if err != nil {
		log.WithError(err).Error("failed to unmarshal xml")
		return domain.ErrInternal
	}

	// save rates
	if err := s.repo.SaveRates(ctx, date, rates.Items); err != nil {
		log.WithError(err).Error("failed to save rates")
		return domain.ErrInternal
	}
	return nil
}

func (s *CurrencyService) GetRatesByDate(ctx context.Context, date string) ([]models.Rate, error) {
	log := logger.GetLoggerFromCtx(ctx).WithField("op", "CurrencyService.GetRatesByDate")

	rates, err := s.repo.GetRatesByDate(ctx, date)
	if err != nil {
		log.WithError(err).WithField("date", date).Error("failed to get rates by date")
		return nil, domain.ErrInternal
	}

	return rates, nil
}

func (s *CurrencyService) GetRatesByDateAndCode(ctx context.Context, date string, code string) ([]models.Rate, error) {
	log := logger.GetLoggerFromCtx(ctx).WithField("op", "CurrencyService.GetRatesByDateAndCode")

	rates, err := s.repo.GetRatesByDateAndCode(ctx, date, code)
	if err != nil {
		log.WithError(err).WithFields(logrus.Fields{
			"date": date,
			"code": code,
		}).Error("failed to get rates by date and code")
		return nil, domain.ErrInternal
	}

	return rates, nil
}
