package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/vkuzmenkova/currency-rates/configs"

	"github.com/rs/zerolog"

	"github.com/gorilla/mux"
	"github.com/vkuzmenkova/currency-rates/internal/controller"
	"github.com/vkuzmenkova/currency-rates/middleware"

	httpSwagger "github.com/swaggo/http-swagger/v2"
	_ "github.com/vkuzmenkova/currency-rates/docs"

	"github.com/rs/zerolog/log"
)

const ShutdownTimeout = 5 * time.Second

// @title           Swagger Currency Rates API
// @version         1.0
// @description     Base currency is USD. The following currencies are supported: USD, EUR, MXN, JPY, GBP, CHF. Info is provided by VAT comply.
// @contact.name   Valentina Kuzmenkova
// @contact.email  valentinakuzmenkova@gmail.com
// @host      localhost:8080
// @BasePath  /api/v1
// @externalDocs.description  OpenAPI
// @externalDocs.url          https://swagger.io/resources/open-api/
func main() {
	ctx := context.Background()
	//time.Local = time.UTC
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	config, err := configs.NewConfig("./configs", "config")
	if err != nil {
		log.Fatal().Err(err).Msg("cannot read config")
	}

	c, err := controller.NewController(ctx, config)
	if err != nil {
		log.Fatal().Err(err).Msg("unable to create controller")
	}

	router := mux.NewRouter().StrictSlash(true)
	subrouter := router.PathPrefix("/api/v1").Subrouter()

	subrouter.Handle(
		"/rates/{code}/update",
		middleware.LoggingRequest(http.HandlerFunc(c.UpdateRate)),
	).Methods(http.MethodPut)
	subrouter.Handle(
		"/rates/{code}",
		middleware.LoggingRequest(http.HandlerFunc(c.GetLastRate)),
	).Methods(http.MethodGet)
	subrouter.Handle(
		"/rates",
		middleware.LoggingRequest(http.HandlerFunc(c.GetRateByUUID)),
	).Methods(http.MethodGet)

	subrouter.PathPrefix("/swagger").Handler(httpSwagger.WrapHandler)

	srv := http.Server{
		Handler:      subrouter,
		Addr:         fmt.Sprintf("%s:%s", config.Service.Host, config.Service.Port),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Info().Msg("Server started at :8080")

	log.Fatal().Err(srv.ListenAndServe())

}
