package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/jlucaspains/sharp-cert-checker/handlers"
	"github.com/jlucaspains/sharp-cert-checker/midlewares"
	"github.com/joho/godotenv"
)

func loadEnv() {
	// outside of local environment, variables should be
	// OS environment variables
	env := os.Getenv("ENV")
	if err := godotenv.Load(); err != nil && env == "" {
		log.Fatal(fmt.Printf("Error loading .env file: %s", err))
	}
}

func main() {
	loadEnv()

	handlers := &handlers.Handlers{}
	handlers.SiteList = handlers.GetConfigSites()

	router := mux.NewRouter()
	router.HandleFunc("/api/check-url", handlers.CheckStatus).Methods("GET")
	router.HandleFunc("/api/site-list", handlers.GetSiteList).Methods("GET")
	router.HandleFunc("/health", handlers.HealthCheck).Methods("GET")
	router.PathPrefix("/").Handler(http.FileServer(http.Dir("./public/")))

	logMiddleware := midlewares.NewLogMiddleware(log.Default())
	router.Use(logMiddleware.Func())

	hostPort, ok := os.LookupEnv("WEB_HOST_PORT")
	if !ok {
		hostPort = ":3000"
	}

	useTls := false
	certFile, ok := os.LookupEnv("TLS_CERT_FILE")
	useTls = ok

	certKeyFile, ok := os.LookupEnv("TLS_CERT_KEY_FILE")
	useTls = useTls && ok

	log.Printf("Starting TLS server on port: %s; use tls: %t", hostPort, useTls)
	if useTls {
		log.Fatalln(http.ListenAndServeTLS(hostPort, certFile, certKeyFile, router))
	} else {
		log.Fatalln(http.ListenAndServe(hostPort, router))
	}
}
