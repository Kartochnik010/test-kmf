// Package classification of Product API
//
// Documentation for Product API
//
//  Schemes: http
//  BasePath: /
//  Version: 1.0.0
//
//  Consumes:
//  - application/json
//
//  Produces:
//  - application/json
// swagger:meta

package handler

import (
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/kartochnik010/test-kmf/docs"
	"github.com/kartochnik010/test-kmf/internal/config"
	_ "github.com/kartochnik010/test-kmf/internal/domain/models"
	"github.com/kartochnik010/test-kmf/internal/pkg/js"
	"github.com/kartochnik010/test-kmf/internal/pkg/logger"
	"github.com/kartochnik010/test-kmf/internal/service"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
)

type Handler struct {
	client *http.Client
	s      *service.Service
	l      *logrus.Logger
	cfg    *config.Config
}

func NewHandler(client *http.Client, s *service.Service, l *logrus.Logger, cfg *config.Config) *Handler {
	return &Handler{
		client: client,
		s:      s,
		l:      l,
		cfg:    cfg,
	}
}

// @Summary Fetch and store rates
// @Description Fetch and store rates by date
// @Tags rates
// @Accept json
// @Produce json
// @Param date path string true "date. Example: '01-01-2022'"
// @Success 200 {object} js.JSON{success=bool}
// @Failure 500 {object} js.JSON{error=string}
// @Router /currency/save/{date} [get]
func (h *Handler) FetchAndStore(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	date := vars["date"]
	if date == "" {
		js.WriteJSON(w, 400, js.JSON{"success": "false", "error": "date is required. Example: '01-01-2022'"}, nil)
		return
	}

	js.WriteJSON(w, 200, js.JSON{"success": true}, nil)

	go func() {
		log := logger.GetLoggerFromCtx(r.Context()).WithField("op", "Handler.FetchAndStore")
		ctx := context.WithValue(context.Background(), logger.ContextKeyLogger, log)
		err := h.s.Currency.FetchAndSaveRates(ctx, date)
		if err != nil {
			log.WithError(err).Error("failed to get rates")
		}
	}()
}

// @Summary Get rates
// @Description Get rates by date
// @Tags rates
// @Accept json
// @Produce json
// @Param date path string true "date. Example: '01-01-2022'"
// @Param code path string false "code. Example: 'USD'"
// @Success 200 {object} js.JSON{rates=[]models.Rate}
// @Failure 500 {object} js.JSON{error=string}
// @Router /currency/{date}/{code} [get]
func (h *Handler) GetRates(w http.ResponseWriter, r *http.Request) {
	log := logger.GetLoggerFromCtx(r.Context()).WithField("op", "Handler.GetRates")

	vars := mux.Vars(r)
	date := vars["date"]
	code, ok := vars["code"]
	if ok {
		rates, err := h.s.Currency.GetRatesByDateAndCode(r.Context(), date, code)
		if err != nil {
			log.WithError(err).Error("failed to get rates")
			js.WriteJSON(w, 500, js.JSON{"error": err}, nil)
			return
		}

		js.WriteJSON(w, 200, js.JSON{"rates": rates}, nil)
		return
	}

	rates, err := h.s.Currency.GetRatesByDate(r.Context(), date)
	if err != nil {
		log.WithError(err).Error("failed to get rates")
		js.WriteJSON(w, 500, js.JSON{"error": err}, nil)
		return
	}

	js.WriteJSON(w, 200, js.JSON{"rates": rates}, nil)
}
