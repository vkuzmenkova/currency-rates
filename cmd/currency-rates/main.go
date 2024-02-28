package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/spf13/viper"
	"github.com/vkuzmenkova/currency-rates/internal/controller"
	"github.com/vkuzmenkova/currency-rates/middleware"

	httpSwagger "github.com/swaggo/http-swagger/v2"
	_ "github.com/vkuzmenkova/currency-rates/docs"
)

// @title           Swagger Currency Rates API
// @version         1.0
// @description     This is a currency rates service.
// @contact.name   Valentina Kuzmenkova
// @contact.email  valentinakuzmenkova@gmail.com
// @host      localhost:8080
// @BasePath  /api/v1
// @externalDocs.description  OpenAPI
// @externalDocs.url          https://swagger.io/resources/open-api/
func main() {
	ctx := context.Background()
	//time.Local = time.UTC

	c, err := controller.NewController(ctx)
	if err != nil {
		log.Fatalln(err)
	}

	router := mux.NewRouter().StrictSlash(true)
	subrouter := router.PathPrefix("/api/v1").Subrouter()

	// указать только нужные параметры base и uuid
	subrouter.Handle(
		"/rates/{code}/update",
		middleware.LoggingRequest(http.HandlerFunc(c.UpdateRate)),
	).Methods(http.MethodGet)
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
		Addr:         fmt.Sprintf("%s:%s", viper.GetString("server.host"), viper.GetString("server.port")),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Println("Server started at :8080")

	log.Fatal(srv.ListenAndServe())

}
