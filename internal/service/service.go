package service

import (
	"net/http"

	"github.com/kartochnik010/test-kmf/internal/repository"
)

type Service struct {
	Currency *CurrencyService
}

func NewService(repo repository.Repository, c *http.Client) *Service {
	return &Service{
		Currency: NewCurrencyService(c, repo),
	}
}
