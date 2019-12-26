package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"

	"github.com/buracevs/weather-api/cmd/weather/database"
	_ "github.com/buracevs/weather-api/cmd/weather/docs"
	"github.com/buracevs/weather-api/cmd/weather/handlers"
	httpSwagger "github.com/swaggo/http-swagger"
)

const (
	swaggerHost = "localhost"
	appPort     = ":80"
)

// @title Weather api for weather stations and users
// @version 1.0
// @contact.email buracevs@gmail.com
// @description Swagger api for Weather api
// @licence.name MIT

func main() {
	swaggerUrl := fmt.Sprintf("http://%s%s/swagger/doc.json", swaggerHost, appPort)
	var router = mux.NewRouter()

	dbSrc := database.MakeMssqlDao("localhost", "sa", "qIz5^n78DX)+hY")
	hlndr := handlers.NewHandlerForHttp(dbSrc)

	router.HandleFunc("/{id:[0-9]+}/add-data", logRequests(hlndr.SaveToDataBase)).Methods(http.MethodPost)
	router.HandleFunc("/get-data/{id:[0-9]+}", logRequests(hlndr.GetDataRange)).Methods(http.MethodGet)

	router.PathPrefix("/swagger/").Handler(httpSwagger.Handler(httpSwagger.URL(swaggerUrl)))

	srv := &http.Server{
		Handler:        router,
		Addr:           appPort,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		IdleTimeout:    10 * time.Second,
		MaxHeaderBytes: 1 << 17,
	}

	if err := srv.ListenAndServe(); err != nil {
		log.Println("error occured:", err)
		os.Exit(1)
	}
}

func logRequests(next http.HandlerFunc) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		log.Printf("level=%s", "DEBUG")
		next.ServeHTTP(writer, request)
	}
}
