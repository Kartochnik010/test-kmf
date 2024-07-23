package handler

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
)

func Routes(h *Handler) *mux.Router {
	r := mux.NewRouter()
	r.Use(h.AssignLoggerMiddleware, h.rateLimit, h.WriteToConsole)

	r.PathPrefix("/swagger/").Handler(httpSwagger.Handler(
		httpSwagger.URL(fmt.Sprintf("http://localhost:%v/swagger/doc.json", h.cfg.Port)),
		httpSwagger.DeepLinking(true),
		httpSwagger.DocExpansion("none"),
		httpSwagger.DomID("swagger-ui"),
	)).Methods(http.MethodGet)

	// regex so date could only be in this format 06.06.2022 or 01-01-2022
	r.HandleFunc("/currency/save/{date:[0-9]{1,2}.[0-9]{1,2}.[0-9]{4}}", h.FetchAndStore).Methods(http.MethodGet)
	r.HandleFunc("/currency/{date:[0-9]{1,2}.[0-9]{1,2}.[0-9]{4}}/{code:[A-Z]{3}}", h.GetRates).Methods(http.MethodGet)
	r.HandleFunc("/currency/{date:[0-9]{1,2}.[0-9]{1,2}.[0-9]{4}}", h.GetRates).Methods(http.MethodGet)
	return r
}
